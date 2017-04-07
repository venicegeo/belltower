package feeders

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/btorm"
)

func TestSimpleFeeder(t *testing.T) {
	assert := assert.New(t)

	// make a Feed (a feeder instance), call the runloop w/ the post-event function
	value := 7
	feed := &btorm.Feed{}
	feed.Id = "10"
	feed.Name = "mysimplefeed"
	feed.FeederId = SimpleFeederId
	feed.PollingInterval = 1
	feed.PollingEndAt = time.Now().Add(3 * time.Second)
	feed.Settings = map[string]interface{}{
		"Value": value,
	}

	// check the SimpleFeeder (a feeder type)
	feeder, err := feederFactory.create(feed)
	assert.NoError(err)

	assert.Equal("integer", feeder.SettingsSchema()["Value"])
	assert.Equal("integer", feeder.EventSchema()["Hits"])
	assert.Equal("integer", feeder.EventSchema()["Value"])

	// define the post-event function, which takes a Feeder's event object
	hits := 0
	poster := func(event *Event) error {
		hits++
		assert.Equal(hits, event.Data.(*SimpleEventData).Hits)
		assert.Equal(value*value, event.Data.(*SimpleEventData).Value)
		return nil
	}

	err = RunFeed(feed, feeder, poster)
	assert.Error(err)

	// verify the event was posted
	assert.Equal(3, hits)
}
