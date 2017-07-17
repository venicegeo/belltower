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
			Name: "START",
			Type: "Starter",
		},
		&common.Component{
			Name: "STOP",
			Type: "Stopper",
		},
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
				"remaps": map[string]string{"Count": "Value"},
			},
		},
	}

	connections := []*common.Connection{
		&common.Connection{Source: "START.Output", Destination: "myticker.Input"},
		&common.Connection{Source: "myticker.Output", Destination: "myremapper.Input"},
		&common.Connection{Source: "myremapper.Output", Destination: "myadder.Input"},
		&common.Connection{Source: "myadder.Output", Destination: "STOP.Input"},
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
