package components

import (
	"fmt"
	"sync"

	"github.com/venicegeo/belltower/common"

	"time"
)

func init() {
	common.Factory.Register("Ticker", &Ticker{})
}

type TickerConfigData struct {
	//   Time between each tick event, expressed as a Duration.
	//   Default is "1s".
	Interval time.Duration

	//   If false, the interval lengths will be constant, using the Interval setting.
	//   If true, the interval lengths will be random, in the range [0..Interval).
	//   Default is false.
	IsRandomized bool

	//   How many ticks should be issued before stopping.
	//   If zero, will never stop.
	Limit int
}

// Nope.
type TickerInputData struct{}

// implements Serializer
type TickerOutputData struct {
	// Number of ticks sent, including this one. The count starts at 1.
	Count int
}

func (m *TickerOutputData) Validate() error               { return nil } // TODO
func (m *TickerOutputData) ReadFromJSON(jsn string) error { return common.ReadFromJSON(jsn, m) }
func (m *TickerOutputData) WriteToJSON() (string, error)  { return common.WriteToJSON(m) }

type Ticker struct {
	common.ComponentCore

	Input  <-chan string
	Output chan<- string

	// lock around state change in Run()
	StateLock *sync.Mutex

	// local state
	isRandomized bool
	counter      int
	interval     time.Duration
	limit        int
}

func (ticker *Ticker) Configure() error {

	data := TickerConfigData{}
	err := ticker.Config.ToStruct(&data)
	if err != nil {
		return err
	}

	ticker.interval = data.Interval
	ticker.limit = data.Limit
	ticker.isRandomized = data.IsRandomized

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

	for {
		time.Sleep(ticker.interval)

		f()

		//	time.Sleep(500 * time.Millisecond)

		if ticker.limit > 0 && ticker.counter >= ticker.limit {
			break
		}
	}
}
