package feeders

import (
	"time"

	"github.com/venicegeo/belltower/btorm"
	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

type Event struct {
	TimeStamp time.Time
	FeedId    common.Ident
	Data      interface{}
}

type Feeder interface {
	Create(settings interface{}) *Feeder // this is a "static" method
	FeedType() btorm.FeedType
	Poll() (*Event, error)
}

type FeederFactory struct {
	registry map[FeedType]*Feeder
}

func NewFeederFactory(feeders ...*Feeder) (*FeederFactory, error) {
	factory := &FeederFactory{}
	factory.registry = map[FeedType]*Feeder{}

	for _, v := range feeders {
		factory.registry[v.FeedType] = v
	}
}
