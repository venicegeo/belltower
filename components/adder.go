package components

import (
	"fmt"

	"github.com/venicegeo/belltower/common"
)

func init() {
	Factory.Register("Adder", &Adder{})
}

type AdderConfigData struct {

	// The value added to the input. Default is zero.
	Addend float64
}

// implements Serializer
type AdderInputData struct {

	// The value added to the addend from the configuration. Default is zero.
	Value float64
}

func (m *AdderInputData) Validate() error               { return nil } // TODO
func (m *AdderInputData) ReadFromJSON(jsn string) error { return common.ReadFromJSON(jsn, m) }
func (m *AdderInputData) WriteToJSON() (string, error)  { return common.WriteToJSON(m) }

// implements Serializer
type AdderOutputData struct {

	// Value of input value added to addend.
	Sum float64
}

func (m *AdderOutputData) Validate() error               { return nil } // TODO
func (m *AdderOutputData) ReadFromJSON(jsn string) error { return common.ReadFromJSON(jsn, m) }
func (m *AdderOutputData) WriteToJSON() (string, error)  { return common.WriteToJSON(m) }

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

	input := &AdderInputData{}
	err := input.ReadFromJSON(inputJson)
	if err != nil {
		panic(err)
	}

	output, err := adder.Run(input)
	if err != nil {
		panic(err)
	}

	outputJson, err := output.(*AdderOutputData).WriteToJSON()
	if err != nil {
		panic(err)
	}

	adder.Output <- outputJson
}

func (adder *Adder) Run(in interface{}) (interface{}, error) {

	input := in.(AdderInputData)
	output := &AdderOutputData{}

	output.Sum = input.Value + adder.addend

	return output, nil
}
