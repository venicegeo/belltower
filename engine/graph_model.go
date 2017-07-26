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

import "fmt"

type GraphModel struct {
	Name string `json:"name"`

	Components  []*ComponentModel  `json:"components"`
	Connections []*ConnectionModel `json:"connections"`
}

type ComponentModel struct {
	Name string `json:"name"`

	Type string `json:"type"`

	Precondition  string `json:"precondition,omitempty"`
	Postcondition string `json:"postcondition,omitempty"`
	Config        ArgMap
	/*Component      Component
	InConnections  []*ConnectionModel
	OutConnections []*ConnectionModel*/
}

type ConnectionModel struct {
	Source      string `json:"source"`      // component.port
	Destination string `json:"destination"` // component.port
}

func (m *GraphModel) Validate() error               { return nil } // TODO
func (m *GraphModel) ReadFromJSON(jsn string) error { return ReadFromJSON(jsn, m) }
func (m *GraphModel) WriteToJSON() (string, error)  { return WriteToJSON(m) }

func (m *ComponentModel) Validate() error               { return nil } // TODO
func (m *ComponentModel) ReadFromJSON(jsn string) error { return ReadFromJSON(jsn, m) }
func (m *ComponentModel) WriteToJSON() (string, error)  { return WriteToJSON(m) }

func (m *ConnectionModel) Validate() error               { return nil } // TODO
func (m *ConnectionModel) ReadFromJSON(jsn string) error { return ReadFromJSON(jsn, m) }
func (m *ConnectionModel) WriteToJSON() (string, error)  { return WriteToJSON(m) }

//---------------------------------------------------------------------

func (g *GraphModel) String() string {
	s := fmt.Sprintf("Name: %s\n", g.Name)
	for _, v := range g.Components {
		s += fmt.Sprintf("Component: %s\n", v)
	}
	for _, v := range g.Connections {
		s += fmt.Sprintf("Connection: %s\n", v)
	}
	return s
}

func (c *ComponentModel) String() string {
	s := fmt.Sprintf("Name: %s, Type: %s", c.Name, c.Type)
	return s
}

func (c *ConnectionModel) String() string {
	s := fmt.Sprintf("Source: %s, Dest: %s", c.Source, c.Destination)
	return s
}

/*
type ComponentVisitor func(*ComponentModel) error
type ConnectionVisitor func(*ConnectionModel) error

type Visitor struct {
	Graph             *GraphModel
	ComponentVisitor  ComponentVisitor
	ConnectionVisitor ConnectionVisitor
	visited           map[string]bool
}

func (v *Visitor) Visit() error {
	v.visited = map[string]bool{}

	for _, component := range v.Graph.Components {
		err := v.visit(component)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Visitor) visit(component *ComponentModel) error {

	if v.visited[component.Name] {
		return nil
	}

	if v.ComponentVisitor != nil {
		err := v.ComponentVisitor(component)
		if err != nil {
			return nil
		}
	}
	v.visited[component.Name] = true

	for _, connection := range component.OutConnections {
		if v.ConnectionVisitor != nil {
			err := v.ConnectionVisitor(connection)
			if err != nil {
				return err
			}
		}

		destName := strings.Split(connection.Destination, ".")[0]
		dest := v.Graph.Components[destName]
		err := v.visit(dest)
		if err != nil {
			return err
		}
	}

	return nil
}
*/
