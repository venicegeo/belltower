package feeders

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/orm"
)

// RandFeed checks to see if a random number from the range [0..limit) is equal to zero. If so,
// it will give the server a message. It sleeps for the given number of seconds between
// checks.

type RandomFeed struct {
	id    uint
	name  string
	sleep time.Duration
	limit int
}

//---------------------------------------------------------------------

func NewRandomFeed(feed *orm.Feed) (*RandomFeed, error) {
	var _ FeedRunner = &RandomFeed{}

	f := &RandomFeed{
		id:   feed.ID,
		name: feed.Name,
	}

	err := f.setVars(feed.Config)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *RandomFeed) setVars(m map[string]interface{}) error {
	var err error

	f.sleep, err = common.GetMapValueAsDuration(m, "sleep")
	if err != nil {
		return err
	}

	f.limit, err = common.GetMapValueAsInt(m, "limit")
	if err != nil {
		return err
	}

	return nil
}

func (rf *RandomFeed) ID() uint {
	return rf.id
}

func (rf *RandomFeed) Name() string {
	return rf.name
}

func (f *RandomFeed) Run(statusF StatusF, mssgF MssgF) error {

	var err error
	ok := true

	go func() {
		for {
			ok, err = statusF("good")
			if err != nil {
				return
			}
			if !ok {
				return
			}

			x := rand.Intn(f.limit)
			if x == 0 {
				m := map[string]string{
					"mssg":  time.Now().Format(time.RFC3339),
					"limit": strconv.Itoa(x),
				}

				err = mssgF(m)
				if err != nil {
					return
				}
			}

			time.Sleep(f.sleep)
		}
	}()

	return nil
}
