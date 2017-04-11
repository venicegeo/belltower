package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdent(t *testing.T) {
	assert := assert.New(t)

	assert.True(NoIdent.String() == "")

	id := NewId()
	assert.False(id.String() == "")

	assert.EqualValues(id, ToIdent(id.String()))
}
