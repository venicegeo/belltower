package feeders

import (
	"encoding/json"
	"log"
)

//---------------------------------------------------------------------

type Controller struct {
}

func (c *Controller) Post(e *Event) error {
	log.Printf("POSTED %T %#v", e, e.Data.(*RandomEventData))
	byts, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	log.Printf("%s", string(byts))

	fe := Event{}
	err = json.Unmarshal(byts, &fe)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v", fe)
	return nil
}

//---------------------------------------------------------------------
