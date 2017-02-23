package dsl

type Eval struct {
}

func (e *Eval) Evaluate(expr ExprNode, env *Environment) (ExprValue, error) {
	result := expr.Eval(env)
	return result, nil
}
