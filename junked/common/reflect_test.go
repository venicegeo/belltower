package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReflection(t *testing.T) {
	assert := assert.New(t)

	type T struct {
		I int
		S string
	}

	{
		i := 9

		x := NewViaReflection(i)
		ii := x.(int)
		assert.Equal(0, ii)

		ii = 4
		assert.Equal(4, ii)
		assert.Equal(9, i)
	}
	{
		i := 9

		x := NewViaReflection(&i)
		ii := x.(*int)
		assert.Equal(0, *ii)

		*ii = 4
		assert.Equal(4, *ii)
		assert.Equal(9, i)
	}

	{
		a := T{}
		aa := NewViaReflection(a)

		aaa, ok := aa.(T)
		assert.True(ok)
		assert.Equal(aaa.S, "")
		assert.Equal(aaa.I, 0)
		aaa.S = "sss"
		aaa.I = 99
		assert.Equal(aaa.S, "sss")
		assert.Equal(aaa.I, 99)
	}
	{
		a := &T{}
		aa := NewViaReflection(a)

		aaa, ok := aa.(*T)
		assert.True(ok)
		assert.Equal(aaa.S, "")
		assert.Equal(aaa.I, 0)
		aaa.S = "sss"
		aaa.I = 99
		assert.Equal(aaa.S, "sss")
		assert.Equal(aaa.I, 99)
	}
	//log.Printf("%T %#v", aa, aa)
}
