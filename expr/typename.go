package expr

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

func FromString(s string) TypeName {
	if isValidTypeName(s) {
		return TypeName(s)
	}
	return TypeNameInvalid
}
