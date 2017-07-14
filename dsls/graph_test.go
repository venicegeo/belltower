package dsls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type tuple struct {
	valid    bool
	obj      Validater
	json     string
	expected Validater
}

var port1 = tuple{
	valid: true,
	obj:   &Port{},
	json: `{
		"id": "123",
		"name": "foo",
		"porttype": "input",
		"datatype": {
			"name": "x",
			"type": "int"
		}
	}`,
	expected: &Port{
		Name: "foo",
		Id:   "123",
		DataType: &DataType{
			Name: "x",
			Type: "int",
		},
		PortType: "input",
	},
}

var metadata1 = tuple{
	valid: true,
	obj:   &Metadata{},
	json: `{
		"contact": "mpg",
		"version": "1.2.3",
		"description": "xyzzy"
	}`,
	expected: &Metadata{
		Contact:     "mpg",
		Version:     "1.2.3",
		Description: "xyzzy",
	},
}

var componentType1 = tuple{
	valid: true,
	obj:   &ComponentType{},
	json: `{
		"id": "123",
		"name": "ticker",
		"metadata": ` + metadata1.json + `,
		"input": ` + port1.json + `,
		"output": ` + port1.json + `,
		"config": ` + port1.json + `
	}`,
	expected: &ComponentType{
		Id:       "123",
		Name:     "ticker",
		Metadata: metadata1.obj.(*Metadata),
		Input:    port1.obj.(*Port),
		Output:   port1.obj.(*Port),
		Config:   port1.obj.(*Port),
	},
}

var component1 = tuple{
	valid: true,
	obj:   &Component{},
	json: `{
		"id": "123",
		"name": "ticker",
		"metadata": ` + metadata1.json + `,
		"type": "Ticker",
		"precondition": "x>y",
		"postcondition": "x<y"
	}`,
	expected: &Component{
		Id:            "123",
		Name:          "ticker",
		Metadata:      metadata1.obj.(*Metadata),
		Type:          "Ticker",
		Precondition:  "x>y",
		Postcondition: "x<y",
	},
}

var component2 = tuple{
	valid: true,
	obj:   &Component{},
	json: `{
		"id": "124",
		"name": "ticker2",
		"metadata": ` + metadata1.json + `,
		"type": "Ticker",
		"precondition": "x>y",
		"postcondition": "x<y"
	}`,
	expected: &Component{
		Id:            "124",
		Name:          "ticker2",
		Metadata:      metadata1.obj.(*Metadata),
		Type:          "Ticker",
		Precondition:  "x>y",
		Postcondition: "x<y",
	},
}

var connection1 = tuple{
	valid: true,
	obj:   &Connection{},
	json: `{
		"id": "123",
		"name": "link",
		"source": "foo.out",
		"destination": "bar.in"
	}`,
	expected: &Connection{
		Id:          "123",
		Name:        "link",
		Source:      "foo.out",
		Destination: "bar.in",
	},
}

var graph1 = tuple{
	valid: true,
	obj:   &Graph{},
	json: `{
		"id": "123",
		"name": "ticker",
		"metadata": ` + metadata1.json + `,
		"components": [
			` + component1.json + `,
			` + component2.json + `
		],
		"connections": [
			` + connection1.json + `
		]
	}`,
	expected: &Graph{
		Id:       "123",
		Name:     "ticker",
		Metadata: metadata1.obj.(*Metadata),
		Components: []*Component{
			component1.expected.(*Component),
			component2.expected.(*Component),
		},
		Connections: []*Connection{
			connection1.expected.(*Connection),
		},
	},
}

func TestGraphTypes(t *testing.T) {
	assert := assert.New(t)

	items := []tuple{
		port1,
		metadata1,
		componentType1,
		component1,
		connection1,
		graph1,
	}

	for _, item := range items {

		//t.Log("----------------\n", item.jsn, "-----------------\n")

		p, err := NewObjectFromJSON(item.json, item.obj)
		if !item.valid {
			assert.Error(err)
			assert.Nil(p)
			continue
		}
		assert.NoError(err)
		assert.NotNil(p)

		assert.Equal(item.expected, p)
	}
}
