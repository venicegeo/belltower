package btorm

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type Core struct {
	Id        common.Ident `json:"id"`         // R
	Name      string       `json:"name"`       // CRU
	CreatedAt time.Time    `json:"created_at"` // R
	UpdatedAt time.Time    `json:"updated_at"` // R
	IsEnabled bool         `json:"is_enabled"` // CRU
	IsPublic  bool         `json:"is_public"`  // CRU
	OwnerId   common.Ident `json:"owner_id"`   // CR
}

func (c *Core) GetId() common.Ident      { return c.Id }
func (c *Core) SetId() common.Ident      { c.Id = common.NewId(); return c.Id }
func (c *Core) GetOwnerId() common.Ident { return c.OwnerId }
func (c *Core) GetIsPublic() bool        { return c.IsPublic }

func (c Core) String() string { return fmt.Sprintf("a.%s: %s", c.Id, c.Name) }

func (c *Core) GetCoreMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{
		"id":         esorm.MappingPropertyFields{Type: "keyword"},
		"name":       esorm.MappingPropertyFields{Type: "keyword"},
		"created_at": esorm.MappingPropertyFields{Type: "date"},
		"updated_at": esorm.MappingPropertyFields{Type: "date"},
		"is_enabled": esorm.MappingPropertyFields{Type: "boolean"},
		"is_public":  esorm.MappingPropertyFields{Type: "boolean"},
		"owner_id":   esorm.MappingPropertyFields{Type: "keyword"},
	}

	return properties
}

// c is set from fields, no return
func (c *Core) SetFieldsForCreate(ownerId common.Ident, fields *Core) error {
	c.Name = fields.Name
	c.IsEnabled = fields.IsEnabled
	c.IsPublic = fields.IsPublic
	c.OwnerId = ownerId
	return nil
}

// fields is set from c and returned
func (c *Core) GetFieldsForRead() (Core, error) {
	fields := Core{
		Id:        c.Id,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		IsEnabled: c.IsEnabled,
		IsPublic:  c.IsPublic,
		OwnerId:   c.OwnerId,
	}
	return fields, nil
}

// c is set from fields, no return
func (c *Core) SetFieldsForUpdate(fields *Core) error {
	c.Name = fields.Name
	c.IsEnabled = fields.IsEnabled
	c.IsPublic = fields.IsPublic
	return nil
}
