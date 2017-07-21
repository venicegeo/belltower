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
package mutils

import "testing"
import "fmt"
import "github.com/stretchr/testify/assert"

func TestSerializer(t *testing.T) {

	// good
	{
		jsn := `{ "X": 7 }`

		a := &A{}
		err := a.ReadFromJSON(jsn)
		assert.NoError(t, err)

		err = a.Validate()
		assert.NoError(t, err)
	}

	// bad
	{
		jsn := `{ "Y": 11 }`

		a := &A{}
		err := a.ReadFromJSON(jsn)
		assert.NoError(t, err)

		err = a.Validate()
		assert.Error(t, err)
		assert.Equal(t, "validation error", err.Error())
	}
}

//---------------------------------------------------------------------

type A struct {
	X int
}

func (a *A) WriteToJSON() (string, error) { return WriteToJSON(a) }

func (a *A) ReadFromJSON(s string) error { return ReadFromJSON(s, a) }

func (a *A) Validate() error {
	if a.X == 0 {
		return fmt.Errorf("validation error")
	}
	return nil
}
