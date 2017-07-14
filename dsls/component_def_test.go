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
				"max": 5,
			},
		},
		&Component{
			Name: "myadder",
			Type: "Adder",
			Config: Map{
				"addend": 5,
			},
		},
	}

	connections := []*Connection{
		&Connection{
			Source:      "myticker.Output",
			Destination: "myadder.Input",
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
