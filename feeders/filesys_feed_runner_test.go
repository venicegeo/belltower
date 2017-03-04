package feeders

import (
	"os"
	"testing"
	"time"

	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/orm"
)

func touchFile(t *testing.T, path string) string {
	assert := assert.New(t)

	r := time.Now().Nanosecond()
	nam := fmt.Sprintf("%s/%d.tmp", path, r)
	f, err := os.OpenFile(nam, os.O_RDWR|os.O_CREATE, 0755)
	assert.NoError(err)
	err = f.Close()
	assert.NoError(err)
	return nam
}

func TestFileSysFeed(t *testing.T) {
	assert := assert.New(t)

	feed := &orm.Feed{
		ID:   17,
		Name: "filesystest",
		Settings: map[string]interface{}{
			"path":  "/tmp",
			"name":  "filesys1",
			"sleep": "1",
		},
	}
	runner, err := NewFileSysFeedRunner(feed)
	assert.NoError(err)

	var addedFile string

	statusCount := 5
	hitCount := 1

	statusF := func(s string) (bool, error) {
		assert.Equal("good", s)
		statusCount--
		if statusCount == 0 {
			return false, nil
		}
		if statusCount == 3 {
			path, err := common.GetMapValueAsString(feed.Settings, "path")
			assert.NoError(err)
			addedFile = touchFile(t, path)
		}
		return true, nil
	}

	mssgF := func(data map[string]string) error {
		//log.Printf("event==> %s ... %s", data["mssg"], data["added"])
		path, err := common.GetMapValueAsString(feed.Settings, "path")
		assert.NoError(err)
		assert.Equal(addedFile, path+"/"+data["added"])
		hitCount--
		return nil
	}

	err = runner.Run(statusF, mssgF)
	assert.NoError(err)

	time.Sleep(7 * time.Second)

	assert.Zero(statusCount)
	assert.Zero(hitCount)
}
