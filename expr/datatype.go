package expr

import (
	"encoding/json"
	"fmt"
	"log"
)

type DataType struct {
	Name    string      `json:"name"`
	Type    TypeName    `json:"type"`
	Element *DataType   `json:"element,omitempty"` // for maps and arrays
	Fields  []*DataType `json:"fields,omitempty"`  // for structs
}

func NewScalarDataType(name string, typ TypeName) *DataType {
	return &DataType{
		Name: name,
		Type: typ,
	}
}

func NewMapDataType(name string, elem *DataType) *DataType {
	return &DataType{
		Name:    name,
		Type:    TypeNameMap,
		Element: elem,
	}
}

func NewArrayDataType(name string, elem *DataType) *DataType {
	return &DataType{
		Name:    name,
		Type:    TypeNameArray,
		Element: elem,
	}
}

func NewStructDataType(name string, fields []*DataType) *DataType {
	return &DataType{
		Name:   name,
		Type:   TypeNameStruct,
		Fields: fields,
	}
}

func NewDataTypeFromJSON(jsn string) (*DataType, error) {

	m := &DataType{}
	err := json.Unmarshal([]byte(jsn), m)
	if err != nil {
		return nil, err
	}

	log.Printf("%#v", m)
	err = m.validate()
	if err != nil {
		return nil, err
	}

	return m, nil
	//	return NewDataTypeFromMap(m)
}

/*
func NewDataTypeFromMap(m *map[string]interface{}) (*DataType, error) {

	return parseMap(m)
}

func parseMap(m *map[string]interface{}) (*DataType, error) {

	// example input:
	//   {"x": "int", "y":"string"}
	// output:
	//   dt: <struct: <int>, <string>>

	dt := &DataType{
		Type:   TypeEnumStruct,
		Fields: []*DataType{},
	}

	keys := []string{}
	for k := range *m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := (*m)[k]

		//rawType := reflect.TypeOf(v).Kind()
		//typ := TypeEnumInvalid

		var t *DataType

		switch v.(type) {
		case string:
			switch v.(string) {
			case "int":
				t = NewDataTypeInt()
			case "float":
				t = NewDataTypeFloat()
			case "bool":
				t = NewDataTypeBool()
			case "string":
				t = NewDataTypeString()
			default:
				return nil, fmt.Errorf("invalid type string: %s", v)
			}
		default:
			return nil, fmt.Errorf("unknown type of map value: %s", v)
		}

		dt.Fields = append(dt.Fields, t)
	}
	return dt, nil
}
*/

func (dt *DataType) validate() error {
	if dt.Name == "" {
		return fmt.Errorf("validation failed: Name is not set")
	}
	if !isValidTypeName(string(dt.Type)) {
		return fmt.Errorf("validation failed for field %s: type %s is invalid", dt.Name, dt.Type)
	}
	if dt.Type == TypeNameArray || dt.Type == TypeNameMap {
		if dt.Fields != nil {
			return fmt.Errorf("validation failed for field %s: type %s has Fields set", dt.Name, dt.Type)
		}
		if dt.Element == nil {
			return fmt.Errorf("validation failed for field %s: type %s does not have Element set", dt.Name, dt.Type)
		}
	}
	if dt.Type == TypeNameStruct {
		if dt.Element != nil {
			return fmt.Errorf("validation failed for field %s: type %s has Elements set", dt.Name, dt.Type)
		}
		if dt.Fields == nil {
			return fmt.Errorf("validation failed for field %s: type %s does not have Fields set", dt.Name, dt.Type)
		}
	}
	return nil
}

func (dt *DataType) String() string {
	switch dt.Type {
	case TypeNameInt, TypeNameFloat, TypeNameBool, TypeNameString:
		return fmt.Sprintf("<%s/%s>", dt.Name, dt.Type)
	case TypeNameMap, TypeNameArray:
		return fmt.Sprintf("<%s/%s: %s>", dt.Name, dt.Type, dt.Element.String())
	case TypeNameStruct:
		s := ""
		for _, f := range dt.Fields {
			if s != "" {
				s += ", "
			}
			s += f.String()
		}
		return fmt.Sprintf("<%s/%s: %s>", dt.Name, dt.Type, s)
	}
	log.Printf("bad type for %s: %s", dt.Name, dt.Type.String())
	panic(9)
}
