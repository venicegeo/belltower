package components

import (
	"fmt"
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

func (adder *Adder) OnInput(inputJson string) {
	fmt.Printf("Adder OnInput: %s\n", inputJson)

	input := AdderInputData{}
	err := FromJSONToStruct(inputJson, &input)
	if err != nil {
		panic(err)
	}

	output, err := adder.Run(input)
	if err != nil {
		panic(err)
	}

	outputJson, err := FromStructToJSON(output)
	if err != nil {
		panic(err)
	}

	adder.Output <- outputJson
}

func (adder *Adder) Run(in interface{}) (interface{}, error) {

	input := in.(AdderInputData)
	output := AdderOutputData{}

	output.Sum = input.Value + adder.addend

	return output, nil
}
