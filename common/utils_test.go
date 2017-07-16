package common

import (
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

// derived from the mapstructure examples
func TestMapStructure(t *testing.T) {

	assert := assert.New(t)

	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
		Dur    time.Duration
	}

	input := map[string]interface{}{
		"name":   "Mitchell",
		"age":    "91",
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"twitter": "mitchellh",
		},
		"dur":      "1h",
		"frobnitz": "xwgx4w5v3",
	}

	var result Person
	var md mapstructure.Metadata
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Metadata:         &md,
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		Result:           &result,
	}

	decoder, err := mapstructure.NewDecoder(config)
	assert.NoError(err)

	{
		err = decoder.Decode(input)
		assert.NoError(err)
		assert.Len(md.Unused, 1)
		assert.Equal("frobnitz", md.Unused[0])
		assert.Equal(91, result.Age)
	}

	input["age"] = "foo"
	{
		err = decoder.Decode(input)
		assert.Error(err)
		assert.Contains(err.Error(), `parsing "foo": invalid syntax`)
	}
}

func TestSetStructFromMap(t *testing.T) {

	assert := assert.New(t)

	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
		Dur    time.Duration
	}

	input := map[string]interface{}{
		"name":   "Mitchell",
		"age":    "91",
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"twitter": "mitchellh",
		},
		"dur":      "1h",
		"frobnitz": "xwgx4w5v3",
	}

	{
		var result Person
		result2, err := SetStructFromMap(input, &result, true)
		assert.NoError(err)
		assert.Equal(91, result2.(*Person).Age)
	}

	input["age"] = "foo"

	{
		var result Person
		_, err := SetStructFromMap(input, &result, true)
		assert.Error(err)
		assert.Contains(err.Error(), `parsing "foo": invalid syntax`)
	}
}
