package feeders

import (
	"time"

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
