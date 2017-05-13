package feeders

import (
	"context"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/btorm"
)

func TestSimpleFeeder(t *testing.T) {
	assert := assert.New(t)

	const secs = 3

	// make a Feed (a feeder instance), call the runloop w/ the post-event function
	feed := &btorm.Feed{}
	feed.Id = "10"
	feed.Name = "mysimplefeed"
	feed.FeederId = SimpleFeederId
	feed.PollingInterval = 1
	feed.PollingEndAt = time.Now().Add(time.Duration(secs) * time.Second)
	feed.Settings = map[string]string{
		"Value": "7",
	}

	// check the SimpleFeeder (a feeder type)
	feeder, err := feederRegistry.data[feed.FeederId].Create(feed)
	assert.NoError(err)

	assert.Equal("integer", feeder.GetSettingsSchema()["Value"])
	assert.Equal("integer", feeder.GetEventSchema()["Hits"])
	assert.Equal("integer", feeder.GetEventSchema()["Value"])

	// define the post-event function, which takes a Feeder's event object
	hits := 0
	poster := func(event *Event) error {
		hits++
		assert.Equal(hits, event.Data.(*SimpleEventData).Hits)
		assert.Equal(49, event.Data.(*SimpleEventData).Value)
		return nil
	}

	ctx, cancel := context.WithDeadline(context.Background(), feed.PollingEndAt)
	assert.NotNil(cancel)

	err = runFeed(ctx, feed, feeder, poster)
	assert.Error(err)

	// verify the event was posted; all times approximate
	//log.Printf("%d -- ", hits)
	assert.True(hits >= secs-1 && hits <= secs+1)
}
