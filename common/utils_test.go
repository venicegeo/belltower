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
		"d":   time.Duration(time.Second * 5),
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

func TestObjectEquality(t *testing.T) {
	assert := assert.New(t)

	// TODO: ObjectsAreEqual is broken, and probably nondeterministic

	a := map[string]string{
		"A": "a",
		"B": "b",
	}
	b := map[string]string{
		"B": "b",
		"A": "a",
	}
	c := map[string]string{
		"A": "a",
	}
	d := map[string]string{
		"B": "b",
	}
	var bb interface{} = b

	assert.True(ObjectsAreEqual(a, b))
	assert.True(ObjectsAreEqualValues(a, b))
	assert.False(ObjectsAreEqual(a, c))
	assert.False(ObjectsAreEqualValues(a, c))
	assert.False(ObjectsAreEqual(a, d))
	assert.False(ObjectsAreEqualValues(a, d))
	assert.True(ObjectsAreEqual(b, bb))
	assert.True(ObjectsAreEqual(a, bb))
	assert.True(ObjectsAreEqualValues(b, bb))
	assert.True(ObjectsAreEqualValues(a, bb))
}

func TestMapObjectEquality(t *testing.T) {
	assert := assert.New(t)
	x := map[string]interface{}{
		"A": "a",
		"B": "b",
	}
	y := map[string]interface{}{
		"B": "b",
		"A": "a",
	}
	z1 := map[string]interface{}{
		"Z": "a",
		"B": "b",
	}
	z2 := map[string]interface{}{
		"A": "z",
		"B": "b",
	}

	assert.True(MapsAreEqualValues(x, x))
	assert.True(MapsAreEqualValues(x, y))
	assert.True(MapsAreEqualValues(y, x))
	assert.True(MapsAreEqualValues(y, y))

	assert.False(MapsAreEqualValues(x, z1))
	assert.False(MapsAreEqualValues(x, z2))
	assert.False(MapsAreEqualValues(z1, z2))

}
