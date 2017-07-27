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
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/venicegeo/belltower/engine"
	"github.com/venicegeo/belltower/mpg/mlog"
)

func init() {
	engine.Factory.Register("Rasterize", &Rasterize{})
}

type RasterizeConfigData struct {
	Path string
}

type RasterizeInputData struct {
	SelectedImage string
}

func (m *RasterizeInputData) Validate() error               { return nil } // TODO
func (m *RasterizeInputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *RasterizeInputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type RasterizeOutputData struct {
	SelectedImage string
}

func (m *RasterizeOutputData) Validate() error               { return nil } // TODO
func (m *RasterizeOutputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *RasterizeOutputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type Rasterize struct {
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

func (r *Rasterize) Configure() error {

	data := RasterizeConfigData{}
	err := r.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	r.path = data.Path
	return nil
}

func (r *Rasterize) OnInput(inJ string) {
	fmt.Printf("Rasterize OnInput: %s\n", inJ)

	inS := &RasterizeInputData{}
	err := inS.ReadFromJSON(inJ)
	if err != nil {
		panic(err)
	}

	outS := &RasterizeOutputData{}

	// the work
	{
		//rm -f rgb.tif
		//gdal_merge.py -separate -o rgb.tif red.tif green.tif blue.tif
		//gdal_rasterize -b 1 -b 2 -b 3 -burn 65535 -burn 65535 -burn 65535 in.geojson rgb.tif

		redTif := r.path + "/" + inS.SelectedImage + "-red.tif"
		greenTif := r.path + "/" + inS.SelectedImage + "-green.tif"
		blueTif := r.path + "/" + inS.SelectedImage + "-blue.tif"
		rgbTif := r.path + "/" + inS.SelectedImage + "-rgb.tif"
		coastTif := r.path + "/" + inS.SelectedImage + "-coast.tif"
		geojson := r.path + "/" + inS.SelectedImage + ".geojson"

		commands := [][]string{
			{"rm", "-f", rgbTif, coastTif},
			{"gdal_merge.py", "-q", "-separate", "-o", rgbTif, redTif, greenTif, blueTif},
			{"cp", "-f", rgbTif, coastTif},
			{"gdal_rasterize", "-q", "-b", "1", "-b", "2", "-b", "3", "-burn", "65535", "-burn", "65535", "-burn", "65535", geojson, coastTif},
		}
		for _, args := range commands {
			mlog.Debug(args)
			stdout, stderr, err := r.run(args[0], args[1:])
			if err != nil {
				panic(err)
			}
			if stdout != "" {
				mlog.Debugf("STDOUT: %s", stdout)
			}
			if stderr != "" {
				mlog.Debugf("STDERR: %s", stderr)
				panic(stderr)
			}
		}
	}

	outS.SelectedImage = inS.SelectedImage

	outJ, err := outS.WriteToJSON()
	if err != nil {
		panic(err)
	}

	r.Output <- outJ
}

func (r *Rasterize) run(nam string, args []string) (string, string, error) {
	var err error

	cmd := exec.Command(nam, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", "", err
	}

	slurpErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return "", "", err
	}

	slurpOut, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", "", err
	}

	err = cmd.Wait()
	if err != nil {
		return "", "", err
	}

	return string(slurpOut), string(slurpErr), nil
}
