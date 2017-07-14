package engine

import (
	"testing"

	"github.com/venicegeo/belltower/common"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	assert := assert.New(t)

	components := []*common.Component{
		&common.Component{
			Name: "myticker",
			Type: "Ticker",
			Config: common.Map{
				"max": 5,
			},
		},
		&common.Component{
			Name: "myadder",
			Type: "Adder",
			Config: common.Map{
				"addend": 5,
			},
		},
	}

	connections := []*common.Connection{
		&common.Connection{
			Source:      "myticker.Output",
			Destination: "myadder.Input",
		},
	}

	g := &common.Graph{
		Components:  components,
		Connections: connections,
	}

	net, err := NewNetwork(g)
	assert.NoError(err)
	assert.NotNil(net)

	err = net.Execute()
	assert.NoError(err)
}
