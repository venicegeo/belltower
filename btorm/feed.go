package btorm

import (
	"time"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type Feed struct {
	Core
	FeederId        common.Ident `json:"feeder_id"`        // CR
	PollingInterval uint         `json:"polling_interval"` // CR // in seconds
	PollingEndAt    time.Time    `json:"polling_end_at"`   // CR
	MessageCount    uint         `json:"message_count"`    // R
	LastMessageAt   time.Time    `json:"last_message_at"`  // R
	Settings        interface{}  `json:"settings"`         // CR
}

//---------------------------------------------------------------------

func (feed *Feed) GetMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{}

	properties["feeder_id"] = esorm.MappingPropertyFields{Type: "keyword"}
	properties["message_count"] = esorm.MappingPropertyFields{Type: "integer"}
	properties["polling_interval"] = esorm.MappingPropertyFields{Type: "integer"}
	properties["polling_end_at"] = esorm.MappingPropertyFields{Type: "date"}
	properties["last_message_at"] = esorm.MappingPropertyFields{Type: "date"}
	properties["settings"] = esorm.MappingPropertyFields{Type: "object", Dynamic: "true"}

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
	feed.PollingEndAt = fields.PollingEndAt
	feed.Settings = fields.Settings
	feed.FeederId = fields.FeederId

	return nil
}

func (feed *Feed) GetFieldsForRead() (interface{}, error) {

	core, err := feed.Core.GetFieldsForRead()
	if err != nil {
		return nil, err
	}

	fields := &Feed{
		Core:          core,
		FeederId:      feed.FeederId,
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
