package dsl

import "fmt"

//---------------------------------------------------------------------------

type TokenId int

const (
	TokenInvalid TokenId = iota
	TokenEquals
	TokenNotEquals
	TokenGreaterThan
	TokenGreaterOrEqualThan
	TokenLessThan
	TokenLessOrEqualThan
	TokenAdd
	TokenSubtract
	TokenMultiply
	TokenDivide
	TokenExponent
	TokenMod
	TokenBitwiseOr
	TokenBitwiseAnd
	TokenLogicalAnd
	TokenLogicalOr
	TokenLeftParen
	TokenRightParen
	TokenLeftBracket
	TokenRightBracket
	TokenSymbol // 16
	TokenNumber
	TokenTypeSlice
	TokenTypeArray
	TokenTypeMap
)

type Token struct {
	Line   int
	Column int
	Text   string
	Id     TokenId
	Value  interface{}
}

func (t *Token) String() string {
	s := fmt.Sprintf("[%d:%d] id=%d text=\"%s\"", t.Line, t.Column, t.Id, t.Text)
	if t.Value != nil {
		s += fmt.Sprintf(" value=<%v>", t.Value)
	}
	return s
}

func convertId(r rune) TokenId {
	switch r {
	case -2:
		return TokenSymbol
	case -3:
		return TokenNumber
	case 40:
		return TokenLeftParen
	case 41:
		return TokenRightParen
	case 42:
		return TokenMultiply
	case 43:
		return TokenAdd
	case 60:
		return TokenLessThan
	case 62:
		return TokenGreaterThan
	case 91:
		return TokenLeftBracket
	case 93:
		return TokenRightBracket
	case 124:
		return TokenBitwiseOr
	default:
		return TokenInvalid
	}
}
