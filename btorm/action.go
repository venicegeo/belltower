package btorm

import (
	"fmt"

	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type Action struct {
	Core
	Settings interface{} `json:"settings" crud:"cr"`
}

func (action *Action) GetIndexName() string { return "action_index" }
func (action *Action) GetTypeName() string  { return "action_type" }

func (action *Action) String() string { return fmt.Sprintf("%#v", action) }

//---------------------------------------------------------------------

func (action *Action) GetMappingProperties() map[string]esorm.MappingProperty {
	properties := map[string]esorm.MappingProperty{}

	properties["settings"] = esorm.MappingProperty{Type: "object", Dynamic: "true"}

	for k, v := range action.Core.GetCoreMappingProperties() {
		properties[k] = v
	}

	return properties
}

//---------------------------------------------------------------------
