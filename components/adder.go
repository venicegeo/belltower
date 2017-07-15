package components

import (
	"log"

	"github.com/venicegeo/belltower/common"
)

func init() {
	Factory.Register("Adder", &Adder{})
}

// -- CONFIG --
//
// Addend int
//   The value added to the input. Default is zero.
//
// -- INPUT --
//
// (none)
//
// -- OUTPUT --
//
// Sum int
//   Value of input added to addend.

type Adder struct {
	ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	addend int
}

func (adder *Adder) localConfigure() error {

	addend, err := adder.config.GetIntOrDefault("addend", 0)
	if err != nil {
		return err
	}

	adder.addend = addend

	return nil
}

func (adder *Adder) OnInput(data string) {
	log.Printf("Adder: OnInput <%s>", data)

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

	input, err := in.GetInt("addend")
	if err != nil {
		return out, err
	}

	out["sum"] = input + adder.addend

	return out, nil
}
