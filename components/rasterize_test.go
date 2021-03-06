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

func TestRasterize(t *testing.T) {
	assert := assert.New(t)

	config := engine.ArgMap{}
	rasterizeX, err := engine.Factory.Create("Rasterize", config)
	assert.NoError(err)
	rasterize := rasterizeX.(*Rasterize)

	// this setup is normally done by goflow itself
	chIn := make(chan string)
	chOut := make(chan string)
	rasterize.Input = chIn
	rasterize.Output = chOut

	inJ := `{"SelectedImage": "LC81880282017198LGN00"}`
	go rasterize.OnInput(inJ)

	// check the returned result
	outJ := <-chOut
	assert.JSONEq(`{"SelectedImage":"LC81880282017198LGN00"}`, outJ)
}
