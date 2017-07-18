package components

import (
	"fmt"

	"github.com/venicegeo/belltower/common"
)

func init() {
	common.Factory.Register("Adder", &Adder{})
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
	common.ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	addend float64
}

func (adder *Adder) Configure() error {

	data := AdderConfigData{}
	err := adder.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	adder.addend = data.Addend

	return nil
}

func (adder *Adder) OnInput(inJ string) {
	fmt.Printf("Adder OnInput: %s\n", inJ)

	inS := &AdderInputData{}
	err := inS.ReadFromJSON(inJ)
	if err != nil {
		panic(err)
	}

	outS := &AdderOutputData{}

	// the work
	{
		outS.Sum = inS.Value + adder.addend
	}

	outJ, err := outS.WriteToJSON()
	if err != nil {
		panic(err)
	}

	adder.Output <- outJ
}
