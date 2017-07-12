package expr

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

type Function govaluate.ExpressionFunction

type Functions struct {
	funcs map[string]govaluate.ExpressionFunction
}

func NewFunctions() *Functions {
	env := &Functions{
		funcs: map[string]govaluate.ExpressionFunction{},
	}
	return env
}

func (env *Functions) SetFunctions(fncs map[string]Function) {
	if fncs != nil {
		for k, v := range fncs {
			env.SetFunction(k, v)
		}
	}
}

func (env *Functions) SetFunction(name string, value Function) {
	env.funcs[name] = govaluate.ExpressionFunction(value)
}

func (env *Functions) GetFunction(name string) Function {
	f := env.funcs[name]
	return Function(f)
}

func (env *Functions) String() string {
	s := ""
	for k, _ := range env.funcs {
		s += fmt.Sprintf("\"%s\": (func)\n", k)
	}
	return s
}
