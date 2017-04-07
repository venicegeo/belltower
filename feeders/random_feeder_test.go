package feeders

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/btorm"
)

func TestRandomFeeder(t *testing.T) {
	assert := assert.New(t)

	feed := &btorm.Feed{}
	feed.Id = "20"
	feed.Name = "myrandomfeed"
	feed.FeederId = RandomFeederId
	feed.PollingInterval = 1
	feed.PollingEndAt = time.Now().Add(5 * time.Second)
	feed.Settings = map[string]interface{}{
		"Target": 50,
		"Seed":   19, // 61, 52, 88, 23, 30, 70
	}

	feeder, err := feederFactory.create(feed)
	assert.NoError(err)

	assert.Equal("integer", feeder.SettingsSchema()["Target"])
	assert.Equal("integer", feeder.SettingsSchema()["Seed"])
	assert.Equal("integer", feeder.EventSchema()["Value"])

	// define the post-event function, which takes a Feeder's event object
	hits := 0
	poster := func(event *Event) error {
		hits++
		return nil
	}

	err = RunFeed(feed, feeder, poster)
	assert.Error(err)

	// verify the event was posted
	assert.Equal(2, hits)
}
