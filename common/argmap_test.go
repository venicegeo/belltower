package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgMap(t *testing.T) {
	assert := assert.New(t)

	m := ArgMap{
		"i": 5,
		"s": "foo",
		"b": true,
		"f": 12.34,
	}

	// positiove tests
	i, err := m.GetInt("i")
	assert.NoError(err)
	assert.Equal(5, i)
	f, err := m.GetFloat("f")
	assert.NoError(err)
	assert.Equal(12.34, f)
	s, err := m.GetString("s")
	assert.NoError(err)
	assert.Equal("foo", s)
	b, err := m.GetBool("b")
	assert.NoError(err)
	assert.Equal(true, b)

	// negative tests
	_, err = m.GetInt("s")
	assert.Error(err)
	_, err = m.GetInt("x")
	assert.Error(err)
	_, err = m.GetFloat("s")
	assert.Error(err)
	_, err = m.GetFloat("x")
	assert.Error(err)
	_, err = m.GetString("i")
	assert.Error(err)
	_, err = m.GetString("x")
	assert.Error(err)
	_, err = m.GetBool("s")
	assert.Error(err)
	_, err = m.GetBool("x")
	assert.Error(err)
}
