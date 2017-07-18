package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponent(t *testing.T) {
	assert := assert.New(t)

	// register our new component type
	Factory.Register("Foo", &foo{})

	// set up config data
	config := ArgMap{
		"i": 17,
	}

	// make the component object (and initialize it)
	x, err := Factory.Create("Foo", config)
	assert.NoError(err)
	assert.NotNil(x)
	c := x.(*foo)
	assert.NotNil(c)

	// did the config step work?
	assert.Equal(17, c.Config["i"])
	assert.Equal(17, c.i)

	// this setup is normally done by goflow itself
	chIn := make(chan string)
	chOut := make(chan string)
	c.Input = chIn
	c.Output = chOut

	// invoke OnLoad manually
	inJ := `{"X": 11}`
	go c.OnInput(inJ)

	// check the returned result
	outJ := <-chOut
	assert.JSONEq(`{"Y":28}`, outJ)
}

//---------------------------------------------------------------------

type fooConfigData struct{ I int }
type fooInputData struct{ X int }
type fooOutputData struct{ Y int }

func (m *fooInputData) ReadFromJSON(jsn string) error { return ReadFromJSON(jsn, m) }
func (m *fooOutputData) WriteToJSON() (string, error) { return WriteToJSON(m) }

type foo struct {
	// required
	ComponentCore

	// required: our component has one input port and one output port
	Input  <-chan string
	Output chan<- string

	// local storage
	i int
}

func (x *foo) Configure() error {
	// get the config data into a proper struct
	data := fooConfigData{}
	err := x.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	// do whatever processing of the config data we need to do
	x.i = data.I

	return nil
}

func (foo *foo) OnInput(inJ string) {

	// get the input into a proper input struct
	inS := &fooInputData{}
	err := inS.ReadFromJSON(inJ)
	if err != nil {
		panic(err)
	}

	// set up a proper output struct
	outS := fooOutputData{}

	// do the actual work
	outS.Y = inS.X + foo.i

	// push out the output
	outJ, err := outS.WriteToJSON()
	if err != nil {
		panic(err)
	}
	foo.Output <- outJ
}
