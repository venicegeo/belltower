package esorm

//---------------------------------------------------------------------

type Mapping struct {
	Settings map[string]interface{}     `json:"settings"`
	Mappings map[string]MappingProperty `json:"mappings"`
}

type MappingProperty struct {
	Type       string                     `json:"type,omitempty"`
	Dynamic    string                     `json:"dynamic,omitempty"`
	Properties map[string]MappingProperty `json:"properties,omitempty"`
}

func NewMapping(e Elasticable) *Mapping {
	mt := MappingProperty{
		Dynamic:    "strict",
		Properties: e.GetMappingProperties(),
	}

	m := &Mapping{
		Settings: map[string]interface{}{},
		Mappings: map[string]MappingProperty{GetTypeName(e): mt},
	}

	return m
}
