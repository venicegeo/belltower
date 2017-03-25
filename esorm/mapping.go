package esorm

//---------------------------------------------------------------------

type Mapping struct {
	Settings MappingSettings            `json:"settings"`
	Mappings map[string]MappingTypeName `json:"mappings"`
}

type MappingSettings map[string]interface{}

type MappingTypeName struct {
	Dynamic    string                           `json:"dynamic,omitempty"`
	Properties map[string]MappingPropertyFields `json:"properties,omitempty"`
}

type MappingPropertyFields struct {
	Type       string                           `json:"type"`
	Dynamic    string                           `json:"dynamic,omitempty"`
	Properties map[string]MappingPropertyFields `json:"properties,omitempty"`
}

func NewMapping(e Elasticable) *Mapping {
	mt := MappingTypeName{
		Dynamic:    "strict",
		Properties: e.GetMappingProperties(),
	}

	m := &Mapping{
		Settings: MappingSettings{},
		Mappings: map[string]MappingTypeName{GetTypeName(e): mt},
	}

	return m
}
