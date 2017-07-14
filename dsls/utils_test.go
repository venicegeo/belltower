package dsls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert := assert.New(t)

	m := Map{
		"i": 5,
		"s": "foo",
		"x": false,
	}

	assert.True(m.has("i"))
	assert.False(m.has("y"))

	assert.True(m.isInt("i"))
	assert.False(m.isInt("s"))
	assert.Equal(5, m.asInt("i"))
	assert.Equal(0, m.asInt("s"))

	assert.True(m.isString("s"))
	assert.False(m.isString("i"))
	assert.Equal("foo", m.asString("s"))
	assert.Equal("", m.asString("i"))
}
