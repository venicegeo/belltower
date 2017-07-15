package components

import (
	"fmt"

	"github.com/venicegeo/belltower/common"
)

func init() {
	Factory.Register("Remapper", &Remapper{})
}

// -- CONFIG --
//
// remapped []{from, to string}
//   The value added to the input. Default is zero.
//
// -- INPUT --
//
// anything
//
// -- OUTPUT --
//
// same as input, but remapped

type RemapPair struct {
	from string
	to   string
}

type Remapper struct {
	ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	remaps map[string]string
}

func (remapper *Remapper) localConfigure() error {

	remaps, err := remapper.config.GetInterfaceOrDefault("remaps", nil)
	if err != nil {
		return err
	}
	if remaps == nil {
		return fmt.Errorf("remaps list cannot be empty")
	}

	remapper.remaps = remaps.(map[string]string)

	return nil
}

func (remapper *Remapper) OnInput(data string) {
	fmt.Printf("Remapper OnInput: %s\n", data)

	in, err := common.NewArgMap(data)
	if err != nil {
		panic(err)
	}

	out, err := remapper.Run(in)
	if err != nil {
		panic(err)
	}

	s, err := out.ToJSON()
	if err != nil {
		panic(err)
	}

	remapper.Output <- s
}

func (remapper *Remapper) Run(in common.ArgMap) (common.ArgMap, error) {

	out := in

	for from, to := range remapper.remaps {

		fromValue, ok := in[from]
		if !ok {
			return nil, fmt.Errorf("field '%s' not found for remapping", from)
		}

		out[to] = fromValue
		delete(out, from)
	}

	return out, nil
}
