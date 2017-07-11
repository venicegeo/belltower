package expr

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

type EnvironmentVars struct {
	vars map[string]interface{}
}

type EnvironmentFuncs struct {
	funcs map[string]govaluate.ExpressionFunction
}

//---------------------------------------------------------------------------

func NewEnvironmentVars() *EnvironmentVars {
	env := &EnvironmentVars{
		vars: map[string]interface{}{},
	}
	return env
}

func (env *EnvironmentVars) SetVars(vars map[string]interface{}) {
	if vars != nil {
		for k, v := range vars {
			env.SetVar(k, v)
		}
	}
}

func (env *EnvironmentVars) SetVar(name string, value interface{}) {
	env.vars[name] = value
}

func (env *EnvironmentVars) GetVar(name string) interface{} {
	return env.vars[name]
}

func (env *EnvironmentVars) String() string {
	s := ""
	for k, v := range env.vars {
		s += fmt.Sprintf("\"%s\": %v\n", k, v)
	}
	return s
}

//---------------------------------------------------------------------------

func NewEnvironmentFuncs() *EnvironmentFuncs {
	env := &EnvironmentFuncs{
		funcs: map[string]govaluate.ExpressionFunction{},
	}
	return env
}

func (env *EnvironmentFuncs) SetFuncs(fncs map[string]govaluate.ExpressionFunction) {
	if fncs != nil {
		for k, v := range fncs {
			env.SetFunc(k, v)
		}
	}
}

func (env *EnvironmentFuncs) SetFunc(name string, value govaluate.ExpressionFunction) {
	env.funcs[name] = value
}

func (env *EnvironmentFuncs) GetFunc(name string) string {
	//f := env.funcs[name]
	return name
}

func (env *EnvironmentFuncs) String() string {
	s := ""
	for k, _ := range env.funcs {
		s += fmt.Sprintf("\"%s\": (func)\n", k)
	}
	return s
}
