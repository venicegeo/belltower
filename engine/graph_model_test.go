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
