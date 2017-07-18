package components

import (
	"fmt"

	"github.com/trustmaster/goflow"
	"github.com/venicegeo/belltower/common"
)

type Component interface {
	// the system calls this to do init work specific to your component type
	// you must implement this
	Configure() error

	// the system calls this to do init work for ComponentCore fields
	// do not implement this yourself
	coreConfigure(config common.ArgMap) error

	// Implement "OnPORT" for each input port, using the right signature and contents. For example:
	//
	// func (adder *Adder) OnInput(inputJson string) {
	//   input := &AdderInputData{}
	//   err := input.ReadFromJSON(inputJson)
	//   if err != nil {...}
	//   output := &AdderOutputData{}
	//   output.x = ...
	//   outputJson, err := output.WriteToJSON()
	//   if err != nil {...}
	//   adder.Output <- outputJson
	// }
	//
	// See adder.go for a complete example.
}

type ComponentCore struct {
	config         common.ArgMap
	precondition   *common.Expression
	postcondition  *common.Expression
	executionCount int

	flow.Component
}

func (c *ComponentCore) coreConfigure(config common.ArgMap) error {
	// TODO: can this be replaced by Init() somehow?

	c.config = config

	c.executionCount = 0

	cond, err := config.GetStringOrDefault("precondition", "")
	if err != nil {
		return err
	}
	if cond != "" {
		e, err := common.NewExpression(cond, nil)
		if err != nil {
			return err
		}
		c.precondition = e
	}

	cond, err = config.GetStringOrDefault("postcondition", "")
	if err != nil {
		return nil
	}
	if cond != "" {
		e, err := common.NewExpression(cond, nil)
		if err != nil {
			return err
		}
		c.postcondition = e
	}

	return nil
}

//---------------------------------------------------------------------

func init() {
	Factory.Register("Starter", &Starter{})
	Factory.Register("Stopper", &Stopper{})
}

type Starter struct {
	ComponentCore
	Input  <-chan string
	Output chan<- string
}

func (*Starter) Configure() error { return nil }
func (s *Starter) OnInput(string) {
	fmt.Printf("Starter OnInput\n")
	s.Output <- "{}"
}

type Stopper struct {
	ComponentCore
	Input  <-chan string
	Output chan<- string
}

func (*Stopper) Configure() error { return nil }

func (s *Stopper) OnInput(string) {
	fmt.Printf("Stopper OnInput\n")
	s.Output <- "{}"
}
