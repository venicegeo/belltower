package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMySender(t *testing.T) {
	assert := assert.New(t)

	// register our new component type
	Factory.Register("MySender", &MySender{})

	// set up config data
	config := ArgMap{
		"i": 17,
	}

	// make the component object (and initialize it)
	mysenderX, err := Factory.Create("MySender", config)
	assert.NoError(err)
	assert.NotNil(mysenderX)
	mysender := mysenderX.(*MySender)
	assert.NotNil(mysender)

	// did the config step work?
	assert.Equal(17, mysender.Config["i"])
	assert.Equal(17, mysender.i)

	// this setup is normally done by goflow itself
	chIn := make(chan string)
	chOut := make(chan string)
	mysender.Input = chIn
	mysender.Output = chOut

	// invoke OnLoad manually
	inJ := "{}"
	go mysender.OnInput(inJ)

	// check the returned result
	outJ := <-chOut
	assert.JSONEq(`{"I":17}`, outJ)
}

func TestMyReceiver(t *testing.T) {
	assert := assert.New(t)

	Factory.Register("MyReceiver", &MyReceiver{})

	myreceiverX, err := Factory.Create("MyReceiver", nil)
	assert.NoError(err)
	assert.NotNil(myreceiverX)
	myreceiver := myreceiverX.(*MyReceiver)
	assert.NotNil(myreceiver)

	chIn := make(chan string)
	chOut := make(chan string)
	myreceiver.Input = chIn
	myreceiver.Output = chOut

	inJ := `{"I": 11}`
	go myreceiver.OnInput(inJ)

	_ = <-chOut

	assert.Equal(11, myreceiver.i)
}

//---------------------------------------------------------------------

type MySenderConfigData struct{ I int }
type MySenderOutputData struct{ I int }

func (m *MySenderOutputData) WriteToJSON() (string, error) { return WriteToJSON(m) }

type MySender struct {
	// required
	ComponentCore

	// required: our component has one input port and one output port
	Input  <-chan string
	Output chan<- string

	// local storage
	i int
}

func (mysender *MySender) Configure() error {
	// get the config data into a proper struct
	data := MySenderConfigData{}
	err := mysender.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	// do whatever processing of the config data we need to do
	mysender.i = data.I

	return nil
}

func (mysender *MySender) OnInput(_ string) {

	// set up a proper output struct
	outS := MySenderOutputData{}

	// do the actual work
	outS.I = mysender.i

	// push out the output
	outJ, err := outS.WriteToJSON()
	if err != nil {
		panic(err)
	}
	mysender.Output <- outJ
}

//---------------------------------------------------------------------

type MyReceiverInputData struct{ I int }

func (m *MyReceiverInputData) ReadFromJSON(jsn string) error { return ReadFromJSON(jsn, m) }

type MyReceiver struct {
	// required
	ComponentCore

	// required: our component has one input port and one output port
	Input  <-chan string
	Output chan<- string

	// local state, for testing
	i int
}

func (myreceiver *MyReceiver) Configure() error { return nil }

func (myreceiver *MyReceiver) OnInput(inJ string) {

	// get the input into a proper input struct
	inS := &MyReceiverInputData{}
	err := inS.ReadFromJSON(inJ)
	if err != nil {
		panic(err)
	}

	myreceiver.i = inS.I

	myreceiver.Output <- "{}"
}

//---------------------------------------------------------------------

func init() {
	Factory.Register("MyCopier", &MyCopier{})
}

type MyCopierConfigData struct{}

type MyCopier struct {
	ComponentCore

	Input   <-chan string
	Output1 chan<- string
	Output2 chan<- string
}

func (mycopier *MyCopier) Configure() error { return nil }

func (mycopier *MyCopier) OnInput(inJ string) {
	//fmt.Printf("Copier OnInput: %s\n", inJ)
	mycopier.Output1 <- inJ
	mycopier.Output2 <- inJ
}

//---------------------------------------------------------------------

func init() {
	Factory.Register("MyAdder", &MyAdder{})
}

type MyAdderConfigData struct {

	// The value added to the input. Default is zero.
	Addend float64
}

// implements Serializer
type MyAdderInputData struct {

	// The value added to the addend from the configuration. Default is zero.
	Value float64
}

func (m *MyAdderInputData) Validate() error               { return nil } // TODO
func (m *MyAdderInputData) ReadFromJSON(jsn string) error { return ReadFromJSON(jsn, m) }
func (m *MyAdderInputData) WriteToJSON() (string, error)  { return WriteToJSON(m) }

// implements Serializer
type MyAdderOutputData struct {

	// Value of input value added to addend.
	Sum float64
}

func (m *MyAdderOutputData) Validate() error               { return nil } // TODO
func (m *MyAdderOutputData) ReadFromJSON(jsn string) error { return ReadFromJSON(jsn, m) }
func (m *MyAdderOutputData) WriteToJSON() (string, error)  { return WriteToJSON(m) }

type MyAdder struct {
	ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	addend float64
}

func (myadder *MyAdder) Configure() error {

	data := MyAdderConfigData{}
	err := myadder.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	myadder.addend = data.Addend

	return nil
}

func (myadder *MyAdder) OnInput(inJ string) {
	fmt.Printf("MyAdder OnInput: %s\n", inJ)

	inS := &MyAdderInputData{}
	err := inS.ReadFromJSON(inJ)
	if err != nil {
		panic(err)
	}

	outS := &MyAdderOutputData{}

	// the work
	{
		outS.Sum = inS.Value + myadder.addend
	}

	outJ, err := outS.WriteToJSON()
	if err != nil {
		panic(err)
	}

	myadder.Output <- outJ
}
