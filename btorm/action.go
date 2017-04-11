package btorm

import (
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type Action struct {
	Core
	Settings interface{} `json:"settings" crud:"cr"`
}

//---------------------------------------------------------------------

func (action *Action) GetMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{}

	properties["settings"] = esorm.MappingPropertyFields{Type: "object", Dynamic: "true"}

	for k, v := range action.Core.GetCoreMappingProperties() {
		properties[k] = v
	}

	return properties
}

//---------------------------------------------------------------------
