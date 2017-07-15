package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

func TestAdder(t *testing.T) {
	assert := assert.New(t)

	Factory.Register("Adder", &Adder{})

	config := common.ArgMap{
		"addend": 3,
	}
	adder, err := Factory.Create("Adder", config)
	assert.NoError(err)

	in := common.ArgMap{
		"value": 7,
	}
	out, err := adder.Run(in)
	assert.NoError(err)

	assert.Equal(10.0, out["sum"])
}
