package engine

import (
	"fmt"
	"strings"
	"text/scanner"
)

type Tokenizer struct {
	tokens []Token
}

type Token struct {
	typ int
	pos scanner.Position
	str string
}

func (t *Tokenizer) Push(tokenType int, pos scanner.Position, tokenString string) {
	tok := Token{
		typ: tokenType,
		pos: pos,
		str: tokenString,
	}

	t.tokens = append(t.tokens, tok)
}

func (t *Tokenizer) Scan(lines string) error {

	t.tokens = []Token{}

	var s scanner.Scanner
	s.Filename = "example"
	s.Init(strings.NewReader(lines))

	var tok rune
	lastLine := 0
	for tok != scanner.EOF {
		tok = s.Scan()
		curLine := s.Pos().Line
		if curLine != lastLine {
			p := s.Pos()
			p.Line -= 1
			t.Push(10, p, "\n")
			lastLine = curLine
		}
		t.Push(int(tok), s.Pos(), s.TokenText())
	}

	for _, v := range t.tokens {
		s := v.str
		if v.typ == 10 {
			s = "\\n"
		}
		fmt.Printf("%d  %s  %s\n", v.typ, v.pos, s)
	}
	return nil
}

func (t *Tokenizer) Parse() error {

	//	currGraph := &GraphModel{}
	//	var currComponent ComponentModel
	//	var currConnection ConnectionModel

	return nil
}
