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
	"fmt"

	"github.com/venicegeo/belltower/engine"
	"github.com/venicegeo/belltower/mpg/mlog"
)

func init() {
	engine.Factory.Register("Remapper", &Remapper{})
}

type RemapperConfigData struct {

	// remap field "key" to field "value"
	Remaps map[string]string
}

// Remapper doesn't use an input or output data object, since it works with fields known only at runtime
//type RemapperInputData struct {}
//type RemapperOutputData struct {}

type Remapper struct {
	engine.ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	remaps map[string]string
}

func (remapper *Remapper) Configure() error {

	remaps := RemapperConfigData{}
	err := remapper.Config.ToStruct(&remaps)
	if err != nil {
		return err
	}

	remapper.remaps = remaps.Remaps

	return nil
}

func (remapper *Remapper) OnInput(inJ string) {
	mlog.Printf("Remapper OnInput: %s\n", inJ)

	inMap, err := engine.NewArgMap(inJ)
	if err != nil {
		panic(err)
	}

	outMap, err := remapper.run(inMap)
	if err != nil {
		panic(err)
	}

	outJ, err := outMap.ToJSON()
	if err != nil {
		panic(err)
	}

	remapper.Output <- outJ
}

func (remapper *Remapper) run(inMap engine.ArgMap) (engine.ArgMap, error) {

	outMap := inMap

	for from, to := range remapper.remaps {

		fromValue, ok := inMap[from]
		if !ok {
			return nil, fmt.Errorf("field '%s' not found for remapping", from)
		}

		outMap[to] = fromValue
		delete(outMap, from)
	}

	return outMap, nil
}
