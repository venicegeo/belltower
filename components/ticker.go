package components

import (
	"fmt"
	"log"

	"github.com/venicegeo/belltower/common"
)

func init() {
	Factory.Register("Ticker", &Ticker{})
}

type Ticker struct {
	ComponentCore

	Input  <-chan string
	Output chan<- string

	max int
	cur int
}

func (ticker *Ticker) Description() *Description {
	configFields := []*common.DataType{
		common.NewScalarDataType("max", common.TypeNameInt),
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
		Name: "ticker",
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

func (ticker *Ticker) localConfigure() error {
	if !ticker.config.IsInt("max") {
		return fmt.Errorf("required ticker field missing: max")
	}
	ticker.max = ticker.config.AsInt("max")

	return nil
}

func (ticker *Ticker) OnInput(data string) {
	log.Printf("OnInput ticker!")
	ticker.Output <- "baz"
}

func (ticker *Ticker) Run(in common.Map) (common.Map, error) {

	out := common.Map{}
	out["y"] = ticker.cur

	ticker.cur++

	return out, nil
}
