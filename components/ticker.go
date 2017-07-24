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
	"log"
	"sync"

	"math/rand"
	"time"

	"github.com/venicegeo/belltower/engine"
)

func init() {
	engine.Factory.Register("Ticker", &Ticker{})
}

type TickerConfigData struct {
	// Time between each tick event, expressed as a Duration.
	// Default is "1s".
	Interval time.Duration

	// If false, the interval lengths will be constant, using the Interval setting.
	// If true, the interval lengths will be random, in the range [0..Interval).
	// Default is false.
	IsRandomized bool

	// How many ticks should be issued before stopping.
	// If zero, will never stop.
	Limit int

	// Initial value of ticker. Default is zero, of course.
	// (This weird feature is useful in some testing scenarios.)
	InitialValue int
}

// Nope.
type TickerInputData struct{}

// implements Serializer
type TickerOutputData struct {
	// Number of ticks sent, including this one. The count starts at 1.
	Count int
}

func (m *TickerOutputData) Validate() error               { return nil } // TODO
func (m *TickerOutputData) ReadFromJSON(jsn string) error { return engine.ReadFromJSON(jsn, m) }
func (m *TickerOutputData) WriteToJSON() (string, error)  { return engine.WriteToJSON(m) }

type Ticker struct {
	engine.ComponentCore

	Input  <-chan string
	Output chan<- string

	// lock around state change in Run()
	StateLock *sync.Mutex

	rnd *rand.Rand

	// local state
	isRandomized bool
	counter      int
	interval     time.Duration
	limit        int
	initialValue int
}

func (ticker *Ticker) Configure() error {

	data := TickerConfigData{}
	err := ticker.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	ticker.rnd = rand.New(rand.NewSource(17))
	ticker.rnd.Seed(17)

	ticker.interval = data.Interval
	ticker.limit = data.Limit
	ticker.isRandomized = data.IsRandomized
	ticker.initialValue = data.InitialValue

	ticker.counter = data.InitialValue

	log.Printf("%#v", ticker.counter)
	return nil
}

func (ticker *Ticker) OnInput(_ string) {
	fmt.Printf("Ticker OnInput\n")

	f := func() {
		ticker.counter++
		fmt.Printf("Ticker.Run: counter=%d\n", ticker.counter)

		outS := &TickerOutputData{
			Count: ticker.counter,
		}

		outJ, err := outS.WriteToJSON()
		if err != nil {
			panic(err)
		}

		ticker.Output <- outJ
	}

	// TODO: should be using time.Ticker
	for {
		secs := ticker.interval.Seconds()
		if ticker.isRandomized {
			// TODO: not tested
			secs = ticker.rnd.Float64() * secs
		}
		time.Sleep(time.Duration(secs) * time.Second)

		f()

		//	time.Sleep(500 * time.Millisecond)

		if ticker.limit > 0 && ticker.counter >= ticker.limit {
			break
		}
	}
}
