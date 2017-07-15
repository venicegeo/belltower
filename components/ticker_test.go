package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

func TestTicker(t *testing.T) {
	assert := assert.New(t)

	Factory.Register("Ticker", &Ticker{})

	config := common.ArgMap{
		"Limit": 3,
	}
	ticker, err := Factory.Create("Ticker", config)
	assert.NoError(err)

	in := common.ArgMap{}
	out, err := ticker.Run(in)
	assert.NoError(err)

	assert.Equal(1.0, out["count"])
}
