package dsls

import (
	"fmt"
)

type Variables struct {
	vars map[string]interface{}
}

func NewVariables() *Variables {
	env := &Variables{
		vars: map[string]interface{}{},
	}
	return env
}

func (env *Variables) SetVariables(vars map[string]interface{}) {
	if vars != nil {
		for k, v := range vars {
			env.SetVariable(k, v)
		}
	}
}

func (env *Variables) SetVariable(name string, value interface{}) {
	env.vars[name] = value
}

func (env *Variables) GetVariable(name string) interface{} {
	return env.vars[name]
}

func (env *Variables) String() string {
	s := ""
	for k, v := range env.vars {
		s += fmt.Sprintf("\"%s\": %v\n", k, v)
	}
	return s
}
