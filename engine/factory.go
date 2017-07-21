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
)

type TypeFactory struct {
	types map[string]Component
}

var Factory = &TypeFactory{
	types: map[string]Component{},
}

func (f *TypeFactory) Register(name string, typ Component) {
	f.types[name] = typ
}

func (f *TypeFactory) Create(typ string, config ArgMap) (Component, error) {

	dummy := f.types[typ]
	if dummy == nil {
		return nil, fmt.Errorf("component factory: invalid name: %s", typ)
	}

	dummy2 := NewViaReflection(dummy)
	if dummy2 == nil {
		return nil, fmt.Errorf("component factory: unable to create: %s", typ)
	}

	component, ok := dummy2.(Component)
	if !ok || component == nil {
		return nil, fmt.Errorf("component factory: really unable to create: %s", typ)
	}

	err := component.coreConfigure(config)
	if err != nil {
		return nil, err
	}

	err = component.Configure()
	if err != nil {
		return nil, err
	}

	return component, nil
}
