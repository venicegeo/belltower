package btorm

import (
	"github.com/venicegeo/belltower/esorm"
)

type Rule struct {
	Core
	Expression string `json:"expression" crud:"cr"`
}

//---------------------------------------------------------------------

func (rule *Rule) GetMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{}

	properties["expression"] = esorm.MappingPropertyFields{Type: "text"}

	for k, v := range rule.Core.GetCoreMappingProperties() {
		properties[k] = v
	}

	return properties
}
