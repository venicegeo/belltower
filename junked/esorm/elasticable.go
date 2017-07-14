package esorm

import (
	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

// every object type wil be stored in its own type in its own index
type Elasticable interface {
	GetId() common.Ident
	SetId(common.Ident)

	GetIndexName() string
	GetTypeName() string

	GetMappingProperties() map[string]MappingProperty

	String() string
}
