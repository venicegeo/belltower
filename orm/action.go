package orm

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
)

type Action struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name         string
	IsEnabled    bool
	SettingsJson string
	OwnerID      uint
}

//---------------------------------------------------------------------

type ActionFieldsForCreate struct {
	Name      string
	IsEnabled bool
	Settings  map[string]interface{}
	OwnerID   uint
}

type ActionFieldsForRead struct {
	ID   uint
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time

	IsEnabled bool
	Settings  map[string]interface{}
	OwnerID   uint
}

type ActionFieldsForUpdate struct {
	Name      string
	IsEnabled bool
}

//---------------------------------------------------------------------

func CreateAction(fields *ActionFieldsForCreate) (*Action, error) {
	settingsJson, err := common.NewJsonFromMap(fields.Settings)
	if err != nil {
		return nil, err
	}
	action := &Action{
		Name:         fields.Name,
		IsEnabled:    fields.IsEnabled,
		SettingsJson: settingsJson.AsString(),
		OwnerID:      fields.OwnerID,
	}

	return action, nil
}

//---------------------------------------------------------------------

func (action *Action) Read() (*ActionFieldsForRead, error) {

	var j, err = common.NewJsonFromString(action.SettingsJson)
	if err != nil {
		return nil, err
	}
	var settingsMap map[string]interface{}
	if j.AsString() != "" {
		settingsMap = j.AsMap()
	}

	read := &ActionFieldsForRead{
		ID:   action.ID,
		Name: action.Name,

		CreatedAt: action.CreatedAt,
		UpdatedAt: action.UpdatedAt,

		IsEnabled: action.IsEnabled,
		Settings:  settingsMap,
		OwnerID:   action.OwnerID,
	}

	return read, nil
}

func (action *Action) Update(update *ActionFieldsForUpdate) error {

	if update.Name != "" {
		action.Name = update.Name
	}

	action.IsEnabled = update.IsEnabled

	return nil
}

//---------------------------------------------------------------------

func (a Action) String() string {
	s := fmt.Sprintf("a.%d: %s", a.ID, a.Name)
	return s
}
