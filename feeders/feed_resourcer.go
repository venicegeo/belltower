package feeders

import (
	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/orm"
)

type FeedResourcer struct {
	orm *orm.Orm
}

func (fr *FeedResourcer) Create(jstr common.JSON) (uint, error) {
	feed := &orm.Feed{}
	err := jstr.ToObject(feed)
	if err != nil {
		return 0, err
	}

	// now make the model
	id, err := fr.orm.AddFeed(feed)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (fr *FeedResourcer) Read(id uint) (common.JSON, error) {
	feed, err := fr.orm.GetFeed(id)
	if err != nil {
		return common.NilJSON, err
	}
	s, err := common.ToJson(feed)
	if err != nil {
		return common.NilJSON, err
	}
	return s, nil
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
