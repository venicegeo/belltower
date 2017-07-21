/* Copyright 2017, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package engine

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
