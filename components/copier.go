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
)

func init() {
	engine.Factory.Register("Copier", &Copier{})
}

type CopierConfigData struct{}

// Copier doesn't use an input or output data object, since it works with fields known only at runtime
//type CopierInputData struct {}
//type CopierOutputData struct {}

type Copier struct {
	engine.ComponentCore

	Input   <-chan string
	Output1 chan<- string
	Output2 chan<- string
}

func (copier *Copier) Configure() error {

	data := CopierConfigData{}
	err := copier.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	return nil
}

func (copier *Copier) OnInput(inJ string) {
	fmt.Printf("Copier OnInput: %s\n", inJ)

	copier.Output1 <- inJ
	copier.Output2 <- inJ
}
