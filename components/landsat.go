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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/venicegeo/belltower/engine"
	"github.com/venicegeo/belltower/mpg/mlog"
)

func init() {
	engine.Factory.Register("Landsat", &Landsat{})
}

type LandsatConfigData struct {
}

type LandsatInputData struct {
	SelectedImage string
}

func (m *LandsatInputData) Validate() error               { return nil } // TODO
func (m *LandsatInputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *LandsatInputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type LandsatOutputData struct {
	SelectedImage string
}

func (m *LandsatOutputData) Validate() error               { return nil } // TODO
func (m *LandsatOutputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *LandsatOutputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type Landsat struct {
	engine.ComponentCore

	Input  <-chan string
	Output chan<- string

	// config data
	domain     string
	auth       string
	planet_key string
	url        string
	auth64     string

	// local state
	addend float64
}

func (ls *Landsat) Configure() error {

	data := LandsatConfigData{}
	err := ls.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	err = ls.readBeachfrontrc()
	if err != nil {
		return err
	}

	ls.url = "https://bf-ia-broker." + ls.domain
	ls.auth64 = base64.StdEncoding.EncodeToString([]byte(ls.auth))

	return nil
}

func (ls *Landsat) readBeachfrontrc() error {
	user := os.Getenv("HOME")
	file, err := os.Open(user + "/.beachfrontrc")
	if err != nil {
		return err
	}

	byts, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	obj := &map[string]string{}
	err = json.Unmarshal(byts, obj)
	if err != nil {
		return err
	}

	f := func(field string) (string, error) {
		s, ok := (*obj)[field]
		if !ok || s == "" {
			return "", fmt.Errorf("Missing item in .beachfrontrc: '%s'", field)
		}
		return s, nil
	}

	ls.domain, err = f("domain")
	if err != nil {
		return err
	}
	ls.auth, err = f("auth")
	if err != nil {
		return err
	}
	ls.planet_key, err = f("planet_key")
	if err != nil {
		return err
	}

	return nil
}

func (ls *Landsat) OnInput(inJ string) {
	mlog.Printf("Landsat OnInput: %s\n", inJ)

	inS := &LandsatInputData{}
	err := inS.ReadFromJSON(inJ)
	if err != nil {
		panic(err)
	}

	outS := &LandsatOutputData{}

	// the work
	{
		path := "/planet/landsat/" + inS.SelectedImage + "?PL_API_KEY=" + ls.planet_key
		byts, status, err := ls.httpRequest("GET", ls.url+path, true)
		if err != nil {
			panic(err)
		}
		if status != 200 {
			panic(status)
		}

		//mlog.Debugf(string(byts))
		data := map[string]interface{}{}
		err = json.Unmarshal(byts, &data)
		if err != nil {
			panic(err)
		}

		//mlog.Debug(data)

		for _, color := range []string{"red", "green", "blue"} {

			url := data["properties"].(map[string]interface{})["bands"].(map[string]interface{})[color].(string)
			mlog.Debug(url)

			byts, status, err = ls.httpRequest("GET", url, false)
			if err != nil {
				panic(err)
			}
			if status != 200 {
				panic(status)
			}
			err = ioutil.WriteFile(inS.SelectedImage+"-"+color+".tif", byts, 0600)
			if err != nil {
				panic(err)
			}
		}
	}

	outS.SelectedImage = inS.SelectedImage

	outJ, err := outS.WriteToJSON()
	if err != nil {
		panic(err)
	}

	ls.Output <- outJ
}

func (ls *Landsat) httpRequest(verb string, urlPath string, useAuth bool) ([]byte, int, error) {

	client := &http.Client{}
	req, err := http.NewRequest(verb, urlPath, nil)
	if err != nil {
		return nil, 0, err
	}

	if useAuth {
		req.Header.Set("Authorization", "Basic "+ls.auth64)
	}

	//mlog.Debugf("req: %#v", req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	status := resp.StatusCode

	byts, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, 0, err
	}
	//mlog.Debugf("%s", body)

	return byts, status, nil
}
