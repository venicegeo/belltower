package components

import (
	"fmt"

	"github.com/venicegeo/belltower/common"
)

func init() {
	Factory.Register("Adder", &Adder{})
}

// -- CONFIG --
//
// addend int
//   The value added to the input. Default is zero.
//
// -- INPUT --
//
// value int
//   The value added to the addend from the configuration. Default is zero.
//
// -- OUTPUT --
//
// sum int
//   Value of input value added to addend.

type Adder struct {
	ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	addend float64
}

func (adder *Adder) localConfigure() error {

	addend, err := adder.config.GetFloatOrDefault("addend", 0.0)
	if err != nil {
		return err
	}

	adder.addend = addend

	return nil
}

func (adder *Adder) OnInput(data string) {
	fmt.Printf("Adder OnInput: %s\n", data)

	in, err := common.NewArgMap(data)
	if err != nil {
		panic(err)
	}

	out, err := adder.Run(in)
	if err != nil {
		panic(err)
	}

	s, err := out.ToJSON()
	if err != nil {
		panic(err)
	}

	adder.Output <- s
}

func (adder *Adder) Run(in common.ArgMap) (common.ArgMap, error) {

	out := common.ArgMap{}

	value, err := in.GetFloat("value")
	if err != nil {
		return out, err
	}

	out["sum"] = value + adder.addend

	return out, nil
}
