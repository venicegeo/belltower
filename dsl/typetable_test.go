package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//--------------------------

func TestTypeTable(t *testing.T) {
	assert := assert.New(t)
	var err error

	tt, err := NewTypeTable()
	assert.NoError(err)
	assert.Equal(0, len(tt.Structs))

	// add a symbol
	err = tt.addStruct("myint", NewTypeNodeInt())
	assert.NoError(err)
	assert.Equal(1, len(tt.Structs))

	// add a symbol again
	err = tt.addStruct("myint", NewTypeNodeInt())
	assert.Error(err)
	assert.Equal(1, len(tt.Structs))

	// test has()
	assert.True(tt.hasStruct("myint"))
	assert.False(tt.hasStruct("foofoofoo"))
}
