package feeders

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/venicegeo/belltower/btorm"
	"github.com/venicegeo/belltower/common"
)


//---------------------------------------------------------------------

type Controller struct {
}

func (c *Controller) Execute() error {

	feederFactory := NewFeederFactory(&RandomFeeder{}, &FileSysFeeder{})

	// read feed instances from DB
	orm := btorm.NewBtOrm()
	var currentFeeds []*btorm.Feed := orm.getallfeeds()

	// for each feed in the DB, make an instance of the correct feeder type

	currentFeeders := []*Feeder{}

	for _, feed := range currentFeeds {
		feeder := feederFactory.Create(feed.settings)

		currentFeeders = append(currentFeeders, feeder)
	}

	// launch the feeders
	for _,feeder := range currentFeeders {
		go c.RunLoop(feeder)
	}

	// block
}

func (c *Controller) RunLoop(feeder *Feeder) {
	errCount := 0

	for {
		if errCount == 3 {
			log.Printf("Feeder aborting (%s)", feeder.GetName())
			break
		}

		event, err := feeder.Poll()
		if err != nil {
			log.Printf("Feeder poll error (%s): %v", feeder.GetName(), err)
			errCount++
		} else {
			if event != nil {
				err = c.Post(event)
				log.Printf("Feeder post error (%s): %v", feeder.GetName(), err)
				// TODO
			}
		}

		time.Sleep(settings.GetInterval())
	}
}

func (c *Controller) Post(e *FeedEvent) error {
	log.Printf("POSTED %T %#v", e, e.Data.(*RandomEventData))
	byts, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	log.Printf("%s", string(byts))

	fe := FeedEvent{}
	err = json.Unmarshal(byts, &fe)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v", fe)
	return nil
}

//---------------------------------------------------------------------
