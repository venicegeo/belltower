/* Copyright 2017, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package engine

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	assert := assert.New(t)

	components := []*ComponentModel{
		&ComponentModel{
			Name: "START",
			Type: "Starter",
		},
		&ComponentModel{
			Name: "STOP",
			Type: "Stopper",
		},
		&ComponentModel{
			Name: "mysender",
			Type: "MySender",
			Config: ArgMap{
				"i": 17,
			},
		},
		&ComponentModel{
			Name: "myreceiver",
			Type: "MyReceiver",
		},
		&ComponentModel{
			Name: "myadder",
			Type: "MyAdder",
			Config: ArgMap{
				"addend": 13.0,
			},
		},
	}

	connections := []*ConnectionModel{
		&ConnectionModel{Source: "START.Output", Destination: "mysender.Input"},
		&ConnectionModel{Source: "mysender.Output", Destination: "myadder.Input"},
		&ConnectionModel{Source: "myadder.Output", Destination: "myreceiver.Input"},
		&ConnectionModel{Source: "myreceiver.Output", Destination: "STOP.Input"},
	}

	g := &GraphModel{
		Components:  components,
		Connections: connections,
	}

	net, err := NewNetwork(g)
	assert.NoError(err)
	assert.NotNil(net)

	err = net.Execute(7)
	assert.NoError(err)

	assert.Equal(13, myReceiver.sum)

	assert.NotNil(net)
}

func TestTwoOutputs(t *testing.T) {
	assert := assert.New(t)

	const logfile = "testtwooutput.log"
	_ = os.Remove(logfile)

	components := []*ComponentModel{
		&ComponentModel{
			Name: "START",
			Type: "Starter",
		},
		&ComponentModel{
			Name: "STOP",
			Type: "Stopper",
		},
		&ComponentModel{
			Name: "mycopier",
			Type: "MyCopier",
		},
		&ComponentModel{
			Name: "mysender",
			Type: "MySender",
			Config: ArgMap{
				"i": 17,
			},
		},
		&ComponentModel{
			Name: "myreceiver",
			Type: "MyReceiver",
		},
	}

	// two outputs tied to same single input
	connections := []*ConnectionModel{
		&ConnectionModel{Source: "START.Output", Destination: "mysender.Input"},
		&ConnectionModel{Source: "mysender.Output", Destination: "mycopier.Input"},
		&ConnectionModel{Source: "mycopier.Output1", Destination: "myreceiver.Input"},
		&ConnectionModel{Source: "mycopier.Output2", Destination: "myreceiver.Input"},
		&ConnectionModel{Source: "myreceiver.Output", Destination: "STOP.Input"},
	}

	g := &GraphModel{
		Components:  components,
		Connections: connections,
	}

	net, err := NewNetwork(g)
	assert.NoError(err)
	assert.NotNil(net)

	err = net.Execute(6)
	assert.NoError(err)

	assert.Equal(34, myReceiver.i)

	_ = os.Remove(logfile)
}

/*
func TestVisitor(t *testing.T) {
	assert := assert.New(t)

	components := []*ComponentModel{
		&ComponentModel{
			Name: "START",
			Type: "Starter",
		},
		&ComponentModel{
			Name: "STOP",
			Type: "Stopper",
		},
		&ComponentModel{
			Name: "mycopier",
			Type: "MyCopier",
		},
		&ComponentModel{
			Name: "mysender",
			Type: "MySender",
			Config: ArgMap{
				"i": 17,
			},
		},
		&ComponentModel{
			Name: "myreceiver",
			Type: "MyReceiver",
		},
	}

	// two outputs tied to same single input
	connections := []*ConnectionModel{
		&ConnectionModel{Source: "START.Output", Destination: "mysender.Input"},
		&ConnectionModel{Source: "mysender.Output", Destination: "mycopier.Input"},
		&ConnectionModel{Source: "mycopier.Output1", Destination: "myreceiver.Input"},
		&ConnectionModel{Source: "mycopier.Output2", Destination: "myreceiver.Input"},
		&ConnectionModel{Source: "myreceiver.Output", Destination: "STOP.Input"},
	}

	g := &GraphModel{
		Components:  components,
		Connections: connections,
	}

	net, err := NewNetwork(g)
	assert.NoError(err)
	assert.NotNil(net)

	visitor := &Visitor{
		Graph: g,
		ComponentVisitor: func(comp *ComponentModel) error {
			mlog.Debug(comp)
			return nil
		},
		ConnectionVisitor: func(conn *ConnectionModel) error {
			mlog.Debug(conn)
			return nil
		},
	}

	err = visitor.Visit()
	assert.NoError(err)
}
*/

//---------------------------------------------------------------------

func readFile(t *testing.T, filename string) string {
	assert := assert.New(t)

	byts, err := ioutil.ReadFile(filename)
	assert.NoError(err)

	return string(byts)
}

func readLines(t *testing.T, filename string) []string {
	assert := assert.New(t)

	lines := []string{}

	file, err := os.Open(filename)
	assert.NoError(err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	assert.NoError(scanner.Err())

	return lines
}

// we want the files to contain the expected lines, but in no particular order
func assertLogContainsLines(t *testing.T, filename string, expected []string) {
	actual := readLines(t, filename)
	assert.Len(t, actual, len(expected))
	for _, v := range expected {
		assert.Contains(t, actual, v)
	}
}
