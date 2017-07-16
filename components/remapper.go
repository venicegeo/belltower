package components

import (
	"encoding/json"
	"fmt"

	"github.com/venicegeo/belltower/common"
)

func init() {
	Factory.Register("Remapper", &Remapper{})
}

type RemapperConfigData struct {

	// remap field "key" to field "value"
	Remaps map[string]string
}

// Remapper doesn't use input or output data object, since it works with fields known only at runtime
//type RemapperInputData struct {}
//type RemapperOutputData struct {}

type Remapper struct {
	ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	remaps map[string]string
}

func (remapper *Remapper) Configure() error {

	remaps := RemapperConfigData{}
	err := remapper.config.ToStruct(&remaps)
	if err != nil {
		return err
	}

	remapper.remaps = remaps.Remaps

	return nil
}

func (remapper *Remapper) OnInput(data string) {
	fmt.Printf("Remapper OnInput: %s\n", data)

	inputMap, err := common.NewArgMap(data)
	if err != nil {
		panic(err)
	}

	outputMap, err := remapper.Run(inputMap)
	if err != nil {
		panic(err)
	}

	buf, err := json.Marshal(outputMap)
	if err != nil {
		panic(err)
	}

	remapper.Output <- string(buf)
}

func (remapper *Remapper) Run(in interface{}) (interface{}, error) {

	inputMap := in.(common.ArgMap)

	outputMap := inputMap

	for from, to := range remapper.remaps {

		fromValue, ok := inputMap[from]
		if !ok {
			return nil, fmt.Errorf("field '%s' not found for remapping", from)
		}

		outputMap[to] = fromValue
		delete(outputMap, from)
	}

	return outputMap, nil
}
