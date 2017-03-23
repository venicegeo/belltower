package orm2

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
)

type Action struct {
	Id        common.Ident           `json:"id"`
	Name      string                 `json:"name"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	IsEnabled bool                   `json:"is_enabled"`
	IsPublic  bool                   `json:"is_public"`
	Settings  map[string]interface{} `json:"settings"`
	OwnerId   common.Ident           `json:"owner_id"`
}

//---------------------------------------------------------------------

type ActionFieldsForCreate struct {
	Id        common.Ident
	Name      string
	IsEnabled bool
	Settings  map[string]interface{}
	IsPublic  bool
}

type ActionFieldsForRead struct {
	Id   common.Ident
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time

	IsEnabled bool
	Settings  map[string]interface{}
	OwnerId   common.Ident
	IsPublic  bool
}

type ActionFieldsForUpdate struct {
	Id        common.Ident
	Name      string
	IsEnabled bool
	IsPublic  bool
}

//---------------------------------------------------------------------

func (action *Action) GetIndexName() string {
	return "action_index"
}

func (action *Action) GetTypeName() string {
	return "action_type"
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
				},
			}
		}
	}
}`

	return mapping
}

func (action *Action) GetId() common.Ident {
	return action.Id
}

func (action *Action) SetId() common.Ident {
	action.Id = common.NewId()
	return action.Id
}

//---------------------------------------------------------------------

func (action *Action) SetFieldsForCreate(ownerId common.Ident, ifields interface{}) error {

	fields := ifields.(*ActionFieldsForCreate)

	action.Id = fields.Id
	action.Name = fields.Name
	action.IsEnabled = fields.IsEnabled
	action.Settings = fields.Settings
	action.OwnerId = ownerId
	action.IsPublic = fields.IsPublic

	return nil
}

func (action *Action) GetFieldsForRead() (interface{}, error) {

	read := &ActionFieldsForRead{
		Id:        action.Id,
		Name:      action.Name,
		CreatedAt: action.CreatedAt,
		UpdatedAt: action.UpdatedAt,
		IsEnabled: action.IsEnabled,
		Settings:  action.Settings,
		OwnerId:   action.OwnerId,
		IsPublic:  action.IsPublic,
	}

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

//---------------------------------------------------------------------

func (action *Action) GetOwnerId() common.Ident {
	return action.OwnerId
}

func (action *Action) GetIsPublic() bool {
	return action.IsPublic
}

func (a Action) String() string {
	s := fmt.Sprintf("a.%s: %s", a.Id, a.Name)
	return s
}
