package engine

import (
	"fmt"

	"strings"

	"time"

	flow "github.com/trustmaster/goflow"
)

// Network represents a running instance of a Graph
// We can have several Networks running at the same time.
type Network struct {
	flow.Graph

	model *GraphModel
}

func NewNetwork(graphModel *GraphModel) (*Network, error) {
	net := &Network{}
	net.model = graphModel

	net.InitGraphState()

	//
	// first add the component nodes
	//
	for _, componentModel := range graphModel.Components {

		c, err := Factory.Create(componentModel.Type, componentModel.Config)
		if err != nil {
			return nil, err
		}

		ok := net.Add(c, componentModel.Name)
		if !ok {
			return nil, fmt.Errorf("failed to add component: %s (type %s)", componentModel.Name, componentModel.Type)
		}
	}

	ok := net.MapInPort("START.Input", "START", "Input")
	if !ok {
		return nil, fmt.Errorf("failed to add InPort")
	}
	ok = net.MapOutPort("STOP.Output", "STOP", "Output")
	if !ok {
		return nil, fmt.Errorf("failed to add OutPort")
	}

	//
	// now add the connection edges
	//
	for _, connectionModel := range graphModel.Connections {

		send := strings.Split(connectionModel.Source, ".")[0]
		sendPort := strings.Split(connectionModel.Source, ".")[1]
		recv := strings.Split(connectionModel.Destination, ".")[0]
		recvPort := strings.Split(connectionModel.Destination, ".")[1]

		if !net.portExists(send, sendPort) {
			return nil, fmt.Errorf("output port '%s.%s' not found", send, sendPort)
		}
		if !net.portExists(recv, recvPort) {
			return nil, fmt.Errorf("input port '%s.%s' not found", recv, recvPort)
		}

		ok := net.Connect(send, sendPort, recv, recvPort)
		if !ok {
			return nil, fmt.Errorf("failed to add connection: %s to %s", connectionModel.Source, connectionModel.Destination)
		}
		fmt.Printf("NewNetwork: connected: %s to %s\n", connectionModel.Source, connectionModel.Destination)
	}

	return net, nil
}

func (net *Network) portExists(componentName string, portName string) bool {
	for _, component := range net.model.Components {
		// TODO: check also for correct port name (which will need to be determined via reflection)
		if component.Name == componentName {
			return true
		}
	}
	return false
}

// timeout is in seconds, zero means no timeout
func (g *Network) Execute(timeout int) error {

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
		time.Sleep(time.Duration(timeout) * time.Second)
		fmt.Printf("closing after %d seconds\n", timeout)
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
