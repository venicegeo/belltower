package dsls

import (
	"fmt"
	"log"

	"strings"

	flow "github.com/trustmaster/goflow"
)

//---------------------------------------------------------------------

type GraphImplementation struct {
	flow.Graph
}

type ComponentImplementation interface {
	// describes the port datatypes, etc
	Description() *ComponentDescription

	// configures the component before any executions start
	Configure(config Map) error

	// perform one execution
	Run(in Map) (out Map, err error)
}

type ComponentImplementationCore struct {
	config         Map
	precondition   *Expression
	postcondition  *Expression
	executionCount int

	flow.Component
}

func (c *ComponentImplementationCore) configureCore(config Map) error {

	c.config = config
	if c.config == nil {
		c.config = Map{}
	}

	c.executionCount = 0

	cond, ok := config.asValidString("precondition")
	if ok {
		e, err := NewExpression(cond, nil)
		if err != nil {
			return err
		}
		c.precondition = e
	}

	cond, ok = config.asValidString("postcondition")
	if ok {
		e, err := NewExpression(cond, nil)
		if err != nil {
			return err
		}
		c.postcondition = e
	}

	return nil
}

type ComponentDescription struct {
	Id       Id        `json:"id"`
	Name     string    `json:"name"`
	Metadata *Metadata `json:"metadata,omitempty"`

	Config *Port `json:"config,omitempty"`
	Input  *Port `json:"input,omitempty"`
	Output *Port `json:"output,omitempty"`
}

//---------------------------------------------------------------------

type AdderComponent struct {
	ComponentImplementationCore

	addend int

	Input  <-chan string
	Output chan<- string
}

func (adder *AdderComponent) Description() *ComponentDescription {
	configFields := []*DataType{
		NewScalarDataType("addend", TypeNameInt),
	}
	config := NewStructDataType("config", configFields)

	inFields := []*DataType{
		NewScalarDataType("x", TypeNameInt),
	}
	in := NewStructDataType("Input", inFields)

	outFields := []*DataType{
		NewScalarDataType("y", TypeNameInt),
	}
	out := NewStructDataType("Output", outFields)

	return &ComponentDescription{
		Name: "adder",
		Config: &Port{
			DataType: config,
		},
		Input: &Port{
			DataType: in,
		},
		Output: &Port{
			DataType: out,
		},
	}
}

func (adder *AdderComponent) Configure(config Map) error {
	err := adder.configureCore(config)
	if err != nil {
		return err
	}

	if !adder.config.isInt("addend") {
		return fmt.Errorf("required ticker field missing: addend")
	}
	adder.addend = adder.config.asInt("addend")

	return nil
}

func (adder *AdderComponent) OnInput(data string) {
	log.Printf("OnInput!")
	adder.Output <- "yow"
}

func (adder *AdderComponent) Run(in Map) (Map, error) {
	x := in["x"].(int)

	y := x + adder.addend

	out := Map{}
	out["y"] = y

	return out, nil
}

//---------------------------------------------------------------------

type TickerComponent struct {
	ComponentImplementationCore

	Input  <-chan string
	Output chan<- string

	max int
	cur int
}

func (ticker *TickerComponent) Description() *ComponentDescription {
	configFields := []*DataType{
		NewScalarDataType("max", TypeNameInt),
	}
	config := NewStructDataType("config", configFields)

	inFields := []*DataType{
		NewScalarDataType("x", TypeNameInt),
	}
	in := NewStructDataType("Input", inFields)

	outFields := []*DataType{
		NewScalarDataType("y", TypeNameInt),
	}
	out := NewStructDataType("Output", outFields)

	return &ComponentDescription{
		Name: "ticker",
		Config: &Port{
			DataType: config,
		},
		Input: &Port{
			DataType: in,
		},
		Output: &Port{
			DataType: out,
		},
	}
}

func (ticker *TickerComponent) Configure(config Map) error {
	err := ticker.configureCore(config)
	if err != nil {
		return err
	}

	if !ticker.config.isInt("max") {
		return fmt.Errorf("required ticker field missing: max")
	}
	ticker.max = ticker.config.asInt("max")

	return nil
}

func (ticker *TickerComponent) OnInput(data string) {
	log.Printf("OnInput ticker!")
	ticker.Output <- "baz"
}

func (ticker *TickerComponent) Run(in Map) (Map, error) {

	out := Map{}
	out["y"] = ticker.cur

	ticker.cur++

	return out, nil
}

//---------------------------------------------------------------------

func ComponentFactory(componentType string, config Map) (ComponentImplementation, error) {
	var ci ComponentImplementation

	switch componentType {
	case "Adder":
		ci = &AdderComponent{}
	case "Ticker":
		ci = &TickerComponent{}
	default:
		return nil, fmt.Errorf("component factory: invalid name: %s", componentType)
	}

	err := ci.Configure(config)
	if err != nil {
		return nil, err
	}

	return ci, nil
}

func InterpretGraph(graph *Graph) (*GraphImplementation, error) {
	g := &GraphImplementation{}
	g.InitGraphState()

	for _, component := range graph.Components {

		c, err := ComponentFactory(component.Type, component.Config)
		if err != nil {
			return nil, err
		}

		ok := g.Add(c, component.Name)
		if !ok {
			return nil, fmt.Errorf("failed to add component")
		}

		if component.Name == "myticker" {
			ok = g.MapInPort("IN", "myticker", "Input")
			if !ok {
				return nil, fmt.Errorf("failed to add InPort")
			}
		}
	}
	log.Printf("here0")
	for _, connection := range graph.Connections {

		send := strings.Split(connection.Source, ".")[0]
		sendPort := strings.Split(connection.Source, ".")[1]
		recv := strings.Split(connection.Destination, ".")[0]
		recvPort := strings.Split(connection.Destination, ".")[1]

		log.Printf("here3 %s %s %s %s", send, sendPort, recv, recvPort)

		ok := g.Connect(send, sendPort, recv, recvPort)
		if !ok {
			return nil, fmt.Errorf("failed to add connection")
		}
		log.Printf("here2 %s", connection.Source)
	}
	log.Printf("here1")

	return g, nil
}

func DoIt(g *GraphImplementation) error {

	// We need a channel to talk to it
	in := make(chan string)
	ok := g.SetInPort("IN", in)
	if !ok {
		return fmt.Errorf("failed to SetInPort")
	}

	// Run the net
	flow.RunNet(g)

	// Now we can send some names and see what happens
	in <- "John"

	// Close the input to shut the network down
	close(in)

	// Wait until the app has done its job
	<-g.Wait()

	return nil
}
