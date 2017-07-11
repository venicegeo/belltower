package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	//assert := assert.New(t)
}
func TestStack(t *testing.T) {
	assert := assert.New(t)

	r := NewTokenRParen(1, 1)
	l := NewTokenLParen(2, 2)

	s := Stack{}
	s.push(r)
	assert.Equal(r, s.peek())
	s.push(l)
	assert.Equal(l, s.peek())
	assert.Equal(l, s.pop())
	assert.Equal(r, s.pop())
	assert.Nil(s.pop())
	assert.Nil(s.pop())
}

func TestParse(t *testing.T) {
	//assert := assert.New(t)

	PARSE("8*x")
}
