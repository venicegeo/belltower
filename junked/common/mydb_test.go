package common

import (
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestSql(t *testing.T) {
	assert := assert.New(t)

	var err error
	var ok bool

	file := "./test.db"
	os.Remove(file)

	mydb := &MyDB{Name: file}
	err = mydb.Open()
	assert.NoError(err)

	// does HasTable() work?
	{
		ok, err := mydb.HasTable("foo")
		assert.NoError(err)
		assert.False(ok)

		sqlStmt := `create table foo (id integer not null primary key, name text);`
		_, err = mydb.DB.Exec(sqlStmt)
		assert.NoError(err)

		ok, err = mydb.HasTable("foo")
		assert.NoError(err)
		assert.True(ok)
	}

	// if we reopen, is table foo still there?
	{
		err = mydb.Close()
		assert.NoError(err)

		err = mydb.Open()
		assert.NoError(err)

		ok, err = mydb.HasTable("foo")
		assert.NoError(err)
		assert.True(ok)
	}

	// can we delete the table?
	{
		err = mydb.DropTable("foo")
		assert.NoError(err)

		ok, err = mydb.HasTable("foo")
		assert.NoError(err)
		assert.False(ok)
	}

	err = mydb.Close()
	assert.NoError(err)
	assert.Nil(mydb.DB)
}
