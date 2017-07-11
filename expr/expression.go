package expr

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

type Expression struct {
	text       string
	expression *govaluate.EvaluableExpression
}

func NewExpression(text string, env *EnvironmentFuncs) (*Expression, error) {
	funcs := map[string]govaluate.ExpressionFunction{}
	if env != nil {
		funcs = env.funcs
	}

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

func (e *Expression) Eval(env *EnvironmentVars) (result interface{}, err error) {
	success := false

	defer func() {
		if !success {
			_ = recover()
			err = fmt.Errorf("evaluation failed")
			result = nil
		}
	}()

	envVars := map[string]interface{}{}
	if env != nil {
		envVars = env.vars
	}

	result, err = e.expression.Evaluate(envVars)
	if err != nil {
		return nil, err
	}
	success = true

	return result, nil
}
