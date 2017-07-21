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

func TestFactory(t *testing.T) {
	assert := assert.New(t)

	Factory.Register("Xyzzy", &xyzzy{})

	{
		// make sure we can init without a config map (e.g. to allow defaults to work)
		x, err := Factory.Create("Xyzzy", nil)
		assert.NoError(err)
		assert.NotNil(x)

		assert.Equal(0, x.(*xyzzy).executionCount)
		assert.Nil(x.(*xyzzy).Config["param"])
		assert.Equal(19, x.(*xyzzy).myint)
	}

	config := ArgMap{
		"param": "seventeen",
	}
	x, err := Factory.Create("Xyzzy", config)
	assert.NoError(err)
	assert.NotNil(x)
	xx := x.(*xyzzy)

	assert.Equal(0, xx.executionCount)
	assert.Equal("seventeen", xx.Config["param"])
	assert.Equal(19, xx.myint)
}

//---------------------------------------------------------------------

type xyzzy struct {
	ComponentCore
	myint int
}

func (x *xyzzy) Configure() error {
	x.myint = 19
	return nil
}

func (x *xyzzy) Run(in interface{}) (interface{}, error) {
	return nil, nil
}
