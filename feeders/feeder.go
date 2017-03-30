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
	FeederId  common.Ident
	Data      interface{}
}

type Feeder interface {
	Create(Feed) (Feeder, error) // factory for this given feeder type
	Id() common.Ident
	GetName() string
	Poll() (interface{}, error)
	SettingsSchema() map[string]string
	EventSchema() map[string]string
}

type EventPosterFunc func(*Event) error

func RunFeed(feed Feed, post EventPosterFunc) error {
	feeder, err := feederFactory.create(feed)
	if err != nil {
		return err
	}

	errCount := 0

	for {
		now := time.Now()
		if !feed.EndDate.IsZero() && now.After(feed.EndDate) {
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

		time.Sleep(feed.Interval)
	}

	// not reached
}

//---------------------------------------------------------------------

type FeederFactoryFunc func(Feed) (Feeder, error)

type FeederFactory struct {
	factories map[common.Ident]FeederFactoryFunc
}

var feederFactory = FeederFactory{}

func (f *FeederFactory) register(feeder Feeder) {
	if f.factories == nil {
		f.factories = map[common.Ident]FeederFactoryFunc{}
	}
	_, ok := f.factories[feeder.Id()]
	if ok {
		panic("factory already registered: " + feeder.Id())
	}
	f.factories[feeder.Id()] = feeder.Create
}

func (f *FeederFactory) create(feed Feed) (Feeder, error) {
	factory := f.factories[feed.FeederId]

	if factory == nil {
		return nil, fmt.Errorf("no such factory: %s", feed.FeederId)
	}

	feeder, err := factory(feed)
	if err != nil {
		return nil, err
	}
	return feeder, nil
}

//---------------------------------------------------------------------

type Feed struct {
	Id             common.Ident
	Name           string
	FeederId       common.Ident
	Interval       time.Duration
	EndDate        time.Time
	SettingsValues interface{}
}
