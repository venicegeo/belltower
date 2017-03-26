package btorm

import (
	"time"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

// Users are never publically visible: only the admin has access rights.
type User struct {
	Core
	Role        Role      `json:"role"`          // CR
	LastLoginAt time.Time `json:"last_login_at"` // R
}

//---------------------------------------------------------------------

func (user *User) GetLoweredName() string { return "user" }

func (user *User) GetMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{
		"role":          esorm.MappingPropertyFields{Type: "string"},
		"last_login_at": esorm.MappingPropertyFields{Type: "date"},
	}

	for k, v := range user.Core.GetCoreMappingProperties() {
		properties[k] = v
	}

	return properties
}

//---------------------------------------------------------------------

func (user *User) SetFieldsForCreate(ownerId common.Ident, ifields interface{}) error {

	fields := ifields.(*User)

	err := user.Core.SetFieldsForCreate(ownerId, &fields.Core)
	if err != nil {
		return err
	}

	user.Role = fields.Role

	return nil
}

func (user *User) GetFieldsForRead() (interface{}, error) {

	core, err := user.Core.GetFieldsForRead()
	if err != nil {
		return nil, err
	}

	fields := &User{
		Core:        core,
		LastLoginAt: user.LastLoginAt,
		Role:        user.Role,
	}

	return fields, nil
}

func (user *User) SetFieldsForUpdate(ifields interface{}) error {

	fields := ifields.(*User)

	err := user.Core.SetFieldsForUpdate(&fields.Core)
	if err != nil {
		return nil
	}

	user.Role = fields.Role

	return nil
}
