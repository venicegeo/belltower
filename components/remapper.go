package components

import (
	"fmt"

	"github.com/venicegeo/belltower/common"
)

func init() {
	common.Factory.Register("Remapper", &Remapper{})
}

type RemapperConfigData struct {

	// remap field "key" to field "value"
	Remaps map[string]string
}

// Remapper doesn't use an input or output data object, since it works with fields known only at runtime
//type RemapperInputData struct {}
//type RemapperOutputData struct {}

type Remapper struct {
	common.ComponentCore

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
	fmt.Printf("Remapper OnInput: %s\n", inJ)

	inMap, err := common.NewArgMap(inJ)
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

func (remapper *Remapper) run(inMap common.ArgMap) (common.ArgMap, error) {

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
