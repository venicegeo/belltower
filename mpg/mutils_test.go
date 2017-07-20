package mutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestX(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("mutils.go:10", SourceFile(0))
	assert.Equal("mutils_test.go:13", SourceFile(1))
}
