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

	assert.Equal(11, out["omega"])
	assert.Equal(22, out["psi"])
	assert.Equal(33, out["gamma"])
	assert.NotContains(out, "alpha")
	assert.NotContains(out, "beta")
}
