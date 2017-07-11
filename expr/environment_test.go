package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentVars(t *testing.T) {
	assert := assert.New(t)

	mystruct := struct {
		frob int
		nitz string
	}{
		frob: 19,
		nitz: "ztin",
	}

	init := map[string]interface{}{
		"A":   "a",
		"B":   12,
		"C":   13.0,
		"x":   false,
		"yy":  true,
		"ary": []int{34, 56, 78},
		"foo": mystruct,
	}
	env := NewEnvironmentVars()
	env.SetVars(init)

	assert.Equal("a", env.GetVar("A"))
	assert.Equal(12, env.GetVar("B"))
	assert.Equal(13.0, env.GetVar("C"))
	assert.Equal(false, env.GetVar("x"))
	assert.Equal(true, env.GetVar("yy"))
	assert.EqualValues([]int{34, 56, 78}, env.GetVar("ary"))
	assert.EqualValues(mystruct, env.GetVar("foo"))

	// can we override?
	env.SetVar("yy", false)
	assert.Equal(!true, env.GetVar("yy"))

	// can we get something that doesn't exist?
	assert.Equal(nil, env.GetVar("xyzzy"))
}
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

	env := NewEnvironmentFuncs()
	env.SetFuncs(m)

	assert.NotNil(env.GetFunc("strlen"))
	assert.NotNil(env.GetFunc("negstrlen"))

	// can we override?
	env.SetFunc("strlen", g)
	assert.NotNil(env.GetFunc("strlen"))

	// can we get something that doesn't exist?
	assert.Nil(env.GetFunc("xyzzy"))
}
