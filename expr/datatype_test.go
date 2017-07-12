package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			jsn: `{
				"x": "array[int]"
			}`,
			expected: "<struct: <array: <int>>>",
		},
		{
			jsn: `{
				"x": "map[bool]"
			}`,
			expected: "<struct: <map: <bool>>>",
		},
		{
			jsn: `{
				"x": {
					"a": "bool",
					"b": "string"
					}
			}`,
			expected: "<struct: <struct: <bool>, <string>>>",
		},
	}

	for _, item := range items {
		i, err := NewDataTypeFromJSON(item.jsn)
		assert.NoError(err, item.expected)
		assert.Equal(item.expected, i.String(), item.expected)
	}
}
