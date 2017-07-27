/* Copyright 2017, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package components

import (
	"bytes"
	"fmt"
	"os"

	"encoding/json"

	"github.com/venicegeo/belltower/engine"
	"github.com/venicegeo/belltower/mpg/mlog"
)

func init() {
	engine.Factory.Register("Logger", &Logger{})
}

type LoggerConfigData struct {

	// File to write to
	// TODO: eventually shove the data off to a URL somewhere or something
	FileName string
}

// TODO: rather than log the whole actual JSON, would be nice to use JSON to describe the log message text
type Logger struct {
	engine.ComponentCore

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
	mlog.Printf("Logger OnInput: %s\n", inJ)

	dst := new(bytes.Buffer)

	src := []byte(inJ)

	err := json.Compact(dst, src)
	if err != nil {
		panic(err)
	}

	var file *os.File
	p := func() { fmt.Fprintf(file, "%s\n", dst) }

	switch logger.filename {
	case "STDOUT":
		file = os.Stdout

	case "STDERR":
		file = os.Stderr
	default:
		file, err = os.OpenFile(logger.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	p()

	logger.Output <- "{}"
}
