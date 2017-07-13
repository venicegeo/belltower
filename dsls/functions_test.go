package dsls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentFuncs(t *testing.T) {
	assert := assert.New(t)

	var f Function = func(args ...interface{}) (interface{}, error) {
		length := len(args[0].(string))
		return (float64)(length), nil
	}
	var g Function = func(args ...interface{}) (interface{}, error) {
		length := len(args[0].(string))
		return -(float64)(length), nil
	}
	m := map[string]Function{
		"strlen":    f,
		"negstrlen": g,
	}

	env := NewFunctions()
	env.SetFunctions(m)

	assert.NotNil(env.GetFunction("strlen"))
	assert.NotNil(env.GetFunction("negstrlen"))

	// can we override?
	env.SetFunction("strlen", g)
	assert.NotNil(env.GetFunction("strlen"))

	// can we get something that doesn't exist?
	assert.Nil(env.GetFunction("xyzzy"))
}
