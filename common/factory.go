package common

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
