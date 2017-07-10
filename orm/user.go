package orm

import (
	"fmt"
	"time"

)

// Users are never publically visible: only the admin has access rights.
type User struct {
	Core
	Role        Role      `json:"role"          crud:"cr"`
	LastLoginAt time.Time `json:"last_login_at" crud:"r"`
}

func (user *User) GetIndexName() string { return "user_index" }
func (user *User) GetTypeName() string  { return "user_type" }

func (user *User) String() string { return fmt.Sprintf("%#v", user) }

//---------------------------------------------------------------------

func (user *User) GetMappingProperties() map[string]esorm.MappingProperty {
	properties := map[string]esorm.MappingProperty{}

	properties["role"] = esorm.MappingProperty{Type: "keyword"}
	properties["last_login_at"] = esorm.MappingProperty{Type: "date"}

	for k, v := range user.Core.GetCoreMappingProperties() {
		properties[k] = v
	}

	return properties
}
