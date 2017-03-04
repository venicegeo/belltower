package feeders

import (
	"os"
	"testing"
	"time"

	"fmt"

	"github.com/stretchr/testify/assert"
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
	assert.True(true)

	tf := FileSysFeed{}

	var addedFile string

	statusCount := 5
	hitCount := 1

	config := Config{
		"path":  "/tmp",
		"name":  "filesys1",
		"sleep": "1",
	}

	statusF := func(s string) (bool, error) {
		assert.Equal("good", s)
		statusCount--
		if statusCount == 0 {
			return false, nil
		}
		if statusCount == 3 {
			addedFile = touchFile(t, config["path"])
		}
		return true, nil
	}

	mssgF := func(data map[string]string) error {
		//log.Printf("event==> %s ... %s", data["mssg"], data["added"])
		assert.Equal(addedFile, config["path"]+"/"+data["added"])
		hitCount--
		return nil
	}

	err := tf.Run(config, statusF, mssgF)
	assert.NoError(err)

	time.Sleep(7 * time.Second)

	assert.Zero(statusCount)
	assert.Zero(hitCount)
}
