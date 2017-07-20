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
package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/engine"
)

func TestRemapper(t *testing.T) {
	assert := assert.New(t)

	config := engine.ArgMap{
		"remaps": map[string]string{
			"alpha": "omega",
			"beta":  "psi",
		},
	}
	remapperX, err := engine.Factory.Create("Remapper", config)
	assert.NoError(err)
	remapper := remapperX.(*Remapper)

	// this setup is normally done by goflow itself
	chIn := make(chan string)
	chOut := make(chan string)
	remapper.Input = chIn
	remapper.Output = chOut

	inJ := `{
		"alpha": 11.0,
		"beta":  22.0,
		"gamma": 33.0
	}`
	go remapper.OnInput(inJ)

	// check the returned result
	outJ := <-chOut
	outM, err := engine.NewArgMap(outJ)
	assert.NoError(err)

	assert.Equal(11.0, outM["omega"])
	assert.Equal(22.0, outM["psi"])
	assert.Equal(33.0, outM["gamma"])
	assert.NotContains(outM, "alpha")
	assert.NotContains(outM, "beta")
}
