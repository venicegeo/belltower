package server

import (
	"log"

	"github.com/venicegeo/belltower/orm"
)

func Server() error {

	model, err := orm.NewOrm()
	if err != nil {
		return err
	}

	//
	// Usr
	//
	userAttrs := &orm.UserAttributes{
		Name:      "alice",
		IsAdmin:   false,
		IsEnabled: true,
	}
	userId, err := model.AddUser(userAttrs)
	if err != nil {
		return err
	}
	log.Printf("userid: %d", userId)

	//
	// Feed
	//
	feedAttrs := &orm.FeedAttributes{
		Name:      "randomfeed",
		FeedType:  "RandomFeed",
		IsEnabled: true,
		//PersistenceDuration: 0,
		ConfigInfo: map[string]string{
			"name":  "randomfeed2",
			"limit": "3",
			"sleep": "5",
		},
	}
	feedId, err := model.AddFeed(feedAttrs)
	log.Printf("feedid: %d", feedId)

	//
	// Action
	//

	//
	// Rule
	//

	// time.Sleep(10 * time.Second)

	return nil
}
