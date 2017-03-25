package esorm

import (
	"fmt"
	"strings"

	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

// every object type wil be stored in its own type in its own index
type Elasticable interface {
	GetMappingProperties() map[string]MappingPropertyFields

	GetId() common.Ident
	SetId() common.Ident

	SetFieldsForCreate(ownerId common.Ident, fields interface{}) error
	GetFieldsForRead() (interface{}, error)
	SetFieldsForUpdate(fields interface{}) error

	String() string
}

func getLoweredName(x interface{}) string {
	s := fmt.Sprintf("%T", x)
	dot := strings.Index(s, ".")
	t := strings.ToLower(s[dot+1:])
	return t
}

func GetIndexName(x interface{}) string { return getLoweredName(x) + "_index" }
func GetTypeName(x interface{}) string  { return getLoweredName(x) + "_type" }
