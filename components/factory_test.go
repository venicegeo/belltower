package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

func TestFactory(t *testing.T) {
	assert := assert.New(t)

	Factory.Register("Xyzzy", &xyzzy{})

	{
		// make sure we can init without a config map (e.g. to allow defaults to work)
		x, err := Factory.Create("Xyzzy", nil)
		assert.NoError(err)
		assert.NotNil(x)

		assert.Equal("Bob", x.Description().Name)
		assert.Equal(0, x.(*xyzzy).executionCount)
		assert.Nil(x.(*xyzzy).config["param"])
		assert.Equal(19, x.(*xyzzy).myint)
	}

	config := common.Map{
		"param": "seventeen",
	}
	x, err := Factory.Create("Xyzzy", config)
	assert.NoError(err)
	assert.NotNil(x)

	assert.Equal("Bob", x.Description().Name)
	assert.Equal(0, x.(*xyzzy).executionCount)
	assert.Equal("seventeen", x.(*xyzzy).config["param"])
	assert.Equal(19, x.(*xyzzy).myint)
}

//---------------------------------------------------------------------

type xyzzy struct {
	ComponentCore
	myint int
}

func (x *xyzzy) Description() *Description {
	return &Description{
		Name: "Bob",
	}
}

func (x *xyzzy) localConfigure() error {
	x.myint = 19
	return nil
}

func (x *xyzzy) Run(in common.Map) (common.Map, error) {
	return nil, nil
}
