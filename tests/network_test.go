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

package tests

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/venicegeo/belltower/components"
	"github.com/venicegeo/belltower/engine"

	"github.com/stretchr/testify/assert"
)

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

//---------------------------------------------------------------------

func TestFlow(t *testing.T) {
	assert := assert.New(t)

	const logfile = "testflow.log"
	_ = os.Remove(logfile)

	/*
		components := []*engine.ComponentModel{
			&engine.ComponentModel{
				Name: "START",
				Type: "Starter",
			},
			&engine.ComponentModel{
				Name: "STOP",
				Type: "Stopper",
			},
			&engine.ComponentModel{
				Name: "myticker",
				Type: "Ticker",
				Config: engine.ArgMap{
					"limit": 5.0,
				},
			},
			&engine.ComponentModel{
				Name: "mylogger",
				Type: "Logger",
				Config: engine.ArgMap{
					"FileName": logfile,
				},
			},
			&engine.ComponentModel{
				Name: "myadder",
				Type: "Adder",
				Config: engine.ArgMap{
					"addend": 10.0,
				},
			},
			&engine.ComponentModel{
				Name: "myremapper",
				Type: "Remapper",
				Config: engine.ArgMap{
					"remaps": map[string]string{"Count": "Value"},
				},
			},
		}

		connections := []*engine.ConnectionModel{
			&engine.ConnectionModel{Source: "START.Output", Destination: "myticker.Input"},
			&engine.ConnectionModel{Source: "myticker.Output", Destination: "myremapper.Input"},
			&engine.ConnectionModel{Source: "myremapper.Output", Destination: "myadder.Input"},
			&engine.ConnectionModel{Source: "myadder.Output", Destination: "mylogger.Input"},
			&engine.ConnectionModel{Source: "mylogger.Output", Destination: "STOP.Input"},
		}

		g := &engine.GraphModel{
			Components:  components,
			Connections: connections,
		}
	*/

	dsl := `
	graph mygraph
		component START
			type: Starter
		end
		component STOP
			type: Stopper
		end
		component myticker
			type: Ticker
			config
				limit: 5.0
			end
		end
		component mylogger
			type: Logger
			config
				fileName: logfile
			end
		end
		component myadder
			type: Adder
			config
				addend: 10.0
			end
		end
		component myremapper
			type: Remapper
			config
				remaps: {"Count": "Value"}
			end
		end
	
 		START.Output -> myticker.Input
 		myticker.Output -> myremapper.Input
 		myremapper.Output -> myadder.Input
 		myadder.Output -> mylogger.Input
 		mylogger.Output -> STOP.Input
	end
	`

	g, err := engine.ParseDSL(dsl)
	assert.NoError(err)

	net, err := engine.NewNetwork(g)
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
	assertLogContainsLines(t, logfile, expected)

	_ = os.Remove(logfile)
}

func TestTwoOutputs(t *testing.T) {
	assert := assert.New(t)

	const logfile = "testtwooutput.log"
	_ = os.Remove(logfile)

	components := []*engine.ComponentModel{
		&engine.ComponentModel{
			Name: "START",
			Type: "Starter",
		},
		&engine.ComponentModel{
			Name: "STOP",
			Type: "Stopper",
		},
		&engine.ComponentModel{
			Name: "mycopier",
			Type: "Copier",
		},
		&engine.ComponentModel{
			Name: "myticker1",
			Type: "Ticker",
			Config: engine.ArgMap{
				"limit": 3.0,
			},
		},
		&engine.ComponentModel{
			Name: "myticker10",
			Type: "Ticker",
			Config: engine.ArgMap{
				"initialValue": 10,
				"limit":        13.0,
			},
		},
		&engine.ComponentModel{
			Name: "mylogger",
			Type: "Logger",
			Config: engine.ArgMap{
				"FileName": logfile,
			},
		},
	}

	// two outputs tied to same single input
	connections := []*engine.ConnectionModel{
		&engine.ConnectionModel{Source: "START.Output", Destination: "mycopier.Input"},
		&engine.ConnectionModel{Source: "mycopier.Output1", Destination: "myticker1.Input"},
		&engine.ConnectionModel{Source: "mycopier.Output2", Destination: "myticker10.Input"},
		&engine.ConnectionModel{Source: "myticker1.Output", Destination: "mylogger.Input"},
		&engine.ConnectionModel{Source: "myticker10.Output", Destination: "mylogger.Input"},
		&engine.ConnectionModel{Source: "mylogger.Output", Destination: "STOP.Input"},
	}

	g := &engine.GraphModel{
		Components:  components,
		Connections: connections,
	}

	net, err := engine.NewNetwork(g)
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
	assertLogContainsLines(t, logfile, expected)

	_ = os.Remove(logfile)
}
