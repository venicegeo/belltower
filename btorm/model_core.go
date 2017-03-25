package btorm

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type Common struct {
	Id        common.Ident `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	IsEnabled bool         `json:"is_enabled"`
	IsPublic  bool         `json:"is_public"`
	OwnerId   common.Ident `json:"owner_id"`
}

func (c *Common) GetId() common.Ident      { return c.Id }
func (c *Common) SetId() common.Ident      { c.Id = common.NewId(); return c.Id }
func (c *Common) GetOwnerId() common.Ident { return c.OwnerId }
func (c *Common) GetIsPublic() bool        { return c.IsPublic }

func (c Common) String() string { return fmt.Sprintf("a.%s: %s", c.Id, c.Name) }

func (c *Common) GetCommonMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{
		"id":         esorm.MappingPropertyFields{Type: "string"},
		"name":       esorm.MappingPropertyFields{Type: "string"},
		"created_at": esorm.MappingPropertyFields{Type: "date"},
		"updated_at": esorm.MappingPropertyFields{Type: "date"},
		"is_enabled": esorm.MappingPropertyFields{Type: "boolean"},
		"is_public":  esorm.MappingPropertyFields{Type: "boolean"},
		"owner_id":   esorm.MappingPropertyFields{Type: "string"},
	}

	return properties
}
