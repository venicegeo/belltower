package orm2

import (
	"fmt"
	"time"
)

type Action struct {
	Id        string                 `json:"id"`
	Name      string                 `json:"name"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	IsEnabled bool                   `json:"is_enabled"`
	IsPublic  bool                   `json:"is_public"`
	Settings  map[string]interface{} `json:"settings"`
	OwnerId   string                 `json:"owner_id"`
}

//---------------------------------------------------------------------

type ActionFieldsForCreate struct {
	Id        string
	Name      string
	IsEnabled bool
	Settings  map[string]interface{}
	IsPublic  bool
}

type ActionFieldsForRead struct {
	Id   string
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time

	IsEnabled bool
	Settings  map[string]interface{}
	OwnerId   string
	IsPublic  bool
}

type ActionFieldsForUpdate struct {
	Id        string
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

func (action *Action) GetId() string {
	return action.Id
}

func (action *Action) SetId() string {
	action.Id = NewId()
	return action.Id
}

//---------------------------------------------------------------------

func CreateAction(requestorId string, fields *ActionFieldsForCreate) (*Action, error) {

	action := &Action{
		Id:        fields.Id,
		Name:      fields.Name,
		IsEnabled: fields.IsEnabled,
		Settings:  fields.Settings,
		OwnerId:   requestorId,
		IsPublic:  fields.IsPublic,
	}

	return action, nil
}

func (action *Action) GetOwnerId() string {
	return action.OwnerId
}

func (action *Action) GetIsPublic() bool {
	return action.IsPublic
}

//---------------------------------------------------------------------

func (action *Action) Read() (*ActionFieldsForRead, error) {

	read := &ActionFieldsForRead{
		Id:   action.Id,
		Name: action.Name,

		CreatedAt: action.CreatedAt,
		UpdatedAt: action.UpdatedAt,

		IsEnabled: action.IsEnabled,
		Settings:  action.Settings,
		OwnerId:   action.OwnerId,
		IsPublic:  action.IsPublic,
	}

	return read, nil
}

func (action *Action) Update(update *ActionFieldsForUpdate) error {

	if update.Name != "" {
		action.Name = update.Name
	}

	action.IsEnabled = update.IsEnabled
	action.IsPublic = update.IsPublic

	return nil
}

//---------------------------------------------------------------------

func (a Action) String() string {
	s := fmt.Sprintf("a.%d: %s", a.Id, a.Name)
	return s
}
