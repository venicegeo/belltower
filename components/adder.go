package components

import (
	"fmt"
	"log"

	"github.com/venicegeo/belltower/common"
)

func init() {
	Factory.Register("Adder", &Adder{})
}

type Adder struct {
	ComponentCore

	addend int

	Input  <-chan string
	Output chan<- string
}

func (adder *Adder) Description() *Description {
	configFields := []*common.DataType{
		common.NewScalarDataType("addend", common.TypeNameInt),
	}
	config := common.NewStructDataType("config", configFields)

	inFields := []*common.DataType{
		common.NewScalarDataType("x", common.TypeNameInt),
	}
	in := common.NewStructDataType("Input", inFields)

	outFields := []*common.DataType{
		common.NewScalarDataType("y", common.TypeNameInt),
	}
	out := common.NewStructDataType("Output", outFields)

	return &Description{
		Name: "adder",
		Config: &common.Port{
			DataType: config,
		},
		Input: &common.Port{
			DataType: in,
		},
		Output: &common.Port{
			DataType: out,
		},
	}
}

func (adder *Adder) localConfigure() error {
	if !adder.config.IsInt("addend") {
		return fmt.Errorf("required ticker field missing: addend")
	}
	adder.addend = adder.config.AsInt("addend")

	return nil
}

func (adder *Adder) OnInput(data string) {
	log.Printf("OnInput!")
	adder.Output <- "yow"
}

func (adder *Adder) Run(in common.Map) (common.Map, error) {
	x := in["x"].(int)

	y := x + adder.addend

	out := common.Map{}
	out["y"] = y

	return out, nil
}
