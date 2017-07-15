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
			Config: common.ArgMap{
				"limit": 5.0,
			},
		},
		&common.Component{
			Name: "myadder",
			Type: "Adder",
			Config: common.ArgMap{
				"addend": 10.0,
			},
		},
		&common.Component{
			Name: "myremapper",
			Type: "Remapper",
			Config: common.ArgMap{
				"remaps": map[string]string{"count": "value"},
			},
		},
	}

	connections := []*common.Connection{
		&common.Connection{
			Source:      "myticker.Output",
			Destination: "myremapper.Input",
		},
		&common.Connection{
			Source:      "myremapper.Output",
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
