package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

func TestTicker(t *testing.T) {
	assert := assert.New(t)

	config := common.ArgMap{
		"Limit": 3,
	}
	tickerX, err := common.Factory.Create("Ticker", config)
	assert.NoError(err)
	ticker := tickerX.(*Ticker)

	// this setup is normally done by goflow itself
	chIn := make(chan string)
	chOut := make(chan string)
	ticker.Input = chIn
	ticker.Output = chOut

	inJ := "{}"
	go ticker.OnInput(inJ)

	outJ := <-chOut
	assert.JSONEq(`{"Count":1}`, outJ)
}
