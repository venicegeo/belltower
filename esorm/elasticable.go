package esorm

import (
	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

// every object type wil be stored in its own type in its own index
type Elasticable interface {
	GetLoweredName() string
	GetIndexName() string
	GetTypeName() string
	GetMappingProperties() map[string]MappingPropertyFields

	GetId() common.Ident
	SetId() common.Ident

	SetFieldsForCreate(ownerId common.Ident, fields interface{}) error
	GetFieldsForRead() (interface{}, error)
	SetFieldsForUpdate(fields interface{}) error

	String() string
}
