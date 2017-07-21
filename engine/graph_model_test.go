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
	"testing"

	"github.com/stretchr/testify/assert"
)

type tuple struct {
	valid    bool
	obj      Serializer
	json     string
	expected Serializer
}

var component1 = tuple{
	valid: true,
	obj:   &ComponentModel{},
	json: `{
		"name": "ticker",
		"type": "Ticker",
		"precondition": "x>y",
		"postcondition": "x<y"
	}`,
	expected: &ComponentModel{
		Name:          "ticker",
		Type:          "Ticker",
		Precondition:  "x>y",
		Postcondition: "x<y",
	},
}

var component2 = tuple{
	valid: true,
	obj:   &ComponentModel{},
	json: `{
		"name": "ticker2",
		"type": "Ticker",
		"precondition": "x>y",
		"postcondition": "x<y"
	}`,
	expected: &ComponentModel{
		Name:          "ticker2",
		Type:          "Ticker",
		Precondition:  "x>y",
		Postcondition: "x<y",
	},
}

var connection1 = tuple{
	valid: true,
	obj:   &ConnectionModel{},
	json: `{
		"source": "foo.out",
		"destination": "bar.in"
	}`,
	expected: &ConnectionModel{
		Source:      "foo.out",
		Destination: "bar.in",
	},
}

var graph1 = tuple{
	valid: true,
	obj:   &GraphModel{},
	json: `{
		"name": "ticker", 
		"components": [
			` + component1.json + `,
			` + component2.json + `
		],
		"connections": [
			` + connection1.json + `
		]
	}`,
	expected: &GraphModel{
		Name: "ticker",
		Components: []*ComponentModel{
			component1.expected.(*ComponentModel),
			component2.expected.(*ComponentModel),
		},
		Connections: []*ConnectionModel{
			connection1.expected.(*ConnectionModel),
		},
	},
}

func TestGraphTypes(t *testing.T) {
	assert := assert.New(t)

	items := []tuple{
		component1,
		connection1,
		graph1,
	}

	for _, item := range items {

		//t.Log("----------------\n", item.jsn, "-----------------\n")

		err := item.obj.ReadFromJSON(item.json)
		if !item.valid {
			assert.Error(err)
			continue
		}
		assert.NoError(err)

		assert.Equal(item.expected, item.obj)
	}
}
