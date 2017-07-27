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

func TestNotifier(t *testing.T) {
	assert := assert.New(t)

	config := engine.ArgMap{
		"path": "/tmp/bf/outputs",
	}
	notifierX, err := engine.Factory.Create("Notifier", config)
	assert.NoError(err)
	notifier := notifierX.(*Notifier)

	// this setup is normally done by goflow itself
	chIn := make(chan string)
	chOut := make(chan string)
	notifier.Input = chIn
	notifier.Output = chOut

	inJ := `{"SelectedImage": "LC81340472017028LGN00"}`
	go notifier.OnInput(inJ)

	// check the returned result
	outJ := <-chOut
	assert.True(outJ != "")
}
