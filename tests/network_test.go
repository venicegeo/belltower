package tests

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
			Name: "myticker",
			Type: "Ticker",
			Config: common.ArgMap{
				"limit": 5.0,
			},
		},
		&common.ComponentModel{
			Name: "mylogger",
			Type: "Logger",
			Config: common.ArgMap{
				"FileName": logfile,
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
		&common.ConnectionModel{Source: "myadder.Output", Destination: "mylogger.Input"},
		&common.ConnectionModel{Source: "mylogger.Output", Destination: "STOP.Input"},
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
			Name: "myticker1",
			Type: "Ticker",
			Config: common.ArgMap{
				"limit": 3.0,
			},
		},
		&common.ComponentModel{
			Name: "myticker10",
			Type: "Ticker",
			Config: common.ArgMap{
				"initialValue": 10,
				"limit":        13.0,
			},
		},
		&common.ComponentModel{
			Name: "mylogger",
			Type: "Logger",
			Config: common.ArgMap{
				"FileName": logfile,
			},
		},
	}

	// two outputs tied to same single input
	connections := []*common.ConnectionModel{
		&common.ConnectionModel{Source: "START.Output", Destination: "mycopier.Input"},
		&common.ConnectionModel{Source: "mycopier.Output1", Destination: "myticker1.Input"},
		&common.ConnectionModel{Source: "mycopier.Output2", Destination: "myticker10.Input"},
		&common.ConnectionModel{Source: "myticker1.Output", Destination: "mylogger.Input"},
		&common.ConnectionModel{Source: "myticker10.Output", Destination: "mylogger.Input"},
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
