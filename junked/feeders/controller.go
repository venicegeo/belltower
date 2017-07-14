package feeders

import (
	"context"
	"log"
	"time"

	"fmt"

	"github.com/venicegeo/belltower/btorm"
	"github.com/venicegeo/belltower/common"
)

type Controller struct {
	Post         EventPosterFunc
	orm          *btorm.BtOrm
	feeds        []*btorm.Feed
	ctx          context.Context
	cancel       context.CancelFunc
	ctxMap       map[common.Ident]context.Context    // one for each feeder id
	ctxCancelMap map[common.Ident]context.CancelFunc // one for each feeder id
}

func (ctl *Controller) load() error {
	// TODO: just read 100 for now
	feeds, _, err := ctl.orm.ReadFeeds(0, 100)
	if err != nil {
		return err
	}

	log.Printf("got %d feeds", len(feeds))

	ctl.feeds = feeds

	return nil
}

// Start looks at each Feed from the database and makes a Feeder for it.
// The Feeder runs in a goroutine. When all Feeders have been made, the function
// returns.
func (ctl *Controller) Start() error {

	if ctl.Post == nil {
		ctl.Post = func(event *Event) error {
			log.Printf("HIT: %#v", event)
			return nil
		}
	}

	orm := &btorm.BtOrm{}
	err := orm.Open()
	if err != nil {
		return err
	}

	ctl.orm = orm

	if err = ctl.load(); err != nil {
		return err
	}

	ctl.ctx, ctl.cancel = context.WithCancel(context.Background())
	ctl.ctxMap = map[common.Ident]context.Context{}
	ctl.ctxCancelMap = map[common.Ident]context.CancelFunc{}

	for _, feed := range ctl.feeds {

		feederInfo := feederRegistry.data[feed.FeederId]
		if feederInfo == nil {
			return fmt.Errorf("feed's feeder not registered")
		}

		feeder, err := feederInfo.Create(feed)
		if err != nil {
			return err
		}

		var ctx context.Context
		var cancel context.CancelFunc
		if feed.PollingEndAt.IsZero() {
			ctx, cancel = context.WithCancel(ctl.ctx)
		} else {
			ctx, cancel = context.WithDeadline(ctl.ctx, feed.PollingEndAt)
		}
		ctl.ctxMap[feed.FeederId] = ctx
		ctl.ctxCancelMap[feed.FeederId] = cancel

		go runFeed(ctx, feed, feeder, ctl.Post)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ctl *Controller) Stop() error {
	ctl.cancel()
	return nil
}

func (ctl *Controller) Status() error {
	panic(99)
}

func runFeed(ctx context.Context, feed *btorm.Feed, feeder Feeder, post EventPosterFunc) error {
	errCount := 0

	for {
		eventData, err := feeder.Poll()

		errored := false

		if err != nil {
			log.Printf("Feeder poll error (%s): %v", feeder.GetName(), err)
			errored = true
		} else if eventData == nil {
			// no error and no data - not a hit
		} else {
			// a hit! a hit!
			event := &Event{
				TimeStamp: time.Now(),
				FeedId:    feed.Id,
				FeederId:  feed.FeederId,
				Data:      eventData,
			}
			err = post(event)
			if err != nil {
				log.Printf("Feeder post error (%s): %v", feeder.GetName(), err)
				errored = true
			}
		}

		if errored {
			if errCount == 3 {
				log.Printf("Feeder aborting (%s)", feeder.GetName())
				return fmt.Errorf("Feeder aborting -- too many errors")
			}
		} else {
			errCount = 0
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		d := time.Duration(feed.PollingInterval)
		time.Sleep(d * time.Second)
	}

	// not reached
}
