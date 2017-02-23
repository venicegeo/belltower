package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {

	assert := assert.New(t)

	testData := []struct {
		node   ExprNode
		values map[string]ExprValue
		typs   map[string]TypeNode
		result ExprValue
	}{
		{ // "a * (b + c )"
			node: NewExprNodeMultiply(
				NewExprNodeAdd(NewExprNodeSymbolRef("c"), NewExprNodeSymbolRef("b")),
				NewExprNodeSymbolRef("a")),
			values: map[string]ExprValue{
				"a": ExprValue{Type: IntType, Value: 2},
				"b": ExprValue{Type: IntType, Value: 3},
				"c": ExprValue{Type: IntType, Value: 4},
			},
			typs: map[string]TypeNode{
				"a": NewTypeNodeInt(),
				"b": NewTypeNodeInt(),
				"c": NewTypeNodeInt(),
			},
			result: ExprValue{Type: IntType, Value: 14},
		},
	}

	for _, item := range testData {

		typeTable, err := NewTypeTable()
		assert.NoError(err)

		env := NewEnvironment()

		err = typeTable.addStruct("x", nil)
		assert.NoError(err)
		for k, v := range item.typs {
			err = typeTable.addField("x", FieldName(k), v)
			assert.NoError(err)
		}

		for k, v := range item.values {
			env.set(k, v)
		}

		eval := &Eval{}
		result, err := eval.Evaluate(item.node, env)
		assert.NoError(err)
		assert.NotNil(result)
		assert.EqualValues(item.result, result)
	}

}
