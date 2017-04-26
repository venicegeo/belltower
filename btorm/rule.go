package btorm

import (
	"fmt"

	"github.com/venicegeo/belltower/esorm"
)

type Rule struct {
	Core
	Expression string `json:"expression" crud:"cr"`
}

func (rule *Rule) GetIndexName() string { return "rule_index" }
func (rule *Rule) GetTypeName() string  { return "rule_type" }

func (rule *Rule) String() string { return fmt.Sprintf("%#v", rule) }

//---------------------------------------------------------------------

func (rule *Rule) GetMappingProperties() map[string]esorm.MappingProperty {
	properties := map[string]esorm.MappingProperty{}

	properties["expression"] = esorm.MappingProperty{Type: "text"}

	for k, v := range rule.Core.GetCoreMappingProperties() {
		properties[k] = v
	}

	return properties
}
