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

	// positive tests
	i, err := m.GetFloat("i")
	assert.NoError(err)
	assert.Equal(float64(5), i)
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
