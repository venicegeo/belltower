package btorm

import (
	"time"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type FeedType string

type Feed struct {
	Core
	FeedType        FeedType      `json:"feed_type"`        // CR
	PollingInterval time.Duration `json:"polling_interval"` // CR
	MessageCount    uint          `json:"message_count"`    // R
	LastMessageAt   time.Time     `json:"last_message_at"`  // R
	Settings        interface{}   `json:"settings"`         // CR
}

//---------------------------------------------------------------------

func (feed *Feed) GetLoweredName() string { return "feed" }

func (feed *Feed) GetMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{
		"feed_type":        esorm.MappingPropertyFields{Type: "string"},
		"message_count":    esorm.MappingPropertyFields{Type: "integer"},
		"polling_interval": esorm.MappingPropertyFields{Type: "integer"},
		"last_message_at":  esorm.MappingPropertyFields{Type: "date"},
		"settings":         esorm.MappingPropertyFields{Type: "object", Dynamic: "true"},
	}

	for k, v := range feed.Core.GetCoreMappingProperties() {
		properties[k] = v
	}

	return properties
}

//---------------------------------------------------------------------

func (feed *Feed) SetFieldsForCreate(ownerId common.Ident, ifields interface{}) error {

	fields := ifields.(*Feed)

	err := feed.Core.SetFieldsForCreate(ownerId, &fields.Core)
	if err != nil {
		return err
	}

	feed.PollingInterval = fields.PollingInterval
	feed.Settings = fields.Settings
	feed.FeedType = fields.FeedType

	return nil
}

func (feed *Feed) GetFieldsForRead() (interface{}, error) {

	core, err := feed.Core.GetFieldsForRead()
	if err != nil {
		return nil, err
	}

	fields := &Feed{
		Core:          core,
		FeedType:      feed.FeedType,
		Settings:      feed.Settings,
		MessageCount:  feed.MessageCount,
		LastMessageAt: feed.LastMessageAt,
	}

	return fields, nil
}

func (feed *Feed) SetFieldsForUpdate(ifields interface{}) error {

	fields := ifields.(*Feed)

	err := feed.Core.SetFieldsForUpdate(&fields.Core)
	if err != nil {
		return nil
	}

	return nil
}
