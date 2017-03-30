package feeders

import (
	"log"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestFeed(t *testing.T) {
	assert := assert.New(t)

	// check registry in general
	assert.True(len(feedersRegistry) > 0)
	for _, feed := range feedersRegistry {
		assert.Implements((*Feeder)(nil), feed)
	}

	// check the SimpleFeeder (a feeder type)
	feeder := feedersRegistry["101"]
	assert.Equal("integer", feeder.SettingsSchema()["Value"])
	assert.Equal("integer", feeder.EventSchema()["Hits"])
	assert.Equal("integer", feeder.EventSchema()["Value"])

	// make a Feed (a feeder instance), call the runloop w/ the post-event function
	value := 7
	limit := 3
	feed := &Feed{
		Id:       "99",
		Name:     "myfeed",
		FeederId: "101",
		Interval: 500 * time.Millisecond,
		SettingsValues: SimpleSettings{
			Value: value,
			Limit: limit,
		},
	}

	// define the post-event function, which takes a Feeder's event object
	hits := 0
	poster := func(event *Event) error {
		log.Printf("EVENT for FeedID: %s", event.FeedId)
		log.Printf("EVENT Data: %#v", event.Data.(*SimpleEventData))
		hits++
		assert.Equal(hits, event.Data.(*SimpleEventData).Hits)
		assert.Equal(value*value, event.Data.(*SimpleEventData).Value)
		return nil
	}

	err := RunFeed(feed, poster)
	assert.Error(err)

	// verify the event was posted
	assert.Equal(3, hits)
}
