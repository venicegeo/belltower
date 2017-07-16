package components

import (
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

func (remapper *Remapper) OnInput(inputJson string) {
	fmt.Printf("Remapper OnInput: %s\n", inputJson)

	inputMap, err := common.NewArgMap(inputJson)
	if err != nil {
		panic(err)
	}

	outputMap, err := remapper.Run(inputMap)
	if err != nil {
		panic(err)
	}

	outputJson, err := outputMap.(common.ArgMap).ToJSON()
	if err != nil {
		panic(err)
	}

	remapper.Output <- outputJson
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
