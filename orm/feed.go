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

	FeedType      string
	IsEnabled     bool
	SettingsJson  *common.Json
	MessageCount  uint
	LastMessageAt *time.Time
	OwnerID       uint
	IsPublic      bool
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
	Settings  interface{}
	IsPublic  bool
}

type FeedFieldsForRead struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Name      string
	FeedType  string
	IsEnabled bool
	Settings  interface{}
	OwnerID   uint

	MessageCount  uint
	LastMessageAt *time.Time
	IsPublic      bool
}

type FeedFieldsForUpdate struct {
	Name      string
	IsEnabled bool
	IsPublic  bool
	Settings  interface{}
}

func CreateFeed(requestorID uint, fields *FeedFieldsForCreate) (*Feed, error) {
	settingsJson, err := common.NewJsonFromObject(fields.Settings)
	if err != nil {
		return nil, err
	}
	feed := &Feed{
		Name:         fields.Name,
		FeedType:     fields.FeedType,
		IsEnabled:    fields.IsEnabled,
		SettingsJson: settingsJson,
		OwnerID:      requestorID,
		IsPublic:     fields.IsPublic,
	}

	return feed, nil
}

func (feed *Feed) GetOwnerID() uint {
	return feed.OwnerID
}

func (feed *Feed) GetIsPublic() bool {
	return feed.IsPublic
}

//---------------------------------------------------------------------

func (feed *Feed) Read() (*FeedFieldsForRead, error) {

	var settingsMap map[string]interface{}
	if feed.SettingsJson != nil {
		settingsMap = feed.SettingsJson.AsMap()
	}

	read := &FeedFieldsForRead{
		ID:   feed.ID,
		Name: feed.Name,

		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,

		FeedType:      feed.FeedType,
		IsEnabled:     feed.IsEnabled,
		Settings:      settingsMap,
		MessageCount:  feed.MessageCount,
		LastMessageAt: feed.LastMessageAt,
		OwnerID:       feed.OwnerID,
		IsPublic:      feed.IsPublic,
	}

	return read, nil
}

func (feed *Feed) Update(update *FeedFieldsForUpdate) error {

	if update.Name != "" {
		feed.Name = update.Name
	}

	feed.IsEnabled = update.IsEnabled
	feed.IsPublic = update.IsPublic

	return nil
}

func (f Feed) String() string {
	s := fmt.Sprintf("f.%d: %s", f.ID, f.Name)
	return s
}
