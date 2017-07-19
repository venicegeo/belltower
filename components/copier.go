package components

import (
	"fmt"

	"github.com/venicegeo/belltower/engine"
)

func init() {
	engine.Factory.Register("Copier", &Copier{})
}

type CopierConfigData struct{}

// Copier doesn't use an input or output data object, since it works with fields known only at runtime
//type CopierInputData struct {}
//type CopierOutputData struct {}

type Copier struct {
	engine.ComponentCore

	Input   <-chan string
	Output1 chan<- string
	Output2 chan<- string
}

func (copier *Copier) Configure() error {

	data := CopierConfigData{}
	err := copier.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	return nil
}

func (copier *Copier) OnInput(inJ string) {
	fmt.Printf("Copier OnInput: %s\n", inJ)

	copier.Output1 <- inJ
	copier.Output2 <- inJ
}
