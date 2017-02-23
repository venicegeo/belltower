package dsl

import (
	"encoding/json"
	"fmt"
)

type TypeTokenizer struct {
	typeTable *TypeTable
}

func NewTypeTokenizer() (*TypeTokenizer, error) {
	typeTable, err := NewTypeTable()
	if err != nil {
		return nil, err
	}

	tp := &TypeTokenizer{
		typeTable: typeTable,
	}

	return tp, nil
}

// DeclBlock is a JSON string declares the types used in the event.
//
// A DeclBlock is a map from type names to struct type definitions.
//
// A struct type definition is a map from field names to simple type definitions.
//
// A simple type definition is either a base type or a base type collection.
// The base types are int, float, string, and bool. A base type collection
// is a map, an array of fixed size, or an array of arbitrary size.
//
// The DeclBlock must contain a struct type definition with the name "Main".
// By convention, type names are capitalized, but field names are not.
//
// Here is a DeclBlock that uses struct types:
//
// {
//     "Height": {
//         "h":        "int",
//         "inMeters": "bool",
//     },
//     "Point3": {
//         "x": "float",
//         "y": "float",
//         "z": "Height",
//     },
//     "Main": {
//	       "ul":        "Point3",
//	       "lr":        "Point3",
//         "timestamp": "string",
//     }
// }
//
// Here is a DeclBlock that contains only two struct types, one of which uses
// a variety of base types and base type collections:
//
// {
//     "Point2": {
//         "x": "float",
//         "y": "float",
//     },
//     "Main": {
//	       "i":       "int",
//         "f":       "float",
//         "s":       "string",
//         "b":       "bool",
//	       "f4":      "[4]float", // array of four floats
//         "iN":      "[]int", // array of 0 to N ints
//         "bmap":    "[map]bool", // map of string to bool
//         "pmapmap": "[map][map]Point2", // map of (string) -> (map of (string) -> Point2)
//     }
// }

type StructName string
type FieldName string
type FieldDecl string
type DeclBlock map[StructName]*StructDecl
type StructDecl map[FieldName]*FieldDecl

// ParseJson takes a declaration block expressed as JSON string and parses it.
func (p *TypeTokenizer) ParseJson(s string) (*TypeTable, error) {
	declBlock := &DeclBlock{}
	err := json.Unmarshal([]byte(s), declBlock)
	if err != nil {
		return nil, err
	}

	return p.Parse(declBlock)
}

// Parse takes a declaration block expressed as a DeclBlock object and parses it.
func (p *TypeTokenizer) Parse(block *DeclBlock) (*TypeTable, error) {

	for structName, structDecl := range *block {

		_, err := p.parseStruct(structName, structDecl)
		if err != nil {
			return nil, err
		}

	}

	return p.typeTable, nil
}

func (p *TypeTokenizer) parseStruct(structName StructName, structDecl *StructDecl) (TypeNode, error) {
	var err error
	var tnode TypeNode

	structNode := NewTypeNodeStruct()

	err = p.typeTable.addStruct(structName, structNode)
	if err != nil {
		return nil, err
	}

	for fieldName, fieldDecl := range *structDecl {

		tnode, err = p.parseField(structName, fieldName, fieldDecl)
		if err != nil {
			return nil, err
		}
		err = p.typeTable.addField(structName, fieldName, tnode)
		if err != nil {
			return nil, err
		}
		structNode.Fields[fieldName] = NewTypeNodeField(fieldName, tnode)
	}

	return structNode, nil
}

func (p *TypeTokenizer) parseField(structName StructName, fieldName FieldName, fieldDecl *FieldDecl) (TypeNode, error) {
	var err error
	var tnode TypeNode

	tnode, err = p.parseFieldDecl(fieldDecl)
	if err != nil {
		return nil, err
	}
	//log.Printf("$$ %s $$ %v", name, tnode)

	return tnode, nil
}

func (p *TypeTokenizer) parseFieldDecl(decl *FieldDecl) (TypeNode, error) {

	scanner := Scanner{}

	s := string(*decl)

	toks, err := scanner.Scan(s)
	if err != nil {
		return nil, err
	}

	tnode, err := parseTheTokens(toks, p.typeTable)
	if err != nil {
		return nil, err
	}

	//log.Printf("$$ %s $$ %s $$ %v", name, stringDecl, tnode)
	return tnode, nil
}

//---------------------------------------------------------------------------

func parseTheTokens(toks []Token, typeTable *TypeTable) (TypeNode, error) {

	t0 := toks[0]
	t1ok := len(toks) > 1
	//t2ok := len(toks) > 2

	var out TypeNode

	switch t0.Id {

	case TokenSymbol:
		if t1ok {
			return nil, fmt.Errorf("extra token after %v\n\t%v", t0, toks[1])
		}
		switch t0.Text {
		case "int":
			out = NewTypeNodeInt()
		case "float":
			out = NewTypeNodeFloat()
		case "string":
			out = NewTypeNodeString()
		case "bool":
			out = NewTypeNodeBool()
		default:
			out = NewTypeNodeName(t0.Text)
		}

	case TokenTypeSlice:
		if !t1ok {
			return nil, fmt.Errorf("no token after %v", t0)
		}
		next, err := parseTheTokens(toks[1:], typeTable)
		if err != nil {
			return nil, err
		}
		out = NewTypeNodeSlice(next)

	case TokenTypeMap:
		if !t1ok {
			return nil, fmt.Errorf("no token after %v", t0)
		}
		next, err := parseTheTokens(toks[1:], typeTable)
		if err != nil {
			return nil, err
		}
		out = NewTypeNodeMap(NewTypeNodeString(), next)

	case TokenTypeArray:
		if !t1ok {
			return nil, fmt.Errorf("no token after %v", t0)
		}
		next, err := parseTheTokens(toks[1:], typeTable)
		if err != nil {
			return nil, err
		}
		out = NewTypeNodeArray(next, t0.Value.(int))

	default:
		return nil, fmt.Errorf("unhandled token: " + t0.String())
	}

	return out, nil
}
