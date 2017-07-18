package engine

import (
	"testing"

	"github.com/venicegeo/belltower/common"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	assert := assert.New(t)

	components := []*common.ComponentModel{
		&common.ComponentModel{
			Name: "START",
			Type: "Starter",
		},
		&common.ComponentModel{
			Name: "STOP",
			Type: "Stopper",
		},
		&common.ComponentModel{
			Name: "myticker",
			Type: "Ticker",
			Config: common.ArgMap{
				"limit": 5.0,
			},
		},
		&common.ComponentModel{
			Name: "myadder",
			Type: "Adder",
			Config: common.ArgMap{
				"addend": 10.0,
			},
		},
		&common.ComponentModel{
			Name: "myremapper",
			Type: "Remapper",
			Config: common.ArgMap{
				"remaps": map[string]string{"Count": "Value"},
			},
		},
	}

	connections := []*common.ConnectionModel{
		&common.ConnectionModel{Source: "START.Output", Destination: "myticker.Input"},
		&common.ConnectionModel{Source: "myticker.Output", Destination: "myremapper.Input"},
		&common.ConnectionModel{Source: "myremapper.Output", Destination: "myadder.Input"},
		&common.ConnectionModel{Source: "myadder.Output", Destination: "STOP.Input"},
	}

	g := &common.GraphModel{
		Components:  components,
		Connections: connections,
	}

	net, err := NewNetwork(g)
	assert.NoError(err)
	assert.NotNil(net)

	err = net.Execute(10)
	assert.NoError(err)
}
func TestCopier(t *testing.T) {
	assert := assert.New(t)

	components := []*common.ComponentModel{
		&common.ComponentModel{
			Name: "START",
			Type: "Starter",
		},
		&common.ComponentModel{
			Name: "STOP",
			Type: "Stopper",
		},
		&common.ComponentModel{
			Name: "myticker",
			Type: "Ticker",
			Config: common.ArgMap{
				"limit": 3.0,
			},
		},
		&common.ComponentModel{
			Name:   "mycopier",
			Type:   "Copier",
			Config: common.ArgMap{},
		},
	}

	connections := []*common.ConnectionModel{
		&common.ConnectionModel{Source: "START.Output", Destination: "myticker.Input"},
		&common.ConnectionModel{Source: "myticker.Output", Destination: "mycopier.Input"},
		&common.ConnectionModel{Source: "mycopier.Output1", Destination: "STOP.Input"},
		&common.ConnectionModel{Source: "mycopier.Output2", Destination: "STOP.Input"},
	}

	g := &common.GraphModel{
		Components:  components,
		Connections: connections,
	}

	net, err := NewNetwork(g)
	assert.NoError(err)
	assert.NotNil(net)

	err = net.Execute(6)
	assert.NoError(err)
}
