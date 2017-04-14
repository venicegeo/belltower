package feeders

import (
	"math/rand"
	"time"

	"strconv"

	"github.com/venicegeo/belltower/btorm"
	"github.com/venicegeo/belltower/common"
)

// RandFeed checks to see if a random number from the range [0..limit) is equal to zero. If so,
// it will give the server a message. It sleeps for the given number of seconds between
// checks.

const RandomFeederId common.Ident = "73c9b03e-5455-4aa6-8c64-b1abeddf3763"

// implements Feeder
type RandomFeeder struct {
	random   *rand.Rand
	settings map[string]string
}

type RandomEventData struct {
	Value int
}

func (f *RandomFeeder) GetName() string { return "RandomFeeder" }

func (f *RandomFeeder) GetId() common.Ident { return RandomFeederId }

func (f *RandomFeeder) GetSettingsSchema() map[string]string {
	return map[string]string{
		"Target": "integer",
		"Seed":   "integer",
	}
}

func (f *RandomFeeder) GetEventSchema() map[string]string {
	return map[string]string{
		"Value": "integer",
	}
}

func init() {
	info := &FeederInfo{
		FeederId:    RandomFeederId,
		Description: "random feeder",
		Create:      RandomFeederCreate,
	}
	feederRegistry.register(info)
}

func RandomFeederCreate(feed *btorm.Feed) (Feeder, error) {

	seed32, err := strconv.Atoi(feed.Settings["Seed"])
	if err != nil {
		return nil, err
	}
	seed64 := int64(seed32)
	if seed64 == 0 {
		seed64 = time.Now().UnixNano()
	}

	f := &RandomFeeder{
		settings: feed.Settings,
		random:   rand.New(rand.NewSource(seed64)),
	}

	return f, nil
}

func (f *RandomFeeder) Poll() (interface{}, error) {

	x := f.random.Intn(100)
	target, err := strconv.Atoi(f.settings["Target"])
	if err != nil {
		return nil, err
	}
	if x > target {
		// not a hit
		return nil, nil
	}

	e := &RandomEventData{
		Value: x,
	}

	return e, nil
}
