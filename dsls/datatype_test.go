package dsls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeEnum(t *testing.T) {
	assert := assert.New(t)

	items := []struct {
		src string
		ok  bool
	}{
		{src: "int", ok: true},
		{src: "float", ok: true},
		{src: "bool", ok: true},
		{src: "string", ok: true},
		{src: "map", ok: true},
		{src: "array", ok: true},
		{src: "struct", ok: true},
		{src: "", ok: false},
		{src: "record", ok: false},
	}

	for _, item := range items {
		t := TypeNameFromString(item.src)
		if item.ok {
			assert.NotEqual(TypeNameInvalid, t)
			s := t.String()
			assert.Equal(item.src, s)
		} else {
			assert.Equal(TypeNameInvalid, t)
		}
	}
}

func TestDataType(t *testing.T) {
	assert := assert.New(t)

	i := NewScalarDataType("ii", TypeNameInt)
	assert.Equal("<ii/int>", i.String())

	f := NewScalarDataType("ff", TypeNameFloat)
	assert.Equal("<ff/float>", f.String())

	s := NewScalarDataType("ss", TypeNameString)
	assert.Equal("<ss/string>", s.String())

	b := NewScalarDataType("bb", TypeNameBool)
	assert.Equal("<bb/bool>", b.String())

	a := NewArrayDataType("v1", b)
	assert.Equal("<v1/array: <bb/bool>>", a.String())

	m := NewMapDataType("v2", a)
	assert.Equal("<v2/map: <v1/array: <bb/bool>>>", m.String())

	fs := []*DataType{i, s}
	st := NewStructDataType("v3", fs)
	assert.Equal("<v3/struct: <ii/int>, <ss/string>>", st.String())
}

func TestDataTypeJSON(t *testing.T) {
	assert := assert.New(t)

	items := []struct {
		jsn      string
		expected string
	}{
		{
			jsn: `
			{
				"name": "IN",
				"type": "struct",
				"fields": [
					{
						"name": "x",
						"type": "int"
					},
					{
						"name": "y",
						"type": "float"
					}
				]
			}`,
			expected: "<IN/struct: <x/int>, <y/float>>",
		},
		{
			jsn: `
			{
				"name": "in",
				"type": "struct",
				"fields": [
					{
						"name": "x",
						"type": "array",
						"element": {
							"type": "int"
						}
					}
				]
			}`,
			expected: "<in/struct: <x/array: <int>>>",
		},
		{
			jsn: `
			{
				"name": "x",
				"type": "struct",
				"fields": [
					{
						"name": "xx",
						"type": "map",
						"element": {
							"type": "bool"
						}
					}
				]
			}`,
			expected: "<x/struct: <xx/map: <bool>>>",
		},
		{
			jsn: `{
				"name": "x",
				"type": "struct",
				"fields": [
					{
						"name": "y",
						"type": "struct",
						"fields": [
							{
								"name": "a",
								"type": "bool"
							},
							{
								"name": "b",
								"type": "string"
							}
						]
					}
				]
			}`,
			expected: "<x/struct: <y/struct: <a/bool>, <b/string>>>",
		},
	}

	for _, item := range items {
		i, err := NewDataTypeFromJSON(item.jsn)
		assert.NoError(err, item.expected)
		assert.Equal(item.expected, i.String(), item.expected)
	}
}
