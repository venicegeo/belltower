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
	"io/ioutil"
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/engine"
)

func TestLogger(t *testing.T) {
	assert := assert.New(t)

	// TODO: add tests to make sure stdout and stderr output targets work
	const file = "testlogger.log" // or "STDOUT" or "STDERR"

	err := os.Remove(file)
	assert.NoError(err)

	config := engine.ArgMap{
		"FileName": file,
	}
	loggerX, err := engine.Factory.Create("Logger", config)
	assert.NoError(err)
	logger := loggerX.(*Logger)

	// this setup is normally done by goflow itself
	chIn := make(chan string)
	chOut := make(chan string)
	logger.Input = chIn
	logger.Output = chOut

	lines := []string{
		`{ "Alpha" : 1 }`,
		`{ "Beta"  : 2 }`,
		`{ "Gamma" : 3 }`,
	}
	go logger.OnInput(lines[0])
	go logger.OnInput(lines[1])
	go logger.OnInput(lines[2])

	// ignore the returned result
	_ = <-chOut
	_ = <-chOut
	_ = <-chOut

	// verify the log file
	result, err := ioutil.ReadFile(file)
	assert.NoError(err)

	linesC := []string{
		`{"Alpha":1}`,
		`{"Beta":2}`,
		`{"Gamma":3}`,
	}

	assert.Contains(string(result), linesC[0])
	assert.Contains(string(result), linesC[1])
	assert.Contains(string(result), linesC[2])

}
