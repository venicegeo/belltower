package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/engine"
)

func TestAdder(t *testing.T) {
	assert := assert.New(t)

	config := engine.ArgMap{
		"addend": 3,
	}
	adderX, err := engine.Factory.Create("Adder", config)
	assert.NoError(err)
	adder := adderX.(*Adder)

	// this setup is normally done by goflow itself
	chIn := make(chan string)
	chOut := make(chan string)
	adder.Input = chIn
	adder.Output = chOut

	inJ := `{"Value": 7}`
	go adder.OnInput(inJ)

	// check the returned result
	outJ := <-chOut
	assert.JSONEq(`{"Sum":10.0}`, outJ)
}
