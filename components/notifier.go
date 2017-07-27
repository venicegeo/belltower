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
	"log"

	"github.com/deckarep/gosx-notifier"
	"github.com/venicegeo/belltower/engine"
	"github.com/venicegeo/belltower/mpg/mlog"
)

func init() {
	engine.Factory.Register("Notifier", &Notifier{})
}

type NotifierConfigData struct {

	// The value added to the input. Default is zero.
	Path string
}

type NotifierInputData struct {

	// The value added to the addend from the configuration. Default is zero.
	SelectedImage string
}

func (m *NotifierInputData) Validate() error               { return nil } // TODO
func (m *NotifierInputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *NotifierInputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type NotifierOutputData struct {

	// Value of input value added to addend.
	SelectedImage string
}

func (m *NotifierOutputData) Validate() error               { return nil } // TODO
func (m *NotifierOutputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *NotifierOutputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type Notifier struct {
	engine.ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	path string
}

func (notifier *Notifier) Configure() error {

	data := NotifierConfigData{}
	err := notifier.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	notifier.path = data.Path

	return nil
}

func (notifier *Notifier) OnInput(inJ string) {
	mlog.Printf("Notifier OnInput: %s\n", inJ)

	inS := &NotifierInputData{}
	err := inS.ReadFromJSON(inJ)
	if err != nil {
		panic(err)
	}

	outS := &NotifierOutputData{}

	// the work
	{
		path := "file://" + notifier.path + "/"
		link := path + inS.SelectedImage + "-coast.tif"
		log.Print(link)
		home := "/Users/mpgerlek/venicegeo/belltower/etc/images/"

		note := gosxnotifier.NewNotification(inS.SelectedImage)
		note.Title = "Beachfront Job Complete"
		note.Sender = "com.apple.Maps"
		note.Link = link
		note.AppIcon = home + "bf-icon.png"
		note.ContentImage = home + "content-icon.png"

		err := note.Push()
		if err != nil {
			panic(err)
		}
	}

	outJ, err := outS.WriteToJSON()
	if err != nil {
		panic(err)
	}

	notifier.Output <- outJ
}
