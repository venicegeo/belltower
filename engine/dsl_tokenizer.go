package engine

import (
	"fmt"
	"strings"
	"text/scanner"
)

type Token struct {
	typ int
	pos scanner.Position
	str string
}

const (
	ColonMarker   = 58
	PeriodMarker  = 46
	HyphenMarker  = 45
	GreaterMarker = 62
)

func (t Token) isEOL() bool {
	return t.typ == 10 && t.str == "EOL"
}

func (t Token) isEOF() bool {
	return t.typ == -1
}

func (t Token) isIdent() bool {
	return t.typ == -2
}

func (t Token) isIdentN(s string) bool {
	return t.isIdent() && t.str == s
}

func (t Token) String() string {
	return fmt.Sprintf("%s: %s (%d)", t.pos, t.str, t.typ)
}

//---------------------------------------------------------------------

type Tokenizer struct {
	tokens     []Token
	tokenIndex int
}

// pop from front
func (t *Tokenizer) Pop() Token {
	x := t.tokens[t.tokenIndex]
	t.tokenIndex++
	return x
}

// return top, but don't pop
func (t *Tokenizer) Peek() Token {
	x := t.tokens[t.tokenIndex]
	return x
}

// push onto front!
func (t *Tokenizer) PutBack(token Token) {
	t.tokens = append([]Token{token}, t.tokens...)
}

// push onto end
func (t *Tokenizer) Push(token Token) {

	t.tokens = append(t.tokens, token)
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
			t.Push(Token{typ: 10, pos: p, str: "EOL"})
			lastLine = curLine
		}
		t.Push(Token{typ: int(tok), pos: s.Pos(), str: s.TokenText()})
	}

	//t.Push(Token{typ: 0, pos: scanner.Position{}, str: "EOF"})

	for _, v := range t.tokens {
		fmt.Printf("%d  %s  %s\n", v.typ, v.pos, v.str)
	}
	return nil
}
