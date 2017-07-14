package common

import "strconv"

var globalId int

type Ident string

const NoIdent Ident = ""

func NewId() Ident {
	globalId++
	s := strconv.Itoa(globalId)
	return Ident(s)
}

func (id Ident) String() string {
	return string(id)
}

func ToIdent(s string) Ident {
	return Ident(s)
}
