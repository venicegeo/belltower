package engine

import (
	"fmt"
	"strings"
	"text/scanner"

	"github.com/venicegeo/belltower/mpg/mlog"
)

type Token struct {
	typ int
	pos scanner.Position
	str string
}

func (t Token) isEOL() bool {
	return t.typ == 10 && t.str == "EOL"
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

	graph       *GraphModel
	inGraph     bool
	component   *ComponentModel
	inComponent bool
	connection  *ConnectionModel
	inConn      bool
}

// pop from front
func (t *Tokenizer) Pop() Token {
	x := t.tokens[t.tokenIndex]
	t.tokenIndex++
	return x
}

// push onto end
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
			t.Push(10, p, "EOL")
			lastLine = curLine
		}
		t.Push(int(tok), s.Pos(), s.TokenText())
	}

	for _, v := range t.tokens {
		fmt.Printf("%d  %s  %s\n", v.typ, v.pos, v.str)
	}
	return nil
}

func (t *Tokenizer) Parse() error {
	var err error

	for {

		token := t.Pop()
		switch {
		case token.isEOL():
			break
		case token.isIdentN("graph"):
			err = t.actionGraph(token)
		default:
			err = t.actionError(token)
			break
		}

		if err != nil {
			break
		}
	}
	return err
}

func (t *Tokenizer) actionError(tok Token) error {
	fmt.Printf("ERROR: %s", tok)
	return fmt.Errorf("ERROR: %s\n", tok)
}

func (t *Tokenizer) actionGraph(tok Token) error {
	mlog.Debug("at actionGraph")
	var err error

	t.graph = &GraphModel{}
	t.inGraph = true

	token := t.Pop()
	if !token.isIdent() {
		err = t.actionError(token)
		return err
	}
	t.actionGraphName(token)

	err = t.skipEOLs()

	switch {
	case token.isIdentN("component"):
		t.actionComponent(token)

	}
}
