package engine

import (
	"os"
	"testing"

	"github.com/venicegeo/belltower/common"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	assert := assert.New(t)

	const logfile = "testflow.log"
	_ = os.Remove(logfile)

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
			Name: "mysender",
			Type: "MySender",
		},
		&common.ComponentModel{
			Name: "myreceiver",
			Type: "MyReceiver",
		},
		&common.ComponentModel{
			Name: "myadder",
			Type: "MyAdder",
			Config: common.ArgMap{
				"addend": 10.0,
			},
		},
	}

	connections := []*common.ConnectionModel{
		&common.ConnectionModel{Source: "START.Output", Destination: "mysender.Input"},
		&common.ConnectionModel{Source: "mysender.Output", Destination: "myadder.Input"},
		&common.ConnectionModel{Source: "myadder.Output", Destination: "myreceiver.Input"},
		&common.ConnectionModel{Source: "myreceiver.Output", Destination: "STOP.Input"},
	}

	g := &common.GraphModel{
		Components:  components,
		Connections: connections,
	}

	net, err := NewNetwork(g)
	assert.NoError(err)
	assert.NotNil(net)

	err = net.Execute(7)
	assert.NoError(err)

	expected := []string{
		`{"Sum":11}`,
		`{"Sum":12}`,
		`{"Sum":13}`,
		`{"Sum":14}`,
		`{"Sum":15}`,
	}
	common.AssertLogContainsLines(t, logfile, expected)

	_ = os.Remove(logfile)
}

func TestTwoOutputs(t *testing.T) {
	assert := assert.New(t)

	const logfile = "testtwooutput.log"
	_ = os.Remove(logfile)

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
			Name: "mycopier",
			Type: "Copier",
		},
		&common.ComponentModel{
			Name: "mysender",
			Type: "MySender",
		},
		&common.ComponentModel{
			Name: "myreceiver",
			Type: "MyReceiver",
		},
	}

	// two outputs tied to same single input
	connections := []*common.ConnectionModel{
		&common.ConnectionModel{Source: "START.Output", Destination: "mysender.Input"},
		&common.ConnectionModel{Source: "mysender.Output", Destination: "mycopier.Input"},
		&common.ConnectionModel{Source: "mycopier.Output1", Destination: "myreceiver.Input"},
		&common.ConnectionModel{Source: "mycopier.Output2", Destination: "myreceiver.Input"},
		&common.ConnectionModel{Source: "mylogger.Output", Destination: "STOP.Input"},
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

	expected := []string{
		`{"Count":1}`,
		`{"Count":2}`,
		`{"Count":3}`,
		`{"Count":11}`,
		`{"Count":12}`,
		`{"Count":13}`,
	}
	common.AssertLogContainsLines(t, logfile, expected)

	_ = os.Remove(logfile)
}
