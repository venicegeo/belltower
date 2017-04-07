package feeders

import (
	"log"
	"time"

	"fmt"

	"github.com/venicegeo/belltower/btorm"
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
	Create(*btorm.Feed) (Feeder, error) // factory for this given feeder type (static function)
	Id() common.Ident
	GetName() string
	Poll() (interface{}, error)
	SettingsSchema() map[string]string
	EventSchema() map[string]string
}

type EventPosterFunc func(*Event) error

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

func checkSchema(schema map[string]string, data map[string]interface{}) error {
	for key, typ := range schema {
		v, ok := data[key]
		if !ok {
			return fmt.Errorf("Settings field '%s' not present", key)
		}

		log.Printf("VVV %T", v)
		switch typ {
		case "integer":
			_, ok := v.(int)
			if !ok {
				return fmt.Errorf("Settings field '%s' is not value of type '%s'", key, typ)
			}
		default:
			return fmt.Errorf("Settings field '%s' has unknown type '%s'", key, typ)
		}
	}
	return nil
}

//---------------------------------------------------------------------

type FeederFactoryFunc func(*btorm.Feed) (Feeder, error)

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

func (f *FeederFactory) create(feed *btorm.Feed) (Feeder, error) {
	factoryFunc := f.factories[feed.FeederId]

	if factoryFunc == nil {
		return nil, fmt.Errorf("no such factory: %s", feed.FeederId)
	}

	feeder, err := factoryFunc(feed)
	if err != nil {
		return nil, err
	}
	return feeder, nil
}

//---------------------------------------------------------------------
