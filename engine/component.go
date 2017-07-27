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
	"github.com/trustmaster/goflow"
	"github.com/venicegeo/belltower/mpg/mlog"
)

type Component interface {
	// the system calls this to do init work specific to your component type
	// you must implement this
	Configure() error

	// the system calls this to do init work for ComponentCore fields
	// do not implement this yourself
	coreConfigure(config ArgMap) error

	// Implement "OnPORT" for each input port, using the right signature and contents. For example:
	//
	// func (adder *Adder) OnInput(inputJson string) {
	//   input := &AdderInputData{}
	//   err := input.ReadFromJSON(inputJson)
	//   if err != nil {...}
	//   output := &AdderOutputData{}
	//   output.x = ...
	//   outputJson, err := output.WriteToJSON()
	//   if err != nil {...}
	//   adder.Output <- outputJson
	// }
	//
	// See adder.go for a complete example.
}

type ComponentCore struct {
	Config         ArgMap
	precondition   *Expression
	postcondition  *Expression
	executionCount int

	flow.Component
}

func (c *ComponentCore) coreConfigure(config ArgMap) error {
	// TODO: can this be replaced by Init() somehow?

	c.Config = config

	c.executionCount = 0

	cond, err := config.GetStringOrDefault("precondition", "")
	if err != nil {
		return err
	}
	if cond != "" {
		e, err := NewExpression(cond, nil)
		if err != nil {
			return err
		}
		c.precondition = e
	}

	cond, err = config.GetStringOrDefault("postcondition", "")
	if err != nil {
		return nil
	}
	if cond != "" {
		e, err := NewExpression(cond, nil)
		if err != nil {
			return err
		}
		c.postcondition = e
	}

	return nil
}

//---------------------------------------------------------------------

func init() {
	Factory.Register("Starter", &Starter{})
	Factory.Register("Stopper", &Stopper{})
}

type Starter struct {
	ComponentCore
	Input  <-chan string
	Output chan<- string
}

func (*Starter) Configure() error { return nil }
func (s *Starter) OnInput(string) {
	mlog.Printf("Starter OnInput\n")
	s.Output <- "{}"
}

type Stopper struct {
	ComponentCore
	Input  <-chan string
	Output chan<- string
}

func (*Stopper) Configure() error { return nil }

func (s *Stopper) OnInput(string) {
	mlog.Printf("Stopper OnInput\n")
	s.Output <- "{}"
}
