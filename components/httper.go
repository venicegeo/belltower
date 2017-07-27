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
	"io/ioutil"
	"net/http"

	"bytes"

	"github.com/venicegeo/belltower/engine"
	"github.com/venicegeo/belltower/mpg/mlog"
)

func init() {
	engine.Factory.Register("HTTPer", &HTTPer{})
}

type HTTPerConfigData struct {

	// The value added to the input. Default is zero.
	Verb string
	Url  string
}

type HTTPerInputData struct {

	// The value added to the addend from the configuration. Default is zero.
	RequestBody string
}

func (m *HTTPerInputData) Validate() error               { return nil } // TODO
func (m *HTTPerInputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *HTTPerInputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type HTTPerOutputData struct {

	// Value of input value added to addend.
	ResponseBody string
}

func (m *HTTPerOutputData) Validate() error               { return nil } // TODO
func (m *HTTPerOutputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *HTTPerOutputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type HTTPer struct {
	engine.ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	verb string
	url  string
}

func (h *HTTPer) Configure() error {

	data := HTTPerConfigData{}
	err := h.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	h.verb = data.Verb
	h.url = data.Url

	return nil
}

func (h *HTTPer) OnInput(inJ string) {
	mlog.Printf("HTTPer OnInput: %s\n", inJ)

	inS := &HTTPerInputData{}
	err := inS.ReadFromJSON(inJ)
	if err != nil {
		panic(err)
	}

	outS := &HTTPerOutputData{}

	// the work
	{
		respBody, err := httpRequest(h.verb, h.url, inS.RequestBody)
		if err != nil {
			panic(err)
		}

		outS.ResponseBody = respBody
	}

	outJ, err := outS.WriteToJSON()
	if err != nil {
		panic(err)
	}

	h.Output <- outJ
}

func httpRequest(verb string, url string, reqBody string) (string, error) {
	buf := bytes.NewBufferString(reqBody)

	client := &http.Client{}
	req, err := http.NewRequest(verb, url, buf)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	mlog.Debugf("%s", body)

	return string(body), nil
}
