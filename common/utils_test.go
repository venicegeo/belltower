package common

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

	assert.True(m.Has("i"))
	assert.False(m.Has("y"))

	assert.True(m.IsInt("i"))
	assert.False(m.IsInt("s"))
	assert.Equal(5, m.AsInt("i"))
	assert.Equal(0, m.AsInt("s"))

	assert.True(m.IsString("s"))
	assert.False(m.IsString("i"))
	assert.Equal("foo", m.AsString("s"))
	assert.Equal("", m.AsString("i"))
}
