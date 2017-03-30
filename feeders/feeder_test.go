package feeders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeed(t *testing.T) {
	assert := assert.New(t)

	// check registry
	assert.True(len(feederFactory.factories) > 0)

	// (real testing of Feeder system done in the derived classes, e.g. SimpleFeeder)
}
