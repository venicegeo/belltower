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
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/venicegeo/belltower/mpg/merr"
)

//---------------------------------------------------------------------

type Function govaluate.ExpressionFunction

type Functions map[string]govaluate.ExpressionFunction

func (fs Functions) String() string {
	s := ""
	for k, _ := range fs {
		s += fmt.Sprintf("\"%s\": (func)\n", k)
	}
	return s
}

//---------------------------------------------------------------------

type Variables map[string]interface{}

func (vs Variables) String() string {
	s := ""
	for k, v := range vs {
		s += fmt.Sprintf("\"%s\": %v\n", k, v)
	}
	return s
}

//---------------------------------------------------------------------

type Expression struct {
	text       string
	expression *govaluate.EvaluableExpression
}

func NewExpression(text string, funcs Functions) (e *Expression, err error) {

	defer func() {
		if r := recover(); r != nil {
			e = nil
			err = merr.Newf("panic inside govaluate.NewEvaluableExpressionWithFunctions: %s", r)
		}
	}()

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(text, funcs)
	if err != nil {
		return nil, err
	}

	e = &Expression{
		text:       text,
		expression: expression,
	}

	return e, nil
}

func (e *Expression) String() string {
	return e.expression.String()
}

func (e *Expression) Eval(vars Variables) (result interface{}, err error) {

	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = merr.Newf("panic inside govaluate.Evaluate: %s", r)
		}
	}()

	result, err = e.expression.Evaluate(vars)
	if err != nil {
		return nil, err
	}

	return result, nil
}
