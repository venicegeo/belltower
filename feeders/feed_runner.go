package feeders

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

type FeedType int

type FeedEvent struct {
	TimeStamp time.Time
	Data      interface{}
}

type FeedRunner interface {
	GetName() string
	GetId() common.Ident
	GetFeedType() FeedType
	Poll() (*FeedEvent, error)
	GetInterval() time.Duration
}

type FeedRunnerCore struct {
	Id       common.Ident
	Name     string
	FeedType FeedType
	Interval time.Duration
}

//---------------------------------------------------------------------

const (
	FeedTypeRandom FeedType = iota
	FeedTypeClock
	FeedTypeHttp
	FeedTypeFileSys
)

type Controller struct {
}

func (c *Controller) Execute() {

	settings := &RandomSettings{
		Name:         "myrandrunner",
		PollInterval: 3 * time.Second,
		Target:       33,
	}

	randRunner := NewRandomRunner(settings)

	feeds := []FeedRunner{randRunner}

	for _, feed := range feeds {
		go c.RunLoop(feed)
	}
}

func (c *Controller) RunLoop(runner FeedRunner) {
	errCount := 0

	for {
		if errCount == 3 {
			log.Printf("Runner aborting (%s)", runner.GetName())
			break
		}

		event, err := runner.Poll()
		if err != nil {
			log.Printf("Runner poll error (%s): %v", runner.GetName(), err)
			errCount++
		} else {
			if event != nil {
				err = c.Post(event)
				log.Printf("Runner post error (%s): %v", runner.GetName(), err)
				// TODO
			}
		}

		time.Sleep(runner.GetInterval())
	}
}

func (c *Controller) Post(e *FeedEvent) error {
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
	return nil
}

//---------------------------------------------------------------------

type RandomRunner struct {
	Core     FeedRunnerCore
	Settings *RandomSettings
}

type RandomEventData struct {
	Value int
}

type RandomSettings struct {
	Name         string
	PollInterval time.Duration
	Target       int // value in range [0..100], 20 means will hit 20% of the time
}

func NewRandomRunner(settings *RandomSettings) *RandomRunner {
	r := &RandomRunner{}
	r.Core.FeedType = FeedTypeRandom
	r.Core.Id = common.NewId()
	r.Core.Interval = time.Second * 3
	r.Core.Name = settings.Name
	r.Settings = settings

	return r
}

func (r *RandomRunner) GetName() string            { return r.Core.Name }
func (r *RandomRunner) GetId() common.Ident        { return r.Core.Id }
func (r *RandomRunner) GetFeedType() FeedType      { return r.Core.FeedType }
func (r *RandomRunner) GetInterval() time.Duration { return r.Core.Interval }

func (r *RandomRunner) Poll() (*FeedEvent, error) {
	x := rand.Intn(100)
	if x > r.Settings.Target {
		// not a hit
		return nil, nil
	}

	e := &FeedEvent{
		TimeStamp: time.Now(),
		Data:      &RandomEventData{Value: x},
	}
	return e, nil
}
