package engine

import (
	"fmt"

	"strings"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/components"

	"time"

	flow "github.com/trustmaster/goflow"
)

// Network represents a running instance of a common.Graph
// We can have several Networks running at the same time.
type Network struct {
	flow.Graph
}

func NewNetwork(graph *common.Graph) (*Network, error) {

	g := &Network{}
	g.InitGraphState()

	//
	// first add the component nodes
	//
	for _, component := range graph.Components {

		c, err := components.Factory.Create(component.Type, component.Config)
		if err != nil {
			return nil, err
		}

		ok := g.Add(c, component.Name)
		if !ok {
			return nil, fmt.Errorf("failed to add component: %s (type %s)", component.Name, component.Type)
		}

		if component.Type == "Ticker" {
			ok = g.MapInPort("INPORT", component.Name, "Input")
			if !ok {
				return nil, fmt.Errorf("failed to add InPort")
			}
		} else if component.Type == "Adder" {
			ok = g.MapOutPort("OUTPORT", component.Name, "Output")
			if !ok {
				return nil, fmt.Errorf("failed to add OutPort")
			}
		}
	}

	//
	// now add the connection edges
	//
	for _, connection := range graph.Connections {

		send := strings.Split(connection.Source, ".")[0]
		sendPort := strings.Split(connection.Source, ".")[1]
		recv := strings.Split(connection.Destination, ".")[0]
		recvPort := strings.Split(connection.Destination, ".")[1]

		ok := g.Connect(send, sendPort, recv, recvPort)
		if !ok {
			return nil, fmt.Errorf("failed to add connection: %s to %s", connection.Source, connection.Destination)
		}
		fmt.Printf("NewNetwork: connected: %s to %s\n", connection.Source, connection.Destination)
	}

	return g, nil
}

func (g *Network) Execute() error {

	// input channel
	in := make(chan string)
	ok := g.SetInPort("INPORT", in)
	if !ok {
		return fmt.Errorf("failed to SetInPort")
	}

	// output channel
	out := make(chan string)
	ok = g.SetOutPort("OUTPORT", out)
	if !ok {
		return fmt.Errorf("failed to SetOutPort")
	}

	// Run the net
	flow.RunNet(g)

	// kick it off
	in <- fmt.Sprintf("{}")

	go func() {
		time.Sleep(10 * time.Second)
		fmt.Printf("closing after 10 seconds\n")
		close(in)
	}()

	for result := range out {
		fmt.Printf("RESULT: %s\n", result)
	}

	// Wait until the app has done its job
	<-g.Wait()

	return nil
}