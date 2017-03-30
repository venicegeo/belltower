package feeders

import (
	"math/rand"
	"time"

	"github.com/venicegeo/belltower/common"
)

// RandFeed checks to see if a random number from the range [0..limit) is equal to zero. If so,
// it will give the server a message. It sleeps for the given number of seconds between
// checks.

func init() {
	registerFeed("100", &RandomFeeder{})
}

type RandomEventData struct {
	Value int
}

type RandomSettings struct {
	Target int // value in range [0..100], 20 means will hit 20% of the time
}

type RandomFeeder struct {
}

func (r *RandomFeeder) GetName() string { return "RandomFeeder" }

func (r *RandomFeeder) Id() common.Ident { return "100" }

func (r *RandomFeeder) SettingsSchema() map[string]string {
	return map[string]string{
		"Target": "integer",
	}
}

func (r *RandomFeeder) EventSchema() map[string]string {
	return map[string]string{
		"Value": "integer",
	}
}

func (r *RandomFeeder) Poll(feedId common.Ident, isettings interface{}) (*Event, error) {
	settings := isettings.(RandomSettings)

	x := rand.Intn(100)
	if x > settings.Target {
		// not a hit
		return nil, nil
	}

	e := &Event{
		TimeStamp: time.Now(),
		FeedId:    feedId,
		Data:      &RandomEventData{Value: x},
	}
	return e, nil
}
