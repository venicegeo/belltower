package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	assert := assert.New(t)

	lines := `
graph

    components

        component
            note: "later this can be changed to a Frobber2 component"
            type: Frobber
            name: myfrobber
            precondition: true  // for now
            // because y can't be bigger than x
            postcondition: x >= y
            config
                x: 5
                y: "foo"
                z: struct
                    alpha: int
                    beta: int
                end struct
            endconfig
        endcomponent

    endcomponents

    connections
        connection
            from: component.port
            to: component.port
        endconnection
    endconnections
endgraph`

	tokenizer := Tokenizer{}
	err := tokenizer.Scan(lines)
	assert.NoError(err)
}
