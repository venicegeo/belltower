package expr

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

type Expression struct {
	text       string
	expression *govaluate.EvaluableExpression
}

func NewExpression(text string) (*Expression, error) {
	expression, err := govaluate.NewEvaluableExpression(text)
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

func (e *Expression) Eval() (result interface{}, err error) {
	success := false

	defer func() {
		if !success {
			_ = recover()
			err = fmt.Errorf("evaluation failed")
			result = nil
		}
	}()

	result, err = e.expression.Evaluate(nil)
	if err != nil {
		return nil, err
	}
	success = true

	return result, nil
}

/*func AsInt(x interface{}) *int {
	v, ok := x.(int)
	if ok {
		return &v
	}
	return nil
}*/

func AsFloat(x interface{}) *float64 {
	v, ok := x.(float64)
	if ok {
		return &v
	}
	return nil
}

func AsBool(x interface{}) *bool {
	v, ok := x.(bool)
	if ok {
		return &v
	}
	return nil
}

func AsString(x interface{}) *string {
	v, ok := x.(string)
	if ok {
		return &v
	}
	return nil
}
