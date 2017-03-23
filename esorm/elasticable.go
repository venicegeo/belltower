package esorm

import "github.com/venicegeo/belltower/common"

//---------------------------------------------------------------------

// every object type wil be stored in its own type in its own index
type Elasticable interface {
	GetIndexName() string
	GetTypeName() string
	GetMapping() string

	GetId() common.Ident
	SetId() common.Ident

	SetFieldsForCreate(ownerId common.Ident, fields interface{}) error
	GetFieldsForRead() (interface{}, error)
	SetFieldsForUpdate(fields interface{}) error

	String() string
}
