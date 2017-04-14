package feeders

import (
	"log"
	"time"

	"fmt"

	"github.com/venicegeo/belltower/btorm"
	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

type Feeder interface {
	GetId() common.Ident
	GetName() string
	GetSettingsSchema() map[string]string
	GetEventSchema() map[string]string
	Poll() (interface{}, error)
}

type CreateFunc func(*btorm.Feed) (Feeder, error)

type FeederInfo struct {
	FeederId    common.Ident
	Description string
	Create      CreateFunc
}

//---------------------------------------------------------------------

type FeederRegistry struct {
	data map[common.Ident]*FeederInfo
}

var feederRegistry = FeederRegistry{}

func (f *FeederRegistry) register(info *FeederInfo) {
	if f.data == nil {
		f.data = map[common.Ident]*FeederInfo{}
	} else if _, ok := f.data[info.FeederId]; ok {
		panic("feeder already registered: " + info.FeederId)
	}
	f.data[info.FeederId] = info
}

//---------------------------------------------------------------------

type Event struct {
	TimeStamp time.Time
	FeedId    common.Ident
	FeederId  common.Ident
	Data      interface{}
}

type EventPosterFunc func(*Event) error

//---------------------------------------------------------------------

func RunFeed(feed *btorm.Feed, feeder Feeder, post EventPosterFunc) error {
	errCount := 0

	for {
		now := time.Now()
		if !feed.PollingEndAt.IsZero() && now.After(feed.PollingEndAt) {
			return fmt.Errorf("end date reached")
		}

		if errCount == 3 {
			log.Printf("Feeder aborting (%s)", feeder.GetName())
			return fmt.Errorf("Feeder aborting -- too many errors")
		}

		eventData, err := feeder.Poll()

		if err != nil {
			log.Printf("Feeder poll error (%s): %v", feeder.GetName(), err)
			errCount++
		} else if eventData == nil {
			// no error and no data - not a hit
		} else {
			// a hit! a hit!
			event := &Event{
				TimeStamp: time.Now(),
				FeedId:    feed.Id,
				FeederId:  feed.FeederId,
				Data:      eventData,
			}
			err = post(event)
			if err != nil {
				log.Printf("Feeder post error (%s): %v", feeder.GetName(), err)
				errCount++
			}
		}

		d := time.Duration(feed.PollingInterval)
		time.Sleep(d * time.Second)
	}

	// not reached
}
