package btorm

import (
	"time"

	"github.com/venicegeo/belltower/esorm"
)

// Users are never publically visible: only the admin has access rights.
type User struct {
	Core
	Role        Role      `json:"role"          crud:"cr"`
	LastLoginAt time.Time `json:"last_login_at" crud:"r"`
}

//---------------------------------------------------------------------

func (user *User) GetMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{}

	properties["role"] = esorm.MappingPropertyFields{Type: "keyword"}
	properties["last_login_at"] = esorm.MappingPropertyFields{Type: "date"}

	for k, v := range user.Core.GetCoreMappingProperties() {
		properties[k] = v
	}

	return properties
}
