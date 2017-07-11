package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
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
