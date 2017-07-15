package components

import (
	"log"

	"time"

	"github.com/venicegeo/belltower/common"
)

func init() {
	Factory.Register("Ticker", &Ticker{})
}

// -- CONFIG --
//
// Interval string
//   Time between each tick event, expressed as a Duration.
//   Default is "1s".
//
// IsRandomized bool
//   If false, the interval lengths will be constant, using the Interval setting.
//   If true, the interval lengths will be random, in the range [0..Interval).
//   Default is false.
//
// Limit int
//   How many ticks should be issued before stopping.
//   If zero, will never stop.
//
// -- INPUT --
//
// (none)
//
// -- OUTPUT --
//
// Count int
//   Number of ticks sent, including this one. The count starts at 1.

type Ticker struct {
	ComponentCore

	Input  <-chan string
	Output chan<- string

	// local state
	isRandomized bool
	counter      int
	interval     time.Duration
	limit        int
}

func (ticker *Ticker) localConfigure() error {

	interval, err := ticker.config.GetStringOrDefault("interval", "1s")
	if err != nil {
		return err
	}

	limit, err := ticker.config.GetIntOrDefault("limit", 0)
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
	log.Printf("Ticker: OnInput <%s>", data)

	in, err := common.NewArgMap(data)
	if err != nil {
		panic(err)
	}

	out, err := ticker.Run(in)
	if err != nil {
		panic(err)
	}

	s, err := out.ToJSON()
	if err != nil {
		panic(err)
	}

	ticker.Output <- s
}

func (ticker *Ticker) Run(in common.ArgMap) (common.ArgMap, error) {

	ticker.counter++

	out := common.ArgMap{}
	out["count"] = ticker.counter

	return out, nil
}
