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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/venicegeo/belltower/engine"
	"github.com/venicegeo/belltower/mpg/mlog"
)

func init() {
	engine.Factory.Register("Beachfront", &Beachfront{})
}

type BeachfrontConfigData struct {
	Path string
}

type BeachfrontInputData struct {
	SelectedImage string
}

func (m *BeachfrontInputData) Validate() error               { return nil } // TODO
func (m *BeachfrontInputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *BeachfrontInputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type BeachfrontOutputData struct {
	SelectedImage string
}

func (m *BeachfrontOutputData) Validate() error               { return nil } // TODO
func (m *BeachfrontOutputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *BeachfrontOutputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type Beachfront struct {
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
	path   string
}

func (bf *Beachfront) Configure() error {

	data := BeachfrontConfigData{}
	err := bf.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	err = bf.readBeachfrontrc()
	if err != nil {
		return err
	}

	bf.url = "https://bf-api." + bf.domain
	bf.auth64 = base64.StdEncoding.EncodeToString([]byte(bf.auth))

	bf.path = data.Path

	return nil
}

func (bf *Beachfront) readBeachfrontrc() error {
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

	bf.domain, err = f("domain")
	if err != nil {
		return err
	}
	bf.auth, err = f("auth")
	if err != nil {
		return err
	}
	bf.planet_key, err = f("planet_key")
	if err != nil {
		return err
	}

	return nil
}

func (bf *Beachfront) OnInput(inJ string) {

	inS := &BeachfrontInputData{}
	err := inS.ReadFromJSON(inJ)
	if err != nil {
		panic(err)
	}

	mlog.Printf("Beachfront OnInput: %s\n", inS)

	outS := &BeachfrontOutputData{}

	// the work
	{
		serviceId, err := bf.getServiceId()
		if err != nil {
			panic(err)
		}

		mlog.Printf("xx serviceId: %s", serviceId)

		jobId, err := bf.submitJob(serviceId, inS.SelectedImage)
		if err != nil {
			panic(err)
		}

		mlog.Printf("xx jobId: %s", jobId)

		err = bf.waitForSuccess(jobId)
		if err != nil {
			panic(err)
		}

		mlog.Printf("Beachfront downloading %s / geojson\n", inS.SelectedImage)
		err = bf.getFiles(jobId, inS.SelectedImage)
		if err != nil {
			panic(err)
		}
	}

	outS.SelectedImage = inS.SelectedImage

	outJ, err := outS.WriteToJSON()
	if err != nil {
		panic(err)
	}

	bf.Output <- outJ
}

func (bf *Beachfront) getFiles(jobId string, image string) error {
	body, code, err := bf.httpRequest("GET", "/v0/job/"+jobId, "")
	if err != nil {
		return err
	}
	byts := []byte(body)

	if code != 200 {
		return fmt.Errorf("GetFiles expected %d, got %d", 200, code)
	}

	err = ioutil.WriteFile(bf.path+"/"+image+".json", byts, 0600)
	if err != nil {
		panic(err)
	}

	body, code, err = bf.httpRequest("GET", "/v0/job/"+jobId+".geojson", "")
	if err != nil {
		return err
	}
	byts = []byte(body)

	if code != 200 {
		return fmt.Errorf("GetFiles expected %d, got %d", 200, code)
	}

	err = ioutil.WriteFile(bf.path+"/"+image+".geojson", byts, 0600)
	if err != nil {
		panic(err)
	}

	return nil
}

func (bf *Beachfront) waitForSuccess(jobId string) error {

	var err error
	status := ""

	const delay = 15
	const iters = 40

	for i := 0; i < iters; i++ {
		status, err = bf.getJobStatus(jobId)
		if err != nil {
			panic(err)
		}

		mlog.Printf("xx jobStatus: %s (%d seconds)", status, i*delay)

		if status == "Success" {
			return nil
		}

		time.Sleep(time.Duration(delay) * time.Second)
	}

	return fmt.Errorf("Job did not complete: %s status after %d seconds", status, iters*delay)
}

func (bf *Beachfront) getJobStatus(jobId string) (string, error) {
	body, code, err := bf.httpRequest("GET", "/v0/job/"+jobId, "")
	if err != nil {
		return "", err
	}

	if code != 200 {
		return "", fmt.Errorf("GetJobStatus expected %d, got %d", 200, code)
	}

	//mlog.Debugf("getJobStatus returned: %s", body)

	rt := respType{}
	err = json.Unmarshal([]byte(body), &rt)
	if err != nil {
		return "", err
	}

	status := rt.Job.Properties.Status
	//mlog.Debugf("getJobStatus returned: %s", status)

	return status, nil
}

type jobProperties struct {
	Status string `json:"status"`
}

type jobType struct {
	ID         string        `json:"id"`
	Properties jobProperties `json:"properties"`
}

type respType struct {
	Job jobType `json:"job"`
}

func (bf *Beachfront) getServiceId() (string, error) {
	body, status, err := bf.httpRequest("GET", "/v0/algorithm", "")
	if err != nil {
		return "", err
	}

	if status != 200 {
		return "", fmt.Errorf("GetServiceId expected %d, got %d", 200, status)
	}

	//mlog.Debugf("getServiceId returned: %s", body)

	type algType struct {
		Name      string `json:"name"`
		ServiceId string `json:"service_id"`
	}
	type algsType struct {
		Algorithms []algType `json:"algorithms"`
	}
	algs := algsType{}
	err = json.Unmarshal([]byte(body), &algs)
	if err != nil {
		return "", err
	}

	//mlog.Debugf("getServiceId returned: %v", data)

	for _, alg := range algs.Algorithms {
		if alg.Name == "NDWI_PY" { // TODO: don't hardcode this
			return alg.ServiceId, nil
		}
	}

	return "", fmt.Errorf("getServiceId: algorithm not found")
}

func (bf *Beachfront) submitJob(serviceId string, selectedImage string) (string, error) {
	payload := `{
		"algorithm_id": "` + serviceId + `",
		"scene_id": "landsat:` + selectedImage + `",
		"name": "xyzzy",
		"planet_api_key": "` + bf.planet_key + `"
	}`

	//mlog.Debugf("submitJob payload: %s", payload)

	respBody, status, err := bf.httpRequest("POST", "/v0/job", payload)
	if err != nil {
		panic(err)
	}
	if status != 201 {
		panic(status)
	}

	//mlog.Debugf("submitJob response: (%d) %s", status, respBody)

	rt := respType{}
	err = json.Unmarshal([]byte(respBody), &rt)
	if err != nil {
		return "", err
	}

	return rt.Job.ID, err
}

func (bf *Beachfront) httpRequest(verb string, urlPath string, reqBody string) (string, int, error) {
	reqBodyBuf := bytes.NewBufferString(reqBody)

	client := &http.Client{Timeout: time.Duration(30 * time.Minute)}
	req, err := http.NewRequest(verb, bf.url+urlPath, reqBodyBuf)
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Authorization", "Basic "+bf.auth64)
	req.Header.Set("Content-Type", "application/json")

	//mlog.Debugf("req: %#v", req)
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}

	status := resp.StatusCode

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", 0, err
	}
	//mlog.Debugf("%s", body)

	return string(body), status, nil
}
