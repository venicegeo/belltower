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
package mlog

import (
	"testing"

	"os"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

func TestMLog(t *testing.T) {
	assert := assert.New(t)

	const filename = "testmlog.log"
	os.Remove(filename)
	SetFlags(0)
	{
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
		assert.NoError(err)
		defer f.Close()
		SetOutput(f)

		Print("one")
		Printf("%s", "two")
		Debug("three")
		Debugf("%s", "four")
		Verbose = true
		Debug("five")
		Debugf("%s", "six")
		SetFlags(Lshortfile)
		Print("seven")
	}

	expected :=
		`one
two
five
six
mlog_test.go:33 seven
`
	byts, err := ioutil.ReadFile(filename)
	assert.NoError(err)
	assert.Equal(expected, string(byts))

	os.Remove(filename)

}
