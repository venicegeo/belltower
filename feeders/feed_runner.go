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
	FeedId    common.Ident
	Data      interface{}
}

type Feeder interface {
	GetName() string
	GetId() common.Ident
	GetFeedType() FeedType
	Poll() (*FeedEvent, error)
	GetInterval() time.Duration
}

type FeederCore struct {
	id       common.Ident
	name     string
	feedType FeedType
	interval time.Duration
}

func (f *FeederCore) GetId() common.Ident        { return f.id }
func (f *FeederCore) GetName() string            { return f.name }
func (f *FeederCore) GetFeedType() FeedType      { return f.feedType }
func (f *FeederCore) GetInterval() time.Duration { return f.interval }

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

	randFeeder := NewRandomFeeder(settings)

	feeds := []Feeder{randFeeder}

	for _, feed := range feeds {
		go c.RunLoop(feed)
	}
}

func (c *Controller) RunLoop(feeder Feeder) {
	errCount := 0

	for {
		if errCount == 3 {
			log.Printf("Feeder aborting (%s)", feeder.GetName())
			break
		}

		event, err := feeder.Poll()
		if err != nil {
			log.Printf("Feeder poll error (%s): %v", feeder.GetName(), err)
			errCount++
		} else {
			if event != nil {
				err = c.Post(event)
				log.Printf("Feeder post error (%s): %v", feeder.GetName(), err)
				// TODO
			}
		}

		time.Sleep(feeder.GetInterval())
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

type RandomFeeder struct {
	FeederCore
	settings *RandomSettings
}

type RandomEventData struct {
	Value int
}

type RandomSettings struct {
	Name         string
	PollInterval time.Duration
	Target       int // value in range [0..100], 20 means will hit 20% of the time
}

func NewRandomFeeder(settings *RandomSettings) *RandomFeeder {
	r := &RandomFeeder{}
	r.feedType = FeedTypeRandom
	r.id = common.NewId()
	r.interval = time.Second * 3
	r.name = settings.Name
	r.settings = settings

	return r
}

func (r *RandomFeeder) Poll() (*FeedEvent, error) {
	x := rand.Intn(100)
	if x > r.settings.Target {
		// not a hit
		return nil, nil
	}

	e := &FeedEvent{
		TimeStamp: time.Now(),
		Data:      &RandomEventData{Value: x},
	}
	return e, nil
}
