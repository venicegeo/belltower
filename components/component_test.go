package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

func TestComponent(t *testing.T) {
	assert := assert.New(t)

	Factory.Register("Foo", &foo{})

	config := fooInputData{
		X: 17,
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
	assert.Equal(11+19+17, out.(fooOutputData).Y)
}

//---------------------------------------------------------------------

type fooConfigData struct {
	I int
}

type fooInputData struct {
	X int
}

type fooOutputData struct {
	Y int
}

type foo struct {
	ComponentCore
	myint      int
	myotherint int
}

func (x *foo) Configure() error {
	data := fooConfigData{}
	err := x.config.ToStruct(&data)
	if err != nil {
		return err
	}

	x.myotherint = data.I
	x.myint = 19

	return nil
}

func (x *foo) Run(in interface{}) (interface{}, error) {
	input := in.(fooInputData)

	out := fooOutputData{
		Y: input.X + x.myint + x.myotherint,
	}

	return out, nil
}
