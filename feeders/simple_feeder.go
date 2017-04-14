package feeders

import (
	"strconv"

	"github.com/venicegeo/belltower/btorm"
	"github.com/venicegeo/belltower/common"
)

// SimpleFeeder just returns a hit count at every poll time.

const SimpleFeederId common.Ident = "7e6c6369-4743-4de9-8407-e6f74198fcb9"

// implements Feeder
type SimpleFeeder struct {
	settings map[string]string
	hits     int
}

type SimpleEventData struct {
	Hits  int
	Value int // the square of the value passed in via Settings
}

func (f *SimpleFeeder) GetId() common.Ident {
	return SimpleFeederId
}

func (f *SimpleFeeder) GetName() string {
	return "SimpleFeeder"
}

func (f *SimpleFeeder) GetSettingsSchema() map[string]string {
	return map[string]string{
		"Value": "integer",
	}
}

func (f *SimpleFeeder) GetEventSchema() map[string]string {
	return map[string]string{
		"Hits":  "integer",
		"Value": "integer",
	}
}

//---------------------------------------------------------------------

func SimpleFeederCreate(feed *btorm.Feed) (Feeder, error) {

	feeder := &SimpleFeeder{
		settings: feed.Settings,
		hits:     0,
	}
	return feeder, nil
}

func (f *SimpleFeeder) Poll() (interface{}, error) {
	f.hits++

	x, err := strconv.ParseFloat(f.settings["Value"], 64)
	if err != nil {
		return nil, err
	}
	e := &SimpleEventData{
		Hits:  f.hits,
		Value: int(x * x),
	}

	return e, nil
}
