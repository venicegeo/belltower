package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//---------------------------------------------------------------------------

func TestEnvironment(t *testing.T) {
	assert := assert.New(t)

	env := NewEnvironment()

	env.set("i", ExprValue{Type: IntType, Value: 12})
	assert.Equal(ExprValue{Type: IntType, Value: 12}, env.get("i"))
}
