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

	FeedType            string
	IsEnabled           bool
	MessageDecl         map[string]interface{} `gorm:"-"`
	MessageDeclInternal *common.Json           // owned by system
	Settings            map[string]interface{} `gorm:"-"`
	SettingsInternal    *common.Json           // owned by system
	MessageCount        uint
	LastMessageAt       *time.Time
	OwnerID             uint
}

//---------------------------------------------------------------------

type FeedFieldsForCreate struct {
	Name        string
	FeedType    string
	IsEnabled   bool
	MessageDecl map[string]interface{} `gorm:"-"`
	Settings    map[string]interface{} `gorm:"-"`
}

type FeedFieldsForRead struct {
	ID   uint
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

	FeedType      string
	IsEnabled     bool
	MessageDecl   map[string]interface{}
	Settings      map[string]interface{}
	MessageCount  uint
	LastMessageAt *time.Time
	OwnerID       uint
}

type FeedFieldsForUpdate struct {
	Name      string
	IsEnabled bool
}

func CreateFeed(fields *FeedFieldsForCreate) (*Feed, error) {
	feed := &Feed{
		Name:        fields.Name,
		FeedType:    fields.FeedType,
		IsEnabled:   fields.IsEnabled,
		MessageDecl: fields.MessageDecl,
		Settings:    fields.Settings,
	}

	return feed, nil
}

//---------------------------------------------------------------------

func (feed *Feed) Read() (*FeedFieldsForRead, error) {
	read := &FeedFieldsForRead{
		ID:   feed.ID,
		Name: feed.Name,

		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		//DeletedAt *time.Time `sql:"index"`

		FeedType:      feed.FeedType,
		IsEnabled:     feed.IsEnabled,
		MessageDecl:   feed.MessageDecl,
		Settings:      feed.Settings,
		MessageCount:  feed.MessageCount,
		LastMessageAt: feed.LastMessageAt,
		OwnerID:       feed.OwnerID,
	}

	return read, nil
}

func (feed *Feed) Update(update *FeedFieldsForUpdate) error {

	if update.Name != "" {
		feed.Name = update.Name
	}

	feed.IsEnabled = update.IsEnabled

	return nil
}

func (f *Feed) BeforeSave() error {
	var err error

	f.MessageDeclInternal, err = common.NewJsonFromMap(f.MessageDecl)
	if err != nil {
		return err
	}

	f.SettingsInternal, err = common.NewJsonFromMap(f.Settings)
	if err != nil {
		return err
	}

	return nil
}

func (f Feed) String() string {
	s := fmt.Sprintf("f.%d: %s", f.ID, f.Name)
	return s
}
