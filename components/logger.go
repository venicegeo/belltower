package components

import (
	"bytes"
	"fmt"
	"os"

	"encoding/json"

	"github.com/venicegeo/belltower/common"
)

func init() {
	common.Factory.Register("Logger", &Logger{})
}

type LoggerConfigData struct {

	// File to write to
	// TODO: eventually shove the data off to a URL somewhere or something
	FileName string
}

// TODO: rather than log the whole actual JSON, would be nice to use JSON to describe the log message text
type Logger struct {
	common.ComponentCore

	Input  <-chan string
	Output chan<- string

	filename string
}

func (logger *Logger) Configure() error {

	data := LoggerConfigData{}
	err := logger.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	logger.filename = data.FileName

	return nil
}

func (logger *Logger) OnInput(inJ string) {
	fmt.Printf("Logger OnInput: %s\n", inJ)

	dst := new(bytes.Buffer)

	src := []byte(inJ)

	err := json.Compact(dst, src)
	if err != nil {
		panic(err)
	}

	var f *os.File
	switch logger.filename {
	case "STDOUT":
		f = os.Stdout
	case "STDERR":
		f = os.Stderr
	default:
		f, err := os.OpenFile(logger.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}
	fmt.Fprintf(f, "%s\n", dst)

	logger.Output <- "{}"
}