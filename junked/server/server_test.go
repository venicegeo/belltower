package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOne(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	err := Server()
	assert.NoError(err)
}
