package orm

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
)

type Feed struct {
	ID   uint `gorm:"primary_key"`
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

	FeedType              string
	IsEnabled             bool
	MessageDecl           map[string]interface{} `gorm:"-"`
	MessageDeclInternal   *common.Json           // owned by system
	Settings              map[string]interface{} `gorm:"-"`
	SettingsInternal      *common.Json           // owned by system
	MessageCount          uint
	LastMessageAt         *time.Time
	LastMessageAtInternal string
	Owner                 User
	OwnerID               uint
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

	f.LastMessageAtInternal = f.LastMessageAt.Format(time.RFC3339)

	return nil
}

func (f Feed) String() string {
	s := fmt.Sprintf("f.%d: %s", f.ID, f.Name)
	return s
}
