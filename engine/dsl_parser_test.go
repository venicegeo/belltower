package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/mpg/mlog"
)

// TODO: if a connection is "component.port -> ..."
// we have an ambiguity with the "component NAME" rule

func TestParser(t *testing.T) {
	assert := assert.New(t)

	lines := `
graph mygraph

        component myfrobber
            type: Frobber
            precondition: "true"  // for now
            // because y can't be bigger than x
            postcondition: "x >= y"
            config
                x: 5
                y: "foo"
                //z: 
                //    alpha: int
                //    beta: int
                //end struct
            end
        end

     myfrobber.port ->  myfrobber.port
end
`

	tokenizer := &Tokenizer{}
	err := tokenizer.Scan(lines)
	assert.NoError(err)

	parser := &Parser{
		tokenizer: tokenizer,
	}
	err = parser.parse()
	assert.NoError(err)

	mlog.Printf("---\n%s---\n", parser.graph.String())
}
