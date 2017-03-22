package orm2

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrm(t *testing.T) {
	assert := assert.New(t)

	orm, err := NewOrm()
	assert.NoError(err)
	assert.NotNil(orm)

	now := time.Now()

	feed := &Feed{
		Name: "Bob",
		Time: now,
		Bool: true,
	}

	id, err := orm.Create(feed)
	assert.NoError(err)
	assert.NotEmpty(id)

	f := &Feed{Id: id}
	g, err := orm.Read(f)
	assert.NoError(err)
	assert.NotNil(g)
	log.Printf("%#v", g)
	assert.EqualValues(id, g.GetID())

	assert.EqualValues("Bob", feed.Name)
	assert.EqualValues(feed.Name, g.(*Feed).Name)
	assert.EqualValues(true, feed.Bool)
	assert.EqualValues(feed.Bool, g.(*Feed).Bool)
	assert.EqualValues(now, feed.Time)
	assert.EqualValues(feed.Time, g.(*Feed).Time)
}
