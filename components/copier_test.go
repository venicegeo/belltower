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

func TestCopier(t *testing.T) {
	assert := assert.New(t)

	config := engine.ArgMap{}
	copierX, err := engine.Factory.Create("Copier", config)
	assert.NoError(err)
	copier := copierX.(*Copier)

	// this setup is normally done by goflow itself
	chIn := make(chan string)
	chOut1 := make(chan string)
	chOut2 := make(chan string)
	copier.Input = chIn
	copier.Output1 = chOut1
	copier.Output2 = chOut2

	inJ := `{
		"alpha": 11.0,
		"beta":  22.0,
		"gamma": 33.0
	}`
	go copier.OnInput(inJ)

	// check the returned result
	out1J := <-chOut1
	out2J := <-chOut2

	assert.JSONEq(inJ, out1J)
	assert.JSONEq(inJ, out2J)
}
