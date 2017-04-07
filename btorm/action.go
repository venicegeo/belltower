package btorm

import (
	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type Action struct {
	Core
	Settings interface{} `json:"settings"` // CR
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

func (action *Action) SetFieldsForCreate(ownerId common.Ident, ifields interface{}) error {

	fields := ifields.(*Action)

	err := action.Core.SetFieldsForCreate(ownerId, &fields.Core)
	if err != nil {
		return err
	}

	action.Settings = fields.Settings

	return nil
}

func (action *Action) GetFieldsForRead() (interface{}, error) {

	core, err := action.Core.GetFieldsForRead()
	if err != nil {
		return nil, err
	}

	fields := &Action{
		Core:     core,
		Settings: action.Settings,
	}

	return fields, nil
}

func (action *Action) SetFieldsForUpdate(ifields interface{}) error {

	fields := ifields.(*Action)

	err := action.Core.SetFieldsForUpdate(&fields.Core)
	if err != nil {
		return nil
	}

	return nil
}
