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

func TestParseFails(t *testing.T) {
	assert := assert.New(t)

	texts := []string{
		"1+",
		"foo.bar + 3",
		"3 + x[2]",
		`x["key"]`,
	}

	for _, text := range texts {
		e, err := NewExpression(text, nil)
		assert.Error(err, text)
		assert.Nil(e, text)
	}
}

func TestEvalFails(t *testing.T) {
	assert := assert.New(t)

	//myarray := []float64{0.1, 1.2, 2.3, 3.4}
	//mymap := map[string]int{"two": 2}

	data := []struct {
		text string
		vars Variables
	}{
		{text: "1+x"},
		//{text: "x[2] * 2", env: NewEnvironment(map[string]interface{}{"x": myarray})},
		//{text: `2 + x["two"]`, env: NewEnvironment(map[string]interface{}{"x": mymap})},
	}

	for _, item := range data {
		e, err := NewExpression(item.text, nil)
		assert.NoError(err, item.text)
		assert.NotNil(e, item.text)

		x, err := e.Eval(item.vars)
		assert.Error(err, item.text)
		assert.Nil(x, item.text)
	}
}

func TestEvals(t *testing.T) {
	assert := assert.New(t)

	mystruct := struct {
		Frob int
	}{
		Frob: 19,
	}

	vars1 := map[string]interface{}{"x": 15.5}
	vars2 := map[string]interface{}{"x": mystruct}

	data := []struct {
		text     string
		expected interface{}
		vars     Variables
	}{
		// misc
		{expected: 3.0, text: "1+2"},

		// constants
		{expected: 17.0, text: "17"},

		// binops
		{expected: 3.0, text: "1+2"},
		{expected: -1.0, text: "1-2"},
		{expected: 6.0, text: "2*3"},
		{expected: 2.0, text: "6/3"},
		{expected: 1.2, text: "6/5"},
		{expected: true, text: "41>5 && (4/5 < 1.1)"},

		// with an environment
		{expected: 19.0, text: "x+3.5", vars: vars1},
		{expected: 22.5, text: "x.Frob + 3.5", vars: vars2},
	}

	for _, item := range data {
		e, err := NewExpression(item.text, nil)
		assert.NoError(err)
		assert.NotNil(e)

		x, err := e.Eval(item.vars)
		assert.NoError(err)
		assert.NotNil(x)

		assert.Equal(item.expected, x)
	}
}

func TestConversion(t *testing.T) {
	assert := assert.New(t)

	{
		e, err := NewExpression("1.23", nil)
		assert.NoError(err)
		x, err := e.Eval(nil)
		assert.NoError(err)
		f, ok := x.(float64)
		assert.True(ok)
		assert.Equal(1.23, f)
	}
	{
		e, err := NewExpression(`"1,2,3"`, nil)
		assert.NoError(err)
		x, err := e.Eval(nil)
		assert.NoError(err)
		s, ok := x.(string)
		assert.True(ok)
		assert.Equal("1,2,3", s)
	}
	{
		e, err := NewExpression("true", nil)
		assert.NoError(err)
		x, err := e.Eval(nil)
		assert.NoError(err)
		b, ok := x.(bool)
		assert.True(ok)
		assert.Equal(true, b)
	}
}
func TestFunction(t *testing.T) {
	assert := assert.New(t)

	f := func(args ...interface{}) (interface{}, error) {
		length := len(args[0].(string))
		return (float64)(length), nil
	}
	g := func(args ...interface{}) (interface{}, error) {
		length := len(args[0].(string))
		return -(float64)(length), nil
	}
	funcs := Functions{
		"strlen":    f,
		"negstrlen": g,
	}

	e, err := NewExpression(`strlen("abc") * negstrlen("4567")`, funcs)
	assert.NoError(err)
	x, err := e.Eval(nil)
	assert.NoError(err)
	v, ok := x.(float64)
	assert.True(ok)
	assert.Equal(-12.0, v)
}
