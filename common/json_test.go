package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJson(t *testing.T) {
	assert := assert.New(t)

	type S struct{ X float64 }
	sObj := S{X: 3.0}
	sString := `{"X":3}`
	sMap := map[string]interface{}{
		"X": 3.0,
	}

	check := func(x *Json) {
		var err error
		assert.NoError(err)
		assert.EqualValues(sMap["X"], x.AsMap()["X"])
		assert.EqualValues(sString, x.AsString())
		xs := &S{}
		err = x.ToObject(xs)
		assert.NoError(err)
		assert.EqualValues(sObj.X, xs.X)
	}

	j1, err := NewJsonFromObject(&sObj)
	assert.NoError(err)
	check(j1)
	j2, err := NewJsonFromString(sString)
	assert.NoError(err)
	check(j2)
	j3, err := NewJsonFromMap(sMap)
	assert.NoError(err)
	check(j3)
}

func TestValidateJson(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(ValidateJsonString("{}"))
	assert.NoError(ValidateJsonString(`{"a":1}`))

	assert.Error(ValidateJsonString("{"))
	assert.Error(ValidateJsonString(`{a:1}`))
}
