package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// this is a DeclBlock, containing both struct decls and string decls
var typeTestString = `
{ 
	"Main": {
        "MyInt": "int",
        "MyFloat": "float",
        "MyBool": "bool",
        "MyString": "string",

        "MyMapInt": "[map]int",
        "MyMapPoint": "[map]Point",
        "MySliceInt": "[]int",
        "MySlicePoint": "[]Point",
        "MyArray10Float": "[10]float",
        "MyArray4Point": "[4]Point"
    },

    "Point": {
        "x": "float",
        "y": "float"
    }, 

    "MyStruct": {
        "alpha": "string",
        "beta": "Point",
        "gamma": "MyStruct"
    }
}`

type typeTestStructItem struct {
	structName StructName
	node       TypeNode
	fields     []typeTestFieldItem
}

type typeTestFieldItem struct {
	fieldName FieldName
	node      TypeNode
}

var typeTestData = []typeTestStructItem{
	typeTestStructItem{
		structName: "Main",
		node: &TypeNodeStruct{
			Fields: map[FieldName]TypeNode{
				"MyInt":          NewTypeNodeInt(),
				"MyFloat":        NewTypeNodeFloat(),
				"MyString":       NewTypeNodeString(),
				"MyBool":         NewTypeNodeBool(),
				"MyMapInt":       NewTypeNodeMap(NewTypeNodeString(), NewTypeNodeInt()),
				"MyMapPoint":     NewTypeNodeMap(NewTypeNodeString(), NewTypeNodeName("Point")),
				"MySliceInt":     NewTypeNodeSlice(NewTypeNodeInt()),
				"MySlicePoint":   NewTypeNodeSlice(NewTypeNodeName("Point")),
				"MyArray10Float": NewTypeNodeArray(NewTypeNodeFloat(), 10),
				"MyArray4Point":  NewTypeNodeArray(NewTypeNodeName("Point"), 4),
			},
		},
		fields: []typeTestFieldItem{
			typeTestFieldItem{
				fieldName: "MyInt",
				node:      NewTypeNodeInt(),
			},
			typeTestFieldItem{
				fieldName: "MyFloat",
				node:      NewTypeNodeFloat(),
			},
			typeTestFieldItem{
				fieldName: "MyBool",
				node:      NewTypeNodeBool(),
			},
			typeTestFieldItem{
				fieldName: "MyString",
				node:      NewTypeNodeString(),
			},
			typeTestFieldItem{
				fieldName: "MyMapInt",
				node:      NewTypeNodeMap(NewTypeNodeString(), NewTypeNodeInt()),
			},
			typeTestFieldItem{
				fieldName: "MyMapPoint",
				node:      NewTypeNodeMap(NewTypeNodeString(), NewTypeNodeName("Point")),
			},
			typeTestFieldItem{
				fieldName: "MySliceInt",
				node:      NewTypeNodeSlice(NewTypeNodeInt()),
			},
			typeTestFieldItem{
				fieldName: "MySlicePoint",
				node:      NewTypeNodeSlice(NewTypeNodeName("Point")),
			},
			typeTestFieldItem{
				fieldName: "MyArray10Float",
				node:      NewTypeNodeArray(NewTypeNodeFloat(), 10),
			},
			typeTestFieldItem{
				fieldName: "MyArray4Point",
				node:      NewTypeNodeArray(NewTypeNodeName("Point"), 4),
			},
		},
	},
	typeTestStructItem{
		structName: "Point",
		node: &TypeNodeStruct{
			Fields: map[FieldName]TypeNode{
				"x": NewTypeNodeFloat(),
				"y": NewTypeNodeFloat(),
			},
		},
		fields: []typeTestFieldItem{
			typeTestFieldItem{
				fieldName: "x",
				node:      NewTypeNodeFloat(),
			},
			typeTestFieldItem{
				fieldName: "y",
				node:      NewTypeNodeFloat(),
			},
		},
	},
	typeTestStructItem{
		structName: "MyStruct",
		node: &TypeNodeStruct{
			Fields: map[FieldName]TypeNode{
				"alpha": NewTypeNodeString(),
				"beta":  NewTypeNodeName("Point"),
				"gamma": NewTypeNodeName("MyStruct"),
			},
		},
		fields: []typeTestFieldItem{
			typeTestFieldItem{
				fieldName: "alpha",
				node:      NewTypeNodeString(),
			},
			typeTestFieldItem{
				fieldName: "beta",
				node:      NewTypeNodeName("Point"),
			},
			typeTestFieldItem{
				fieldName: "gamma",
				node:      NewTypeNodeName("MyStruct"),
			},
		},
	},
}

func TestTypeTokenizer(t *testing.T) {
	assert := assert.New(t)

	tt, err := NewTypeTokenizer()
	assert.NoError(err)
	typeTable, err := tt.ParseJson(typeTestString)
	assert.NoError(err)
	assert.NotNil(typeTable)

	for _, sItem := range typeTestData {

		se, ok := typeTable.Structs[sItem.structName]
		assert.True(ok)
		assert.Equal(sItem.structName, se.Name)

		assert.EqualValues(sItem.node.String(), se.Type.String())

		for _, fItem := range se.Fields {

			fe, ok := se.Fields[fItem.Name]
			assert.True(ok)
			assert.Equal(fItem.Name, fe.Name)
			assert.Equal(fItem.Type, fe.Type)

		}
	}
}
