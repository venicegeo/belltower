package components

import (
	"github.com/trustmaster/goflow"
	"github.com/venicegeo/belltower/common"
)

type Component interface {
	//// describes the port datatypes, etc
	//Description() *Description

	// configuration specific to your component type
	// only called by ComponentCore.configure()
	// you must implement
	localConfigure() error

	// perform one execution
	Run(in common.ArgMap) (out common.ArgMap, err error)

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
