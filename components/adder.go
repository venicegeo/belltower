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
package components

import (
	"github.com/venicegeo/belltower/engine"
	"github.com/venicegeo/belltower/mpg/mlog"
)

func init() {
	engine.Factory.Register("Adder", &Adder{})
}

type AdderConfigData struct {

	// The value added to the input. Default is zero.
	Addend float64
}

// implements Serializer
type AdderInputData struct {

	// The value added to the addend from the configuration. Default is zero.
	Value float64
}

func (m *AdderInputData) Validate() error               { return nil } // TODO
func (m *AdderInputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *AdderInputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

// implements Serializer
type AdderOutputData struct {

	// Value of input value added to addend.
	Sum float64
}

func (m *AdderOutputData) Validate() error               { return nil } // TODO
func (m *AdderOutputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *AdderOutputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type Adder struct {
	engine.ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	addend float64
}

func (adder *Adder) Configure() error {

	data := AdderConfigData{}
	err := adder.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	adder.addend = data.Addend

	return nil
}

func (adder *Adder) OnInput(inJ string) {
	mlog.Printf("Adder OnInput: %s\n", inJ)

	inS := &AdderInputData{}
	err := inS.ReadFromJSON(inJ)
	if err != nil {
		panic(err)
	}

	outS := &AdderOutputData{}

	// the work
	{
		outS.Sum = inS.Value + adder.addend
	}

	outJ, err := outS.WriteToJSON()
	if err != nil {
		panic(err)
	}

	adder.Output <- outJ
}
