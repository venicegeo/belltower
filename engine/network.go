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
	"fmt"

	"strings"

	"time"

	flow "github.com/trustmaster/goflow"
	"github.com/venicegeo/belltower/mpg/merr"
)

// Network represents a running instance of a Graph
// We can have several Networks running at the same time.
type Network struct {
	flow.Graph

	model *GraphModel
}

func NewNetwork(graphModel *GraphModel) (net *Network, err error) {
	defer func() {
		if r := recover(); r != nil {
			net = nil
			err = merr.Newf("panic inside goflow (NewNetwork): %s", r)
		}
	}()

	net = &Network{}
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

		/*componentModel.Component = c*/
	}

	ok := net.MapInPort("START.Input", "START", "Input")
	if !ok {
		return nil, fmt.Errorf("failed to add InPort")
	}
	ok = net.MapOutPort("STOP.Output", "STOP", "Output")
	if !ok {
		return nil, fmt.Errorf("failed to add OutPort")
	}

	/*for _, componentModel := range graphModel.Components {
		for _, connectionModel := range graphModel.Connections {

			send := strings.Split(connectionModel.Source, ".")[0]
			//sendPort := strings.Split(connectionModel.Source, ".")[1]
			recv := strings.Split(connectionModel.Destination, ".")[0]
			//recvPort := strings.Split(connectionModel.Destination, ".")[1]
			name := connectionModel.Source + "." + connectionModel.Destination

			if send == componentModel.Name {
				componentModel.OutConnections[name] = connectionModel
			} else if recv == componentModel.Name {
				componentModel.InConnections[name] = connectionModel
			}
		}
	}*/

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
func (g *Network) Execute(timeout int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = merr.Newf("panic inside goflow (Execute): %s", r)
		}
	}()

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
