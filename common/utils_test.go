package common

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMapValue(t *testing.T) {
	assert := assert.New(t)

	m := map[string]interface{}{
		"b":   true,
		"bx":  "asdf",
		"i":   17,
		"ix":  "asdf",
		"f":   0.5,
		"fx":  true,
		"s":   "asdf",
		"sx":  19,
		"d":   "5s",
		"dx":  "asdf",
		"dxx": 17.19,
	}

	b, err := GetMapValueAsBool(m, "b")
	assert.NoError(err)
	assert.Equal(b, true)
	_, err = GetMapValueAsBool(m, "bx")
	assert.Error(err)
	_, err = GetMapValueAsBool(m, "bb")
	assert.Error(err)

	i, err := GetMapValueAsInt(m, "i")
	assert.NoError(err)
	assert.Equal(i, 17)
	_, err = GetMapValueAsInt(m, "ix")
	assert.Error(err)
	_, err = GetMapValueAsInt(m, "ii")
	assert.Error(err)

	f, err := GetMapValueAsFloat(m, "f")
	assert.NoError(err)
	assert.Equal(f, 0.5)
	_, err = GetMapValueAsFloat(m, "fx")
	assert.Error(err)
	_, err = GetMapValueAsFloat(m, "ff")
	assert.Error(err)

	s, err := GetMapValueAsString(m, "s")
	assert.NoError(err)
	assert.Equal(s, "asdf")
	_, err = GetMapValueAsString(m, "sx")
	assert.Error(err)
	_, err = GetMapValueAsString(m, "ss")
	assert.Error(err)

	d, err := GetMapValueAsDuration(m, "d")
	assert.NoError(err)
	assert.Equal(d, time.Duration(5*time.Second))
	_, err = GetMapValueAsDuration(m, "dx")
	assert.Error(err)
	_, err = GetMapValueAsDuration(m, "dxx")
	assert.Error(err)
	_, err = GetMapValueAsDuration(m, "dd")
	assert.Error(err)
}
