package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDsl(t *testing.T) {
	assert := assert.New(t)

	data := []struct {
		decl   string
		expr   string
		env    *Environment
		result ExprValue
	}{
		{
			decl: `{ "Point": { "x": "float", "y": "float" } }`,
			expr: "a * (b + c )",
			env: &Environment{
				data: map[string]ExprValue{
					"a": ExprValue{Type: IntType, Value: 2},
					"b": ExprValue{Type: IntType, Value: 3},
					"c": ExprValue{Type: IntType, Value: 4},
				},
			},
			result: ExprValue{Type: IntType, Value: 14},
		},
	}

	for _, item := range data {
		d, err := NewDsl()
		assert.NoError(err)

		tId, err := d.ParseDeclaration(item.decl)
		assert.NoError(err)
		assert.NotEqual(InvalidId, tId)

		eId, err := d.ParseExpression(item.expr)
		assert.NoError(err)
		assert.NotEqual(InvalidId, eId)

		result, err := d.Evaluate(eId, tId, item.env)
		assert.NoError(err)
		assert.NotNil(result)
		assert.EqualValues(item.result, result)
	}
}
