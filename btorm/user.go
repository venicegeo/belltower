package btorm

import (
	"time"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

// Users are never publically visible: only the admin has access rights.
type User struct {
	Common
	Role        Role      `json:"role"`
	LastLoginAt time.Time `json:"last_login_at"`
}

//---------------------------------------------------------------------

type UserFieldsForCreate struct {
	Name      string
	IsEnabled bool
	Role      Role
}

type UserFieldsForRead struct {
	Common
	Role        Role
	LastLoginAt time.Time
}

type UserFieldsForUpdate struct {
	Name      string
	IsEnabled bool
	Role      Role
}

//---------------------------------------------------------------------

func (user *User) GetLoweredName() string { return "user" }

func (user *User) GetMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{
		"role":          esorm.MappingPropertyFields{Type: "string"},
		"last_login_at": esorm.MappingPropertyFields{Type: "date"},
	}

	for k, v := range user.Common.GetCommonMappingProperties() {
		properties[k] = v
	}

	return properties
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
		LastLoginAt: user.LastLoginAt,
		Role:        user.Role,
	}
	read.Common = user.Common
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
