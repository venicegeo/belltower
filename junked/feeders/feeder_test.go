package feeders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	assert := assert.New(t)

	// check registry
	assert.True(len(feederRegistry.data) > 0)
}
