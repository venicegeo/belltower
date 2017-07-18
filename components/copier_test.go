package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

func TestCopier(t *testing.T) {
	assert := assert.New(t)

	config := common.ArgMap{}
	copierX, err := common.Factory.Create("Copier", config)
	assert.NoError(err)
	copier := copierX.(*Copier)

	// this setup is normally done by goflow itself
	chIn := make(chan string)
	chOut1 := make(chan string)
	chOut2 := make(chan string)
	copier.Input = chIn
	copier.Output1 = chOut1
	copier.Output2 = chOut2

	inJ := `{
		"alpha": 11.0,
		"beta":  22.0,
		"gamma": 33.0
	}`
	go copier.OnInput(inJ)

	// check the returned result
	out1J := <-chOut1
	out2J := <-chOut2

	assert.JSONEq(inJ, out1J)
	assert.JSONEq(inJ, out2J)
}
