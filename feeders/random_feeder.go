package feeders

import (
	"math/rand"
	"time"

	"github.com/venicegeo/belltower/btorm"
	"github.com/venicegeo/belltower/common"
)

// RandFeed checks to see if a random number from the range [0..limit) is equal to zero. If so,
// it will give the server a message. It sleeps for the given number of seconds between
// checks.

const RandomFeederId common.Ident = "73c9b03e-5455-4aa6-8c64-b1abeddf3763"

func init() {
	feederFactory.register(&RandomFeeder{})
}

type RandomEventData struct {
	Value int
}

type RandomSettings struct {
	Target int   // value in range [0..100], 20 means will hit 20% of the time
	Seed   int64 // to allow for debugging; if zero, will use system default
}

type RandomFeeder struct {
	feed     btorm.Feed
	settings RandomSettings
	random   *rand.Rand
}

func (_ *RandomFeeder) Create(feed btorm.Feed) (Feeder, error) {
	settings := feed.Settings.(RandomSettings)
	seed := settings.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	f := &RandomFeeder{
		feed:     feed,
		settings: settings,
		random:   rand.New(rand.NewSource(seed)),
	}

	return f, nil
}

func (r *RandomFeeder) GetName() string { return "RandomFeeder" }

func (r *RandomFeeder) Id() common.Ident { return RandomFeederId }

func (r *RandomFeeder) SettingsSchema() map[string]string {
	return map[string]string{
		"Target": "integer",
		"Seed":   "integer",
	}
}

func (r *RandomFeeder) EventSchema() map[string]string {
	return map[string]string{
		"Value": "integer",
	}
}

func (r *RandomFeeder) Poll() (interface{}, error) {
	settings := r.feed.Settings.(RandomSettings)

	x := r.random.Intn(100)
	if x > settings.Target {
		// not a hit
		return nil, nil
	}

	e := &RandomEventData{
		Value: x,
	}

	return e, nil
}
