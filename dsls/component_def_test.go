package dsls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	assert := assert.New(t)

	components := []*Component{
		&Component{
			Name: "myticker",
			Type: "Ticker",
			Config: Map{
				"x": 5,
			},
		},
		&Component{
			Name: "myadder",
			Type: "Adder",
			Config: Map{
				"x": 5,
			},
		},
	}

	connections := []*Connection{
		&Connection{
			Source:      "myticker.output",
			Destination: "myadder.input",
		},
	}

	g := &Graph{
		Components:  components,
		Connections: connections,
	}

	gi, err := InterpretGraph(g)
	assert.NoError(err)
	assert.NotNil(gi)

	err = DoIt(gi)
	assert.NoError(err)
}
