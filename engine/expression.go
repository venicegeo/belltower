package engine

import (
	"fmt"

	"github.com/Knetic/govaluate"
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

func NewExpression(text string, funcs Functions) (*Expression, error) {

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(text, funcs)
	if err != nil {
		return nil, err
	}

	e := &Expression{
		text:       text,
		expression: expression,
	}

	return e, nil
}

func (e *Expression) String() string {
	return e.expression.String()
}

func (e *Expression) Eval(vars Variables) (result interface{}, err error) {
	success := false

	defer func() {
		if !success {
			_ = recover()
			err = fmt.Errorf("evaluation failed")
			result = nil
		}
	}()

	result, err = e.expression.Evaluate(vars)
	if err != nil {
		return nil, err
	}
	success = true

	return result, nil
}
