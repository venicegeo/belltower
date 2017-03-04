package feeders

import (
	"github.com/venicegeo/belltower/orm"
)

type FeedResourcer struct {
	orm *orm.Orm
}

func (fr *FeedResourcer) Create(feed *orm.Feed) (uint, error) {
	id, err := fr.orm.AddFeed(feed)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (fr *FeedResourcer) Read(id uint) (*orm.Feed, error) {
	feed, err := fr.orm.GetFeed(id)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

func (fr *FeedResourcer) Update(id uint, feed *orm.Feed) error {
	err := fr.orm.UpdateFeed(id, feed)
	if err != nil {
		return err
	}
	return nil
}

func (fr *FeedResourcer) Delete(id uint) error {
	return fr.orm.DeleteFeed(id)
}
