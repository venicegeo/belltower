package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/engine"
)

func TestTicker(t *testing.T) {
	assert := assert.New(t)

	config := engine.ArgMap{
		"Limit": 3,
	}
	tickerX, err := engine.Factory.Create("Ticker", config)
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
