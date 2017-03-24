package btorm

import (
	"github.com/venicegeo/belltower/common"
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

func (action *Action) GetLoweredName() string {
	return "action"
}

func (action *Action) GetMapping() string {
	mapping := `{
	"settings":{
	},
	"mappings":{
		"action_type":{
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
