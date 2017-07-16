package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

func TestRemapper(t *testing.T) {
	assert := assert.New(t)

	Factory.Register("Remapper", &Remapper{})

	config := common.ArgMap{
		"remaps": map[string]string{
			"alpha": "omega",
			"beta":  "psi",
		},
	}
	remapper, err := Factory.Create("Remapper", config)
	assert.NoError(err)

	in := common.ArgMap{
		"alpha": 11,
		"beta":  22,
		"gamma": 33,
	}
	out, err := remapper.Run(in)
	assert.NoError(err)

	outputMap := out.(common.ArgMap)
	assert.Equal(11, outputMap["omega"])
	assert.Equal(22, outputMap["psi"])
	assert.Equal(33, outputMap["gamma"])
	assert.NotContains(outputMap, "alpha")
	assert.NotContains(outputMap, "beta")
}
