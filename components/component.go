package components

import (
	"encoding/json"
	"fmt"

	"github.com/trustmaster/goflow"
	"github.com/venicegeo/belltower/common"
)

type Component interface {
	//// describes the port datatypes, etc
	//Description() *Description

	// configuration specific to your component type
	// only called by ComponentCore.configure()
	// you must implement
	Configure() error

	// perform one execution
	Run(interface{} /*in common.ArgMap*/) (interface{} /*common.ArgMap*/, error)

	// called by Factory to do init work for ComponentCore fields
	// do not implement yourself
	coreConfigure(config common.ArgMap) error
}

type ComponentCore struct {
	config         common.ArgMap
	precondition   *common.Expression
	postcondition  *common.Expression
	executionCount int

	flow.Component
}

func (c *ComponentCore) coreConfigure(config common.ArgMap) error {

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

type Description struct {
	Id       common.Id        `json:"id"`
	Name     string           `json:"name"`
	Metadata *common.Metadata `json:"metadata,omitempty"`

	Config *common.Port `json:"config,omitempty"`
	Input  *common.Port `json:"input,omitempty"`
	Output *common.Port `json:"output,omitempty"`
}

func FromJSONToStruct(jsn string, obj interface{}) error {

	m := common.ArgMap{}
	err := json.Unmarshal([]byte(jsn), &m)
	if err != nil {
		return err
	}

	err = m.ToStruct(&obj)
	if err != nil {
		return err
	}

	return nil
}

func FromStructToJSON(obj interface{}) (string, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

//---------------------------------------------------------------------

func init() {
	Factory.Register("Starter", &Starter{})
}

type StarterConfigData struct{}

type StarterInputData struct{}

type StarterOutputData struct{}

type Starter struct {
	ComponentCore
	Input  <-chan string
	Output chan<- string
}

func (*Starter) Configure() error {
	return nil
}

func (s *Starter) OnInput(string) {
	fmt.Printf("Starter OnInput\n")
	s.Output <- "{}"
}

func (*Starter) Run(in interface{}) (interface{}, error) { return nil, nil }

//---------------------------------------------------------------------

func init() {
	Factory.Register("Stopper", &Stopper{})
}

type StopperConfigData struct{}

type StopperInputData struct{}

type StopperOutputData struct{}

type Stopper struct {
	ComponentCore
	Input  <-chan string
	Output chan<- string
}

func (*Stopper) Configure() error {
	return nil
}

func (s *Stopper) OnInput(string) {
	fmt.Printf("Stopper OnInput\n")
	s.Output <- "{}"
}

func (*Stopper) Run(in interface{}) (interface{}, error) { return nil, nil }
