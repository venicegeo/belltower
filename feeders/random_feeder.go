package feeders

import (
	"math/rand"
	"time"
)

// RandFeed checks to see if a random number from the range [0..limit) is equal to zero. If so,
// it will give the server a message. It sleeps for the given number of seconds between
// checks.

type RandomEventData struct {
	Value int
}

type RandomSettings struct {
	Target int // value in range [0..100], 20 means will hit 20% of the time
}

type RandomFeeder struct {
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
