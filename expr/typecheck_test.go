package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSameType(t *testing.T) {
	assert := assert.New(t)

	assert.True(sameType("a", "1"))
	assert.False(sameType("a", 1))
}

func TestTypeCheck(t *testing.T) {
	assert := assert.New(t)

	type S struct {
		W string
		X int
		Y string
	}
	a := &S{
		W: "www",
		X: 19,
		Y: "yyy",
	}

	b := map[string]interface{}{
		"W": 12.34,
		"X": 5,
		"Z": true,
	}

	aa, errs := TypeCheck(a, b)

	assert.Len(errs, 3)
	//log.Printf(errs[0])
	//log.Printf(errs[1])
	//log.Printf(errs[2])

	assert.Equal("www", aa.(*S).W)
	assert.Equal(5, aa.(*S).X)
	assert.Equal("yyy", aa.(*S).Y)

	//    t.w = ""
	//    t.x = 5
	//    t.y = ""
	// with these error messages:
	//    t.w expected string but got float
	//    t.y not set
	//    b.z not used

}
