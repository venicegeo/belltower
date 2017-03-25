package btorm

import (
	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type Action struct {
	Common
	Settings interface{} `json:"settings"`
}

//---------------------------------------------------------------------

type ActionFieldsForCreate struct {
	Name      string
	IsEnabled bool
	Settings  interface{}
	IsPublic  bool
}

type ActionFieldsForRead struct {
	Common
	Settings interface{}
}

type ActionFieldsForUpdate struct {
	Name      string
	IsEnabled bool
	IsPublic  bool
}

//---------------------------------------------------------------------

func (action *Action) GetLoweredName() string { return "action" }

func (action *Action) GetMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{
		"settings": esorm.MappingPropertyFields{Type: "object", Dynamic: "true"},
	}

	for k, v := range action.Common.GetCommonMappingProperties() {
		properties[k] = v
	}

	return properties
}

//---------------------------------------------------------------------

func (action *Action) SetFieldsForCreate(ownerId common.Ident, ifields interface{}) error {

	fields := ifields.(*ActionFieldsForCreate)

	action.Name = fields.Name
	action.IsEnabled = fields.IsEnabled
	action.Settings = fields.Settings
	action.OwnerId = ownerId
	action.IsPublic = fields.IsPublic

	return nil
}

func (action *Action) GetFieldsForRead() (interface{}, error) {

	read := &ActionFieldsForRead{
		Settings: action.Settings,
	}
	read.Common = action.Common

	return read, nil
}

func (action *Action) SetFieldsForUpdate(ifields interface{}) error {

	fields := ifields.(*ActionFieldsForUpdate)

	if fields.Name != "" {
		action.Name = fields.Name
	}

	action.IsEnabled = fields.IsEnabled
	action.IsPublic = fields.IsPublic

	return nil
}
