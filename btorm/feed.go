package btorm

import (
	"time"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type Feed struct {
	Common
	FeedType        string        `json:"feed_type"`
	PollingInterval time.Duration `json:"polling_interval"`
	MessageCount    uint          `json:"message_count"`
	LastMessageAt   time.Time     `json:"last_message_at"`
	Settings        interface{}   `json:"settings"`
}

//---------------------------------------------------------------------

type FeedFieldsForCreate struct {
	Name      string
	FeedType  string
	IsEnabled bool
	Settings  map[string]interface{}
	IsPublic  bool
}

type FeedFieldsForRead struct {
	Common
	FeedType      string
	Settings      interface{}
	MessageCount  uint
	LastMessageAt time.Time
}

type FeedFieldsForUpdate struct {
	Name      string
	IsEnabled bool
	IsPublic  bool
	Settings  interface{}
}

//---------------------------------------------------------------------

func (feed *Feed) GetLoweredName() string { return "feed" }

func (feed *Feed) GetMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{
		"feed_type":       esorm.MappingPropertyFields{Type: "string"},
		"message_count":   esorm.MappingPropertyFields{Type: "integer"},
		"last_message_at": esorm.MappingPropertyFields{Type: "date"},
		"settings":        esorm.MappingPropertyFields{Type: "object", Dynamic: "true"},
	}

	for k, v := range feed.GetMappingProperties() {
		properties[k] = v
	}

	return properties
}

//---------------------------------------------------------------------

func (feed *Feed) SetFieldsForCreate(ownerId common.Ident, ifields interface{}) error {

	fields := ifields.(*FeedFieldsForCreate)

	feed.Name = fields.Name
	feed.IsEnabled = fields.IsEnabled
	feed.Settings = fields.Settings
	feed.OwnerId = ownerId
	feed.IsPublic = fields.IsPublic
	feed.FeedType = fields.FeedType

	return nil
}

func (feed *Feed) GetFieldsForRead() (interface{}, error) {

	read := &FeedFieldsForRead{
		FeedType:      feed.FeedType,
		Settings:      feed.Settings,
		MessageCount:  feed.MessageCount,
		LastMessageAt: feed.LastMessageAt,
	}

	read.Common = feed.Common

	return read, nil
}

func (feed *Feed) SetFieldsForUpdate(ifields interface{}) error {

	fields := ifields.(*FeedFieldsForUpdate)

	if fields.Name != "" {
		feed.Name = fields.Name
	}

	feed.IsEnabled = fields.IsEnabled
	feed.IsPublic = fields.IsPublic

	return nil
}
