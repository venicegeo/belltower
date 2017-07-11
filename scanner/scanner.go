package scanner

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strconv"
)

type Stack struct {
	array []Token
}

func (s *Stack) push(t Token) {
	s.array = append(s.array, t)
}

func (s *Stack) pop() Token {
	if len(s.array) == 0 {
		return nil
	}
	t := s.array[len(s.array)-1]
	s.array = s.array[:len(s.array)-1]
	return t
}

func (s *Stack) peek() Token {
	t := s.array[len(s.array)-1]
	return t
}

type TokenType int

const (
	TokenTypeInvalid TokenType = iota
	TokenTypeRParen
	TokenTypeLParen
	TokenTypeComma
	TokenTypeIdent
)

type Token interface {
	String() string
	StartPos() int
	EndPos() int
	Type() TokenType
}

type TokenCommon struct {
	typ      TokenType
	startPos int
	endPos   int
}

func NewTokenCommon(typ TokenType, startPos int, endPos int) (t TokenCommon) {
	return TokenCommon{
		startPos: startPos,
		endPos:   endPos,
		typ:      typ,
	}
}

func (t *TokenCommon) Type() TokenType {
	return t.typ
}

func (t *TokenCommon) StartPos() int {
	return t.startPos
}

func (t *TokenCommon) EndPos() int {
	return t.endPos
}

func (t *TokenCommon) String() string {
	return fmt.Sprintf("[%d:%d] %d", t.startPos, t.endPos, t.typ)
}

type TokenLParen struct {
	TokenCommon
}

func NewTokenLParen(startPos int, endPos int) (t *TokenLParen) {
	return &TokenLParen{
		TokenCommon: NewTokenCommon(TokenTypeLParen, startPos, endPos),
	}
}

type TokenRParen struct {
	TokenCommon
}

func NewTokenRParen(startPos int, endPos int) (t *TokenRParen) {
	return &TokenRParen{
		TokenCommon: NewTokenCommon(TokenTypeRParen, startPos, endPos),
	}
}

func Scan(expr string) {
	/*	var s scanner.Scanner
			s.Filename = "example"
			s.Init(strings.NewReader(expr))
			var tok rune
			for tok != scanner.EOF {
				tok = s.Scan()
				fmt.Println("At position", s.Pos(), ":", s.TokenText(), tok)
				switch tok {
		 case        Ident:
		        case Int
		        case Float
		        case Char
		        case String
		        case RawString
		        case Comment
				}
			}
	*/
}

// adapted from https://thorstenball.com/blog/2016/11/16/putting-eval-in-go/
func PARSE(line string) {
	exp, err := parser.ParseExpr(line)
	if err != nil {
		fmt.Printf("parsing failed: %s\n", err)
		return
	}

	printer.Fprint(os.Stdout, token.NewFileSet(), exp)
	fmt.Printf("\n")

	fmt.Printf("%d\n", Eval(exp))
}

func Eval(exp ast.Expr) int {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return EvalBinaryExpr(exp)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT:
			i, _ := strconv.Atoi(exp.Value)
			return i
		}
	}

	return 0
}

func EvalBinaryExpr(exp *ast.BinaryExpr) int {
	left := Eval(exp.X)
	right := Eval(exp.Y)

	switch exp.Op {
	case token.ADD:
		return left + right
	case token.SUB:
		return left - right
	case token.MUL:
		return left * right
	case token.QUO:
		return left / right
	}

	// fallthrough
	panic(exp)
}
