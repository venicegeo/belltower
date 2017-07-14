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
	config := []*common.DataType{
		common.NewScalarDataType("max", common.TypeNameInt),
	}

	input := []*common.DataType{
		common.NewScalarDataType("x", common.TypeNameInt),
	}

	output := []*common.DataType{
		common.NewScalarDataType("y", common.TypeNameInt),
	}

	return &Description{
		Name: "ticker",
		Config: &common.Port{
			DataType: common.NewStructDataType("Config", config),
		},
		Input: &common.Port{
			DataType: common.NewStructDataType("Input", input),
		},
		Output: &common.Port{
			DataType: common.NewStructDataType("Output", output),
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
	log.Printf("Ticker: OnInput <%s>", data)
	ticker.Output <- "ticker" + data + "baz"
}

func (ticker *Ticker) Run(in common.Map) (common.Map, error) {

	out := common.Map{}
	out["y"] = ticker.cur

	ticker.cur++

	return out, nil
}
