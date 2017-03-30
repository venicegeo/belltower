package feeders

import (
	"time"

	"fmt"

	"github.com/venicegeo/belltower/common"
)

// RandFeed checks to see if a random number from the range [0..limit) is equal to zero. If so,
// it will give the server a message. It sleeps for the given number of seconds between
// checks.

func init() {
	registerFeed("101", &SimpleFeeder{})
}

type SimpleEventData struct {
	Hits  int
	Value int // the square of the value passed in via Settings
}

type SimpleSettings struct {
	Value int
	Limit int
}

type SimpleFeeder struct {
	hits int
}

func (f *SimpleFeeder) GetName() string { return "SimpleFeeder" }

func (f *SimpleFeeder) Id() common.Ident { return "101" }

func (f *SimpleFeeder) SettingsSchema() map[string]string {
	return map[string]string{
		"Value": "integer",
		"Limit": "integer", // produce an err after this many hits
	}
}

func (f *SimpleFeeder) EventSchema() map[string]string {
	return map[string]string{
		"Hits":  "integer",
		"Value": "integer",
	}
}

func (f *SimpleFeeder) Poll(feedId common.Ident, isettings interface{}) (*Event, error) {
	settings := isettings.(SimpleSettings)

	if f.hits == settings.Limit {
		return nil, fmt.Errorf("limit reached")
	}

	f.hits++

	e := &Event{
		TimeStamp: time.Now(),
		FeedId:    feedId,
		Data: &SimpleEventData{
			Hits:  f.hits,
			Value: settings.Value * settings.Value,
		},
	}

	return e, nil
}
