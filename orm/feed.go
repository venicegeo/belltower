package orm

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

type Feed struct {
	ID   uint `gorm:"primary_key"`
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

	FeedType        string
	IsEnabled       bool
	MessageDeclJson *common.Json
	SettingsJson    *common.Json
	MessageCount    uint
	LastMessageAt   *time.Time
	OwnerID         uint
	PublicCanRead   bool
	PublicCanWrite  bool
}

//---------------------------------------------------------------------

type FeedFieldsForCreate struct {
	Name           string
	FeedType       string
	IsEnabled      bool
	MessageDecl    map[string]interface{} `gorm:"-"`
	Settings       map[string]interface{} `gorm:"-"`
	OwnerID        uint
	PublicCanRead  bool
	PublicCanWrite bool
}

type FeedFieldsForRead struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	FeedType    string
	IsEnabled   bool
	MessageDecl map[string]interface{}
	Settings    map[string]interface{}
	OwnerID     uint

	MessageCount   uint
	LastMessageAt  *time.Time
	PublicCanRead  bool
	PublicCanWrite bool
}

type FeedFieldsForUpdate struct {
	Name           string
	IsEnabled      bool
	PublicCanRead  bool
	PublicCanWrite bool
}

func CreateFeed(fields *FeedFieldsForCreate) (*Feed, error) {
	messageJson, err := common.NewJsonFromMap(fields.MessageDecl)
	if err != nil {
		return nil, err
	}
	settingsJson, err := common.NewJsonFromMap(fields.Settings)
	if err != nil {
		return nil, err
	}
	feed := &Feed{
		Name:            fields.Name,
		FeedType:        fields.FeedType,
		IsEnabled:       fields.IsEnabled,
		MessageDeclJson: messageJson,
		SettingsJson:    settingsJson,
		OwnerID:         fields.OwnerID,
		PublicCanRead:   fields.PublicCanRead,
		PublicCanWrite:  fields.PublicCanWrite,
	}

	return feed, nil
}

//---------------------------------------------------------------------

func (feed *Feed) Read() (*FeedFieldsForRead, error) {

	var messageMap, settingsMap map[string]interface{}
	if feed.MessageDeclJson != nil {
		messageMap = feed.MessageDeclJson.AsMap()
	}
	if feed.SettingsJson != nil {
		settingsMap = feed.SettingsJson.AsMap()
	}

	read := &FeedFieldsForRead{
		ID:   feed.ID,
		Name: feed.Name,

		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		//DeletedAt *time.Time `sql:"index"`

		FeedType:       feed.FeedType,
		IsEnabled:      feed.IsEnabled,
		MessageDecl:    messageMap,
		Settings:       settingsMap,
		MessageCount:   feed.MessageCount,
		LastMessageAt:  feed.LastMessageAt,
		OwnerID:        feed.OwnerID,
		PublicCanRead:  feed.PublicCanRead,
		PublicCanWrite: feed.PublicCanWrite,
	}

	return read, nil
}

func (feed *Feed) Update(update *FeedFieldsForUpdate) error {

	if update.Name != "" {
		feed.Name = update.Name
	}

	feed.IsEnabled = update.IsEnabled
	feed.PublicCanRead = update.PublicCanRead
	feed.PublicCanWrite = update.PublicCanWrite

	return nil
}

func (f Feed) String() string {
	s := fmt.Sprintf("f.%d: %s", f.ID, f.Name)
	return s
}
