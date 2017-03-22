package feeders

import (
	"encoding/json"
	"log"
	"time"
)

const (
	FeedTypeRandom FeedType = iota
	FeedTypeClock
	FeedTypeHttp
	FeedTypeFileSys
)

type Controller struct {
}

func (c *Controller) Execute() {

	randSettings := &RandomSettings{}
	randSettings.Interval = 3
	randSettings.Timeout = 12

	randRunner := NewRandomRunner(randSettings)

	feeds := []FeedRunner{randRunner}

	for {
		event := feeds[0].Poll()
		if event != nil {
			c.Post(event)
		}
		d := time.Duration(feeds[0].GetInterval()) * time.Second
		time.Sleep(time.Duration(d))
	}
}

func (c *Controller) Post(e *FeedEvent) {
	log.Printf("POSTED %T %#v", e, e.Data.(*RandomEventData))
	byts, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	log.Printf("%s", string(byts))

	fe := FeedEvent{}
	err = json.Unmarshal(byts, &fe)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v", fe)
}

//---------------------------------------------------------------------

type FeedRunnerCore struct {
	ID       uint
	Name     string
	FeedType FeedType
}

type FeedSettingsCore struct {
	Interval int
	Timeout  int
}

type FeedEvent struct {
	Timestamp time.Time
	Data      interface{}
}

type FeedRunner interface {
	Poll() *FeedEvent
	GetInterval() int
	GetTimeout() int
}

//---------------------------------------------------------------------

type RandomRunner struct {
	FeedRunnerCore
	Settings *RandomSettings
}

type RandomEventData struct {
	Seconds int
}

type RandomSettings struct {
	FeedSettingsCore
}

func NewRandomRunner(settings *RandomSettings) *RandomRunner {
	settings = &RandomSettings{}
	settings.Interval = 3
	settings.Timeout = 12

	r := &RandomRunner{}
	r.Settings = settings

	return r
}

func (r *RandomRunner) GetInterval() int {
	return r.Settings.Interval
}

func (r *RandomRunner) GetTimeout() int {
	return r.Settings.Timeout
}

func (r *RandomRunner) Poll() *FeedEvent {
	now := time.Now()
	secs := now.Second()
	if secs%2 == 0 {
		e := &FeedEvent{
			Timestamp: now,
			Data:      &RandomEventData{Seconds: secs},
		}
		return e
	}
	return nil
}
