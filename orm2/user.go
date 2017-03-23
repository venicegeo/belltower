package orm2

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
)

// Users are never publically visible: only the admin has access rights.
type User struct {
	Id          common.Ident `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Name        string       `json:"name"`
	Role        Role         `json:"role"`
	IsEnabled   bool         `json:"is_enabled"`
	LastLoginAt time.Time    `json:"last_login_at"`
	OwnerId     common.Ident `json:"owner_id"`
}

//---------------------------------------------------------------------

type UserFieldsForCreate struct {
	Name      string
	IsEnabled bool
	Role      Role
}

type UserFieldsForRead struct {
	Id          common.Ident
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	IsEnabled   bool
	Role        Role
	LastLoginAt time.Time
}

type UserFieldsForUpdate struct {
	Name      string
	IsEnabled bool
	Role      Role
}

//---------------------------------------------------------------------

func (user *User) GetIndexName() string {
	return "user_index"
}

func (user *User) GetTypeName() string {
	return "user_type"
}

func (user *User) GetMapping() string {
	mapping := `{
	"settings":{
	},
	"mappings":{
		"user_type":{
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
				"role":{
					"type":"string"
				},
				"last_login_at":{
					"type":"date"
				},
				"owner_id":{
					"type":"string"
				}
			}
		}
	}
}`

	return mapping
}

func (user *User) GetId() common.Ident {
	return user.Id
}

func (user *User) SetId() common.Ident {
	user.Id = common.NewId()
	return user.Id
}

//---------------------------------------------------------------------

func (user *User) SetFieldsForCreate(ownerId common.Ident, ifields interface{}) error {

	fields := ifields.(*UserFieldsForCreate)

	user.Name = fields.Name
	user.IsEnabled = fields.IsEnabled
	user.Role = fields.Role

	return nil
}

func (user *User) GetFieldsForRead() (interface{}, error) {

	read := &UserFieldsForRead{
		Id:          user.Id,
		Name:        user.Name,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsEnabled:   user.IsEnabled,
		LastLoginAt: user.LastLoginAt,
		Role:        user.Role,
	}

	return read, nil
}

func (user *User) SetFieldsForUpdate(ifields interface{}) error {

	fields := ifields.(*UserFieldsForUpdate)

	if fields.Name != "" {
		user.Name = fields.Name
	}

	user.IsEnabled = fields.IsEnabled
	user.Role = fields.Role

	return nil
}

//---------------------------------------------------------------------

func (user *User) GetOwnerId() common.Ident {
	return user.OwnerId
}

func (user *User) GetIsPublic() bool {
	panic(1)
}

func (user User) String() string {
	s := fmt.Sprintf("a.%s: %s", user.Id, user.Name)
	return s
}
