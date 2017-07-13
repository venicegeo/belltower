package dsls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type tuple struct {
	o Validater
	s string
	p Validater
}

var port1 = tuple{
	o: &Port{},
	s: `{
		"id": "123",
		"name": "foo",
		"porttype": "input",
		"datatype": {
			"name": "x",
			"type": "int"
		}
	}`,
	p: &Port{
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
	o: &Metadata{},
	s: `{
		"contact": "mpg",
		"version": "1.2.3",
		"description": "xyzzy"
	}`,
	p: &Metadata{
		Contact:     "mpg",
		Version:     "1.2.3",
		Description: "xyzzy",
	},
}

var componentType1 = tuple{
	o: &ComponentType{},
	s: `{
		"id": "123",
		"name": "ticker",
		"metadata": ` + metadata1.s + `,
		"input": ` + port1.s + `,
		"output": ` + port1.s + `,
		"config": ` + port1.s + `
	}`,
	p: &ComponentType{
		Id:       "123",
		Name:     "ticker",
		Metadata: metadata1.o.(*Metadata),
		Input:    port1.o.(*Port),
		Output:   port1.o.(*Port),
		Config:   port1.o.(*Port),
	},
}

var component1 = tuple{
	o: &Component{},
	s: `{
		"id": "123",
		"name": "ticker",
		"metadata": ` + metadata1.s + `,
		"type": "Ticker",
		"precondition": "x>y",
		"postcondition": "x<y"
	}`,
	p: &Component{
		Id:            "123",
		Name:          "ticker",
		Metadata:      metadata1.o.(*Metadata),
		Type:          "Ticker",
		Precondition:  "x>y",
		Postcondition: "x<y",
	},
}

var component2 = tuple{
	o: &Component{},
	s: `{
		"id": "124",
		"name": "ticker2",
		"metadata": ` + metadata1.s + `,
		"type": "Ticker",
		"precondition": "x>y",
		"postcondition": "x<y"
	}`,
	p: &Component{
		Id:            "124",
		Name:          "ticker2",
		Metadata:      metadata1.o.(*Metadata),
		Type:          "Ticker",
		Precondition:  "x>y",
		Postcondition: "x<y",
	},
}

var graph1 = tuple{
	o: &Graph{},
	s: `{
		"id": "123",
		"name": "ticker",
		"metadata": ` + metadata1.s + `,
		"components": [
			` + component1.s + `,
			` + component2.s + `
		]
	}`,
	p: &Graph{
		Id:       "123",
		Name:     "ticker",
		Metadata: metadata1.o.(*Metadata),
		Components: []*Component{
			component1.p.(*Component),
			component2.p.(*Component),
		},
	},
}

func TestGraphTypes(t *testing.T) {
	assert := assert.New(t)

	items := []struct {
		valid bool
		tuple tuple
	}{
		{
			tuple: port1,
			valid: true,
		},
		{
			tuple: metadata1,
			valid: true,
		},
		{
			tuple: componentType1,
			valid: true,
		},
		{
			tuple: component1,
			valid: true,
		},
		{
			tuple: graph1,
			valid: true,
		},
	}

	for _, item := range items {

		//t.Log("----------------\n", item.jsn, "-----------------\n")

		p, err := NewObjectFromJSON(item.tuple.s, item.tuple.o)
		if !item.valid {
			assert.Error(err)
			assert.Nil(p)
			continue
		}
		assert.NoError(err)
		assert.NotNil(p)

		assert.Equal(item.tuple.p, p)
	}
}
