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
	env := NewVariables()
	env.SetVariables(init)

	assert.Equal("a", env.GetVariable("A"))
	assert.Equal(12, env.GetVariable("B"))
	assert.Equal(13.0, env.GetVariable("C"))
	assert.Equal(false, env.GetVariable("x"))
	assert.Equal(true, env.GetVariable("yy"))
	assert.EqualValues([]int{34, 56, 78}, env.GetVariable("ary"))
	assert.EqualValues(mystruct, env.GetVariable("foo"))

	// can we override?
	env.SetVariable("yy", false)
	assert.Equal(!true, env.GetVariable("yy"))

	// can we get something that doesn't exist?
	assert.Equal(nil, env.GetVariable("xyzzy"))
}
