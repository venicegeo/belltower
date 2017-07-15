package components

import (
	"fmt"
	"sync"

	"time"

	"github.com/venicegeo/belltower/common"
)

func init() {
	Factory.Register("Ticker", &Ticker{})
}

// -- CONFIG --
//
// interval string
//   Time between each tick event, expressed as a Duration.
//   Default is "1s".
//
// isRandomized bool
//   If false, the interval lengths will be constant, using the Interval setting.
//   If true, the interval lengths will be random, in the range [0..Interval).
//   Default is false.
//
// limit int
//   How many ticks should be issued before stopping.
//   If zero, will never stop.
//
// -- INPUT --
//
// (none)
//
// -- OUTPUT --
//
// count int
//   Number of ticks sent, including this one. The count starts at 1.

type Ticker struct {
	ComponentCore

	Input  <-chan string
	Output chan<- string

	// lock around state change in Run()
	StateLock *sync.Mutex

	// local state
	isRandomized bool
	counter      float64
	interval     time.Duration
	limit        float64
}

func (ticker *Ticker) localConfigure() error {

	interval, err := ticker.config.GetStringOrDefault("interval", "1s")
	if err != nil {
		return err
	}

	limit, err := ticker.config.GetFloatOrDefault("limit", 0.0)
	if err != nil {
		return err
	}

	isRandomized, err := ticker.config.GetBoolOrDefault("isRandomized", false)
	if err != nil {
		return err
	}

	ticker.interval, err = time.ParseDuration(interval)
	if err != nil {
		return err
	}
	ticker.limit = limit
	ticker.isRandomized = isRandomized

	return nil
}

func (ticker *Ticker) OnInput(data string) {
	fmt.Printf("Ticker OnInput: %s\n", data)

	_, err := common.NewArgMap(data)
	if err != nil {
		panic(err)
	}

	f := func() {
		out, err := ticker.Run(nil)
		if err != nil {
			panic(err)
		}

		s, err := out.ToJSON()
		if err != nil {
			panic(err)
		}

		ticker.Output <- s
	}

	for {
		time.Sleep(ticker.interval)

		f()

		//	time.Sleep(500 * time.Millisecond)

		if ticker.limit > 0.0 && ticker.counter >= ticker.limit {
			break
		}
	}
}

func (ticker *Ticker) Run(in common.ArgMap) (common.ArgMap, error) {

	ticker.counter++

	fmt.Printf("Ticker.Run: counter=%f\n", ticker.counter)

	out := common.ArgMap{}
	out["count"] = ticker.counter

	return out, nil
}
