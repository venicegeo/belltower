package feeders

import (
	"log"
	"time"

	"fmt"

	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

type Event struct {
	TimeStamp time.Time
	FeedId    common.Ident
	Data      interface{}
}

type Feeder interface {
	Id() common.Ident
	GetName() string
	Poll(feedId common.Ident, settings interface{}) (*Event, error)
	SettingsSchema() map[string]string
	EventSchema() map[string]string
}

var feedersRegistry map[common.Ident]Feeder = map[common.Ident]Feeder{}

func registerFeed(id common.Ident, feeder Feeder) {
	feedersRegistry[id] = feeder
}

func RunFeed(feed *Feed, post func(*Event) error) error {
	feeder := feedersRegistry[feed.FeederId]

	errCount := 0

	for {
		if errCount == 3 {
			log.Printf("Feeder aborting (%s)", feeder.GetName())
			return fmt.Errorf("Feeder aborting -- too many errors")
		}

		event, err := feeder.Poll(feed.Id, feed.SettingsValues)
		if err != nil {
			log.Printf("Feeder poll error (%s): %v", feeder.GetName(), err)
			errCount++
		} else if event == nil {
			// TODO: internal error -- if err nil, then event should not be nil
			errCount++
		} else {
			err = post(event)
			if err != nil {
				log.Printf("Feeder post error (%s): %v", feeder.GetName(), err)
				errCount++
			}
		}

		time.Sleep(feed.Interval)
	}

	// not reached
}

//---------------------------------------------------------------------

type Feed struct {
	Id             common.Ident
	Name           string
	FeederId       common.Ident
	Interval       time.Duration
	SettingsValues interface{}
}
