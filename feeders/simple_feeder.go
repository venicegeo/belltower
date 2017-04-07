package feeders

import (
	"github.com/venicegeo/belltower/btorm"
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

type SimpleFeeder struct {
	feed     *btorm.Feed
	settings map[string]interface{}
	hits     int
}

func (f *SimpleFeeder) Create(feed *btorm.Feed) (Feeder, error) {

	err := checkSchema(f.SettingsSchema(), feed.Settings)
	if err != nil {
		return nil, err
	}

	feeder := &SimpleFeeder{
		feed:     feed,
		settings: feed.Settings,
		hits:     0,
	}
	return feeder, nil
}

func (f *SimpleFeeder) GetName() string { return "SimpleFeeder" }

func (f *SimpleFeeder) Id() common.Ident { return SimpleFeederId }

func (_ *SimpleFeeder) SettingsSchema() map[string]string {
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

	x := f.settings["Value"].(float64)

	e := &SimpleEventData{
		Hits:  f.hits,
		Value: int(x * x),
	}

	return e, nil
}
