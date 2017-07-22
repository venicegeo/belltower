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
	"testing"
	"time"

	"os"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/engine"
)

func creat(path string, name string) error {
	f, err := os.Create(path + "/" + name)
	if err != nil {
		return err
	}
	f.WriteString("hi")
	f.Close()
	return nil
}

func TestWalker(t *testing.T) {
	assert := assert.New(t)

	var err error

	const path = "/tmp/testwalker"

	_ = os.RemoveAll(path)
	err = os.MkdirAll(path, 0700)
	assert.NoError(err)

	err = creat(path, "xxx")
	assert.NoError(err)

	stim := time.Now()

	// TODO: why does this have to be so long?
	time.Sleep(time.Millisecond * 1000)

	err = creat(path, "aaa")
	assert.NoError(err)

	etim := time.Now()

	files, err := walk(path, stim, etim)
	assert.NoError(err)

	assert.Len(files, 1)
	assert.Equal(path+"/aaa", files[0])
	//assert.True(files[0].ModTime().After(stim))
	//assert.True(files[0].ModTime().Before(etim) || files[0].ModTime().Equal(etim))

	_ = os.RemoveAll(path)
}

func TestFileWatcher(t *testing.T) {
	assert := assert.New(t)

	var err error

	const path = "/tmp/testfilewatcher"
	_ = os.RemoveAll(path)
	err = os.MkdirAll(path, 0700)
	assert.NoError(err)

	err = creat(path, "xxx")
	assert.NoError(err)

	time.Sleep(time.Millisecond * 1000)

	config := engine.ArgMap{
		"path": path,
	}
	fwX, err := engine.Factory.Create("FileWatcher", config)
	assert.NoError(err)
	fw := fwX.(*FileWatcher)

	// this setup is normally done by goflow itself
	chIn := make(chan string)
	chOut := make(chan string)
	fw.Input = chIn
	fw.Output = chOut

	err = creat(path, "xxx")
	assert.NoError(err)

	// first run
	{
		inJ := "{}"
		go fw.OnInput(inJ)

		// check the returned result
		outJ := <-chOut
		outS := &FileWatcherOutputData{}
		err = outS.ReadFromJSON(outJ)
		assert.NoError(err)

		assert.True(outS.StartTime.IsZero())
		assert.False(outS.EndTime.IsZero())
		assert.Len(outS.Names, 0)
		assert.Equal(path, outS.Path)
	}

	time.Sleep(time.Millisecond * 1000)
	err = creat(path, "yyy")
	assert.NoError(err)
	time.Sleep(time.Millisecond * 1000)

	// second run
	{
		inJ := "{}"
		go fw.OnInput(inJ)

		// check the returned result
		outJ := <-chOut
		outS := &FileWatcherOutputData{}
		err = outS.ReadFromJSON(outJ)
		assert.NoError(err)

		assert.False(outS.StartTime.IsZero())
		assert.False(outS.EndTime.IsZero())
		assert.Len(outS.Names, 1)
		assert.Equal(path+"/yyy", outS.Names[0])
		assert.Equal(path, outS.Path)
	}

	time.Sleep(time.Millisecond * 1000)
	err = os.MkdirAll(path+"/dir", 0700)
	assert.NoError(err)
	err = creat(path, "zzz0")
	err = creat(path+"/dir", "zzz1")
	assert.NoError(err)
	time.Sleep(time.Millisecond * 1000)

	// third run
	{
		inJ := "{}"
		go fw.OnInput(inJ)

		// check the returned result
		outJ := <-chOut
		outS := &FileWatcherOutputData{}
		err = outS.ReadFromJSON(outJ)
		assert.NoError(err)

		assert.False(outS.StartTime.IsZero())
		assert.False(outS.EndTime.IsZero())
		assert.Len(outS.Names, 2)
		assert.Equal(path+"/dir/zzz1", outS.Names[0])
		assert.Equal(path+"/zzz0", outS.Names[1])
		assert.Equal(path, outS.Path)
	}

	_ = os.RemoveAll(path)
}
