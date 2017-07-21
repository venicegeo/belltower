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
package merr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerr(t *testing.T) {
	assert := assert.New(t)

	saved := UseSourceStamp

	e := New("la dee dah")
	assert.Equal("la dee dah", e.Message)
	assert.Equal("merr_test.go:28", e.Source)
	assert.Equal("merr_test.go:28 la dee dah", e.Error())

	e = Newf("la dee doh")
	assert.Equal("la dee doh", e.Message)
	assert.Equal("merr_test.go:33", e.Source)
	assert.Equal("merr_test.go:33 la dee doh", e.Error())

	UseSourceStamp = false
	assert.Equal("la dee doh", e.Error())

	UseSourceStamp = saved
}
