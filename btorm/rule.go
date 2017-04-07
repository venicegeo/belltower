package btorm

import (
	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

type Rule struct {
	Core
	Expression string `json:"expression"` // CR
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

//---------------------------------------------------------------------

func (rule *Rule) SetFieldsForCreate(ownerId common.Ident, ifields interface{}) error {

	fields := ifields.(*Rule)

	err := rule.Core.SetFieldsForCreate(ownerId, &fields.Core)
	if err != nil {
		return err
	}

	rule.Expression = fields.Expression

	return nil
}

func (rule *Rule) GetFieldsForRead() (interface{}, error) {

	core, err := rule.Core.GetFieldsForRead()
	if err != nil {
		return nil, err
	}

	fields := &Rule{
		Core:       core,
		Expression: rule.Expression,
	}

	return fields, nil
}

func (rule *Rule) SetFieldsForUpdate(ifields interface{}) error {

	fields := ifields.(*Rule)

	err := rule.Core.SetFieldsForUpdate(&fields.Core)
	if err != nil {
		return nil
	}

	return nil
}
