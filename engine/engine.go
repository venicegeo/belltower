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

func NewNetwork(graph *common.GraphModel) (*Network, error) {

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
	}

	ok := g.MapInPort("START.Input", "START", "Input")
	if !ok {
		return nil, fmt.Errorf("failed to add InPort")
	}
	ok = g.MapOutPort("STOP.Output", "STOP", "Output")
	if !ok {
		return nil, fmt.Errorf("failed to add OutPort")
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
	ok := g.SetInPort("START.Input", in)
	if !ok {
		return fmt.Errorf("failed to SetInPort")
	}

	// output channel
	out := make(chan string)
	ok = g.SetOutPort("STOP.Output", out)
	if !ok {
		return fmt.Errorf("failed to SetOutPort")
	}

	// start the net
	flow.RunNet(g)

	// send the initial input
	in <- fmt.Sprintf("{}")

	// set up the shutdown mechanism
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Printf("closing after 10 seconds\n")
		close(in)
	}()

	// read all the outputs
	for result := range out {
		fmt.Printf("RESULT: %s\n", result)
	}

	// Wait until the app has done its job
	<-g.Wait()

	return nil
}
