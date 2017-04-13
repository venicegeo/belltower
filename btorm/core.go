package btorm

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type Core struct {
	Id        common.Ident `json:"id"         crud:"r"`
	Name      string       `json:"name"       crud:"cru"`
	CreatedAt time.Time    `json:"created_at" crud:"r"`
	UpdatedAt time.Time    `json:"updated_at" crud:"r"`
	IsEnabled bool         `json:"is_enabled" crud:"cru"`
	IsPublic  bool         `json:"is_public"  crud:"cru"`
	OwnerId   common.Ident `json:"owner_id"   crud:"cr"`
}

func (c *Core) GetId() common.Ident      { return c.Id }
func (c *Core) SetId() common.Ident      { c.Id = common.NewId(); return c.Id }
func (c *Core) GetOwnerId() common.Ident { return c.OwnerId }
func (c *Core) GetIsPublic() bool        { return c.IsPublic }

func (c Core) String() string { return fmt.Sprintf("a.%s: %s", c.Id, c.Name) }

func (c *Core) GetCoreMappingProperties() map[string]esorm.MappingProperty {
	properties := map[string]esorm.MappingProperty{
		"id":         esorm.MappingProperty{Type: "keyword"},
		"name":       esorm.MappingProperty{Type: "keyword"},
		"created_at": esorm.MappingProperty{Type: "date"},
		"updated_at": esorm.MappingProperty{Type: "date"},
		"is_enabled": esorm.MappingProperty{Type: "boolean"},
		"is_public":  esorm.MappingProperty{Type: "boolean"},
		"owner_id":   esorm.MappingProperty{Type: "keyword"},
	}

	return properties
}
