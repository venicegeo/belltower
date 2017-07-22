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
	"os"
	"path/filepath"
	"time"

	"github.com/venicegeo/belltower/engine"
	"github.com/venicegeo/belltower/mpg/mlog"
)

func init() {
	engine.Factory.Register("FileWatcher", &FileWatcher{})
}

type FileWatcherConfigData struct {
	Path string
}

type FileWatcherInputData struct {
}

func (m *FileWatcherInputData) Validate() error               { return nil } // TODO
func (m *FileWatcherInputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *FileWatcherInputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type FileWatcherOutputData struct {
	Path      string
	Names     []string
	StartTime time.Time
	EndTime   time.Time
}

func (m *FileWatcherOutputData) Validate() error               { return nil } // TODO
func (m *FileWatcherOutputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *FileWatcherOutputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type FileWatcher struct {
	engine.ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	startTime time.Time
	endTime   time.Time
	path      string
}

func (fw *FileWatcher) Configure() error {

	data := FileWatcherConfigData{}
	err := fw.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	fw.path = data.Path

	return nil
}

// When triggered, we run the watcher once; this implies we need to be
// connected to a timer. The interval is taken from the last time we
// were run (and the first time we're run we just record the start time).
func (fw *FileWatcher) OnInput(inJ string) {
	fmt.Printf("FileWatcher OnInput: %s\n", inJ)

	var err error

	now := time.Now()

	outS := &FileWatcherOutputData{}

	// TODO: support checking for new files, deleted files, updated files
	// TODO: return the whole os.FileInfo object, not just the name

	files := []string{}
	if fw.endTime.IsZero() {
		// the first time, so we just establish the starting time
		fw.endTime = now
		mlog.Printf("Not walking: to %s", fw.endTime)

	} else {
		fw.startTime = fw.endTime
		fw.endTime = now

		mlog.Printf("Walking: %s to %s", fw.startTime, fw.endTime)

		files, err = walk(fw.path, fw.startTime, fw.endTime)
		if err != nil {
			panic(err)
		}
	}

	// TODO: lastRun is set too late -- need to have walk() use a time range,
	// and make the new time range start at the end of the old one

	outS.Names = files
	outS.Path = fw.path
	outS.StartTime = fw.startTime
	outS.EndTime = fw.endTime

	outJ, err := outS.WriteToJSON()
	if err != nil {
		panic(err)
	}

	fw.Output <- outJ
}

func inTimeRange(t, start, end time.Time) bool {
	// start <= t < end
	return (t.After(start) || t.Equal(start)) && t.Before(end)
}

func walk(path string, startTime time.Time, endTime time.Time) ([]string, error) {

	files := []string{}

	f := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// TODO: too conservative here
			return nil
		}

		inRange := inTimeRange(info.ModTime(), startTime, endTime)

		if inRange && !info.IsDir() {
			mlog.Printf("HIT: %s[%s] at %s\n", path, info.Name(), info.ModTime())
			files = append(files, path)
		}
		return nil
	}

	err := filepath.Walk(path, f)
	if err != nil {
		return nil, err
	}
	return files, nil
}
