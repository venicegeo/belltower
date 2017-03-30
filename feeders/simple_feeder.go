package feeders

import (
	"github.com/venicegeo/belltower/common"
)

// SimpleFeeder just returns a hit count at every poll time.

const SimpleFeederId common.Ident = "7e6c6369-4743-4de9-8407-e6f74198fcb9"

func init() {
	feederFactory.register(&SimpleFeeder{})
}

type SimpleEventData struct {
	Hits  int
	Value int // the square of the value passed in via Settings
}

type SimpleSettings struct {
	Value int
}

type SimpleFeeder struct {
	feed     Feed
	settings SimpleSettings
	hits     int
}

func (_ *SimpleFeeder) Create(feed Feed) (Feeder, error) {
	f := &SimpleFeeder{
		feed:     feed,
		settings: feed.SettingsValues.(SimpleSettings),
		hits:     0,
	}
	return f, nil
}

func (f *SimpleFeeder) GetName() string { return "SimpleFeeder" }

func (f *SimpleFeeder) Id() common.Ident { return SimpleFeederId }

func (f *SimpleFeeder) SettingsSchema() map[string]string {
	return map[string]string{
		"Value": "integer",
	}
}

func (f *SimpleFeeder) EventSchema() map[string]string {
	return map[string]string{
		"Hits":  "integer",
		"Value": "integer",
	}
}

func (f *SimpleFeeder) Poll() (interface{}, error) {
	f.hits++

	e := &SimpleEventData{
		Hits:  f.hits,
		Value: f.settings.Value * f.settings.Value,
	}

	return e, nil
}
