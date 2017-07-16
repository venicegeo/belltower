package components

import (
	"fmt"

	"encoding/json"

	"github.com/venicegeo/belltower/common"
)

func init() {
	Factory.Register("Adder", &Adder{})
}

type AdderConfigData struct {

	// The value added to the input. Default is zero.
	Addend float64
}

type AdderInputData struct {

	// The value added to the addend from the configuration. Default is zero.
	Value float64
}

type AdderOutputData struct {

	// Value of input value added to addend.
	Sum float64
}

type Adder struct {
	ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	addend float64
}

func (adder *Adder) Configure() error {

	data := AdderConfigData{}
	err := adder.config.ToStruct(&data)
	if err != nil {
		return err
	}

	adder.addend = data.Addend

	return nil
}

func (adder *Adder) OnInput(data string) {
	fmt.Printf("Adder OnInput: %s\n", data)

	inputMap := common.ArgMap{}
	err := json.Unmarshal([]byte(data), &inputMap)
	if err != nil {
		panic(err)
	}

	input := AdderInputData{}
	_, err = common.SetStructFromMap(inputMap, &input, true)
	if err != nil {
		panic(err)
	}

	output, err := adder.Run(input)
	if err != nil {
		panic(err)
	}

	buf, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	adder.Output <- string(buf)
}

func (adder *Adder) Run(in interface{}) (interface{}, error) {

	input := in.(AdderInputData)
	output := AdderOutputData{}

	output.Sum = input.Value + adder.addend

	return output, nil
}
