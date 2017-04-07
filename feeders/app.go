package feeders

import (
	"log"
	"time"

	"github.com/venicegeo/belltower/btorm"
)

type App struct {
	orm   *btorm.BtOrm
	feeds []*btorm.Feed
}

func NewApp() (*App, error) {

	orm, err := btorm.NewBtOrm("", btorm.OrmOptionOpen)
	if err != nil {
		return nil, err
	}

	app := &App{
		orm: orm,
	}

	return app, nil
}

func (app *App) Load() error {
	feeds, count, err := app.orm.ReadFeeds(0, 10)
	if err != nil {
		return err
	}

	if count > int64(len(feeds)) {
		panic("TODO")
	}

	log.Printf("got %d %d feeds", count, len(feeds))

	app.feeds = feeds

	return nil
}

func (app *App) Run() error {

	poster := func(event *Event) error {
		log.Printf("HIT: %#v", event)
		return nil
	}

	for _, feed := range app.feeds {

		feeder, err := feederFactory.create(feed)
		if err != nil {
			return err
		}

		go RunFeed(feed, feeder, poster)
		if err != nil {
			return err
		}
	}

	time.Sleep(10 * time.Second)
	return nil
}
