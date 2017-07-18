package NOTUSED

import (
	"encoding/json"
	"fmt"
	"log"
)

//---------------------------------------------------------------------

type TypeName string

const TypeNameInvalid TypeName = "INVALID"
const TypeNameInt TypeName = "int"
const TypeNameFloat TypeName = "float"
const TypeNameBool TypeName = "bool"
const TypeNameString TypeName = "string"
const TypeNameArray TypeName = "array"
const TypeNameMap TypeName = "map"
const TypeNameStruct TypeName = "struct"

func (t TypeName) String() string {
	return string(t)
}

func isValidTypeName(s string) bool {
	switch TypeName(s) {
	case TypeNameInt, TypeNameFloat, TypeNameBool, TypeNameString,
		TypeNameArray, TypeNameMap, TypeNameStruct:
		return true
	}
	return false
}

func TypeNameFromString(s string) TypeName {
	if isValidTypeName(s) {
		return TypeName(s)
	}
	return TypeNameInvalid
}

//---------------------------------------------------------------------

type DataType struct {
	Name    string      `json:"name,omitempty"`
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

	err = m.validate()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (dt *DataType) validate() error {
	//	if dt.Name == "" {
	//		return fmt.Errorf("validation failed: Name is not set")
	//	}
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
		if dt.Name == "" {
			return fmt.Sprintf("<%s>", dt.Type)
		}
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
