package btorm

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

type Feed struct {
	Id            common.Ident           `json:"id"`
	Name          string                 `json:"name"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	FeedType      string                 `json:"feed_type"`
	IsEnabled     bool                   `json:"is_enabled"`
	Settings      map[string]interface{} `json:"settings"`
	MessageCount  uint                   `json:"message_count"`
	LastMessageAt time.Time              `json:"last_message_at"`
	OwnerId       common.Ident           `json:"owner_id"`
	IsPublic      bool                   `json:"is_public"`
}

type FeedEvent interface{}
type FeedSettings interface{}
type Feeder interface {
	FeedType() string
	Poll() (FeedEvent, error)
	Post(FeedEvent) error
	Sleep() error
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
	Id            common.Ident
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	FeedType      string
	IsEnabled     bool
	Settings      map[string]interface{}
	OwnerId       common.Ident
	MessageCount  uint
	LastMessageAt time.Time
	IsPublic      bool
}

type FeedFieldsForUpdate struct {
	Name      string
	IsEnabled bool
	IsPublic  bool
	Settings  map[string]interface{}
}

//---------------------------------------------------------------------

func (feed *Feed) GetIndexName() string {
	return "feed_index"
}

func (feed *Feed) GetTypeName() string {
	return "feed_type"
}

func (feed *Feed) GetMapping() string {
	mapping := `{
	"settings":{
	},
	"mappings":{
		"feed_type":{
			"dynamic":"strict",
			"properties":{
				"id":{
					"type":"string"
				},
				"name":{
					"type":"string"
				},
				"created_at":{
					"type":"date"
				},
				"updated_at":{
					"type":"date"
				},
				"is_enabled":{
					"type":"boolean"
				},
				"is_public":{
					"type":"boolean"
				},
				"feed_type":{
					"type":"string"
				},
				"message_count":{
					"type":"integer"
				},
				"last_message_at":{
					"type":"date"
				},
				"owner_id":{
					"type":"string"
				},
				"settings":{
					"dynamic":"true",
					"type":"object"
				}
			}
		}
	}
}`

	return mapping
}

func (feed *Feed) GetId() common.Ident {
	return feed.Id
}

func (feed *Feed) SetId() common.Ident {
	feed.Id = common.NewId()
	return feed.Id
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
		Id:            feed.Id,
		Name:          feed.Name,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		FeedType:      feed.FeedType,
		IsEnabled:     feed.IsEnabled,
		Settings:      feed.Settings,
		MessageCount:  feed.MessageCount,
		LastMessageAt: feed.LastMessageAt,
		OwnerId:       feed.OwnerId,
		IsPublic:      feed.IsPublic,
	}

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

//---------------------------------------------------------------------

func (feed *Feed) GetOwnerId() common.Ident {
	return feed.OwnerId
}

func (feed *Feed) GetIsPublic() bool {
	return feed.IsPublic
}

//---------------------------------------------------------------------

func (f Feed) String() string {
	s := fmt.Sprintf("f.%s: %s", f.Id, f.Name)
	return s
}
