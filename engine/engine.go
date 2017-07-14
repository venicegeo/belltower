package engine

import (
	"fmt"
	"log"

	"strings"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/components"

	flow "github.com/trustmaster/goflow"
)

type GraphImplementation struct {
	flow.Graph
}

func InterpretGraph(graph *common.Graph) (*GraphImplementation, error) {
	g := &GraphImplementation{}
	g.InitGraphState()

	for _, component := range graph.Components {

		c, err := components.Factory(component.Type, component.Config)
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
