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
	userAttrs := &orm.UserFieldsForCreate{
		Name:      "alice",
		Role:      orm.AdminRole,
		IsEnabled: true,
	}
	userId, err := model.CreateUser(model.AdminID, userAttrs)
	if err != nil {
		return err
	}
	log.Printf("userid: %d", userId)

	//
	// Feed
	//
	feedAttrs := &orm.FeedFieldsForCreate{
		Name:      "randomfeed",
		FeedType:  "RandomFeed",
		IsEnabled: true,
		IsPublic:  false,
		Settings: map[string]interface{}{
			"name":  "randomfeed2",
			"limit": "3",
			"sleep": "5",
		},
	}
	feedId, err := model.CreateFeed(userId, feedAttrs)
	if err != nil {
		return err
	}
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
