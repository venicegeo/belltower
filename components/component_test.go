package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

func TestComponent(t *testing.T) {
	assert := assert.New(t)

	Factory.Register("Foo", &foo{})

	config := common.ArgMap{
		"myotherint": 17,
	}
	x, err := Factory.Create("Foo", config)
	assert.NoError(err)
	assert.NotNil(x)

	assert.Equal(0, x.(*foo).executionCount)
	assert.Equal(17, x.(*foo).config["myotherint"])
	assert.Equal(19, x.(*foo).myint)

	// does Run() work?
	in := common.ArgMap{"x": 11}
	out, err := x.Run(in)
	assert.Equal(11+19+17, out["y"])
}

//---------------------------------------------------------------------

type foo struct {
	ComponentCore
	myint int
}

func (x *foo) localConfigure() error {
	x.myint = 19
	return nil
}

func (x *foo) Run(in common.ArgMap) (common.ArgMap, error) {
	out := common.ArgMap{}

	out["y"] = in["x"].(int) + x.myint + x.config["myotherint"].(int)

	return out, nil
}
