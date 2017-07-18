package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	assert := assert.New(t)

	Factory.Register("Xyzzy", &xyzzy{})

	{
		// make sure we can init without a config map (e.g. to allow defaults to work)
		x, err := Factory.Create("Xyzzy", nil)
		assert.NoError(err)
		assert.NotNil(x)

		assert.Equal(0, x.(*xyzzy).executionCount)
		assert.Nil(x.(*xyzzy).Config["param"])
		assert.Equal(19, x.(*xyzzy).myint)
	}

	config := ArgMap{
		"param": "seventeen",
	}
	x, err := Factory.Create("Xyzzy", config)
	assert.NoError(err)
	assert.NotNil(x)
	xx := x.(*xyzzy)

	assert.Equal(0, xx.executionCount)
	assert.Equal("seventeen", xx.Config["param"])
	assert.Equal(19, xx.myint)
}

//---------------------------------------------------------------------

type xyzzy struct {
	ComponentCore
	myint int
}

func (x *xyzzy) Configure() error {
	x.myint = 19
	return nil
}

func (x *xyzzy) Run(in interface{}) (interface{}, error) {
	return nil, nil
}
