package feeders

import (
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/btorm"
)

func touchFile(t *testing.T, nam string) {
	assert := assert.New(t)

	f, err := os.OpenFile(nam, os.O_RDWR|os.O_CREATE, 0755)
	assert.NoError(err)
	err = f.Close()
	assert.NoError(err)
}

func deleteFile(t *testing.T, nam string) {
	assert := assert.New(t)
	err := os.Remove(nam)
	assert.NoError(err)
}

func TestFileSysFeeder(t *testing.T) {
	assert := assert.New(t)

	const root = "/tmp"

	feed := btorm.Feed{}
	feed.Id = "30"
	feed.Name = "myfilesysfeed"
	feed.FeederId = FileSysFeederId
	feed.PollingInterval = 1
	feed.PollingEndAt = time.Now().Add(3 * time.Second)
	feed.Settings = FileSysSettings{
		Path: root,
	}

	feeder, err := feederFactory.create(feed)
	assert.NoError(err)

	assert.Equal("string", feeder.SettingsSchema()["Path"])
	assert.Equal("string", feeder.EventSchema()["Added"])
	assert.Equal("string", feeder.EventSchema()["Deleted"])

	// make up some filenames to add and delete
	basename := fmt.Sprintf("filesys.%d.", time.Now().Nanosecond())
	names := make([]string, 100)
	for i := 0; i < 100; i++ {
		names[i] = basename + strconv.Itoa(i)
	}

	// define the post-event function, which takes a Feeder's event object
	count := 0
	added := ""
	deleted := ""
	poster := func(event *Event) error {
		added = event.Data.(*FileSysEventData).Added
		deleted = event.Data.(*FileSysEventData).Deleted

		log.Printf("Added: %s", added)
		log.Printf("Deled: %s", deleted)
		log.Printf("Count: %d", count)

		count++

		// giveth
		touchFile(t, root+"/"+names[count])

		if count > 1 {
			// taketh away
			deleteFile(t, root+"/"+names[count-1])
		}

		return nil
	}

	err = RunFeed(feed, poster)
	assert.Error(err)

	// verify the event was posted
	assert.Equal(names[count-1], added)
	assert.Equal(names[count-2], deleted)
}
