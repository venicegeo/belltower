package dsl

import "fmt"

type TypeTableFieldEntry struct {
	Name FieldName
	Type TypeNode
}

type TypeTableStructEntry struct {
	Name   StructName
	Type   TypeNode
	Fields map[FieldName]*TypeTableFieldEntry
}

type TypeTable struct {
	Structs map[StructName]*TypeTableStructEntry
}

func NewTypeTable() (*TypeTable, error) {

	tt := &TypeTable{
		Structs: map[StructName]*TypeTableStructEntry{},
	}

	return tt, nil
}

//---------------------------------------------------------------------------

type BaseType int

const (
	IntType BaseType = iota
	FloatType
	BoolType
	StringType
)

//---------------------------------------------------------------------------

func (e *TypeTableFieldEntry) String() string {
	s := fmt.Sprintf("%s:", e.Name)
	s += fmt.Sprintf("[%v]", e.Type)
	return s
}

func (e *TypeTableStructEntry) String() string {
	s := fmt.Sprintf("%s:", e.Name)
	s += fmt.Sprintf("[%v]", e.Type)
	for _, f := range e.Fields {
		s += fmt.Sprintf("  %v", f)
	}
	return s
}

func (st *TypeTable) String() string {
	s := ""
	for _, tte := range st.Structs {
		s += fmt.Sprintf("  %v\n", tte)
	}
	return s
}

func (st *TypeTable) addStruct(name StructName, node TypeNode) error {
	if st.hasStruct(name) {
		return fmt.Errorf("type table struct entry already exists: %s", name)
	}
	st.Structs[name] = &TypeTableStructEntry{
		Name:   name,
		Type:   node,
		Fields: map[FieldName]*TypeTableFieldEntry{},
	}
	return nil
}

func (st *TypeTable) getStruct(s StructName) TypeNode {
	v, ok := st.Structs[s]
	if !ok {
		return nil
	}
	return v.Type
}

func (st *TypeTable) hasStruct(s StructName) bool {
	_, ok := st.Structs[s]
	return ok
}

func (st *TypeTable) addField(sn StructName, fn FieldName, node TypeNode) error {
	se, ok := st.Structs[sn]
	if !ok {
		return fmt.Errorf("type table field struct entry does not exist: %s", sn)
	}
	if st.hasField(sn, fn) {
		return fmt.Errorf("type table field entry already exists: %s.%s", sn, fn)
	}
	se.Fields[fn] = &TypeTableFieldEntry{
		Name: fn,
		Type: node,
	}
	return nil
}

func (st *TypeTable) getField(sn StructName, fn FieldName) TypeNode {
	se, ok := st.Structs[sn]
	if !ok {
		return nil
	}
	fe, ok := se.Fields[fn]
	if !ok {
		return nil
	}
	return fe.Type
}

func (st *TypeTable) hasField(sn StructName, fn FieldName) bool {
	se, ok := st.Structs[sn]
	if !ok {
		return false
	}
	_, ok = se.Fields[fn]
	return ok
}
