package mlog

import (
	"testing"

	"os"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

func TestMLog(t *testing.T) {
	assert := assert.New(t)

	const filename = "testmlog.log"
	os.Remove(filename)
	SetFlags(0)
	{
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
		assert.NoError(err)
		defer f.Close()
		SetOutput(f)

		Print("one")
		Printf("%s", "two")
		Debug("three")
		Debugf("%s", "four")
		Verbose = true
		Debug("five")
		Debugf("%s", "six")
		SetFlags(Lshortfile)
		Print("seven")
	}

	expected :=
		`one
two
five
six
mlog_test.go:33 seven
`
	byts, err := ioutil.ReadFile(filename)
	assert.NoError(err)
	assert.Equal(expected, string(byts))

	os.Remove(filename)

}
