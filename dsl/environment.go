package dsl

type Environment struct {
	data map[string]ExprValue
}

func NewEnvironment() *Environment {
	env := &Environment{
		data: map[string]ExprValue{},
	}
	return env
}

func (env *Environment) set(name string, value ExprValue) {
	env.data[name] = value
}

func (env *Environment) get(name string) ExprValue {
	return env.data[name]
}
