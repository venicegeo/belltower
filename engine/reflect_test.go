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

func TestReflection(t *testing.T) {
	assert := assert.New(t)

	type T struct {
		I int
		S string
	}

	{
		i := 9

		x := NewViaReflection(i)
		ii := x.(int)
		assert.Equal(0, ii)

		ii = 4
		assert.Equal(4, ii)
		assert.Equal(9, i)
	}
	{
		i := 9

		x := NewViaReflection(&i)
		ii := x.(*int)
		assert.Equal(0, *ii)

		*ii = 4
		assert.Equal(4, *ii)
		assert.Equal(9, i)
	}

	{
		a := T{}
		aa := NewViaReflection(a)

		aaa, ok := aa.(T)
		assert.True(ok)
		assert.Equal(aaa.S, "")
		assert.Equal(aaa.I, 0)
		aaa.S = "sss"
		aaa.I = 99
		assert.Equal(aaa.S, "sss")
		assert.Equal(aaa.I, 99)
	}
	{
		a := &T{}
		aa := NewViaReflection(a)

		aaa, ok := aa.(*T)
		assert.True(ok)
		assert.Equal(aaa.S, "")
		assert.Equal(aaa.I, 0)
		aaa.S = "sss"
		aaa.I = 99
		assert.Equal(aaa.S, "sss")
		assert.Equal(aaa.I, 99)
	}
	//log.Printf("%T %#v", aa, aa)
}
