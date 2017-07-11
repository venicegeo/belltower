package expr

// expression parser started life from https://thorstenball.com/blog/2016/11/16/putting-eval-in-go/

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

type Expression struct {
	text string
	tree ast.Expr
}

func NewExpression(text string) (*Expression, error) {
	tree, err := parser.ParseExpr(text)
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %s", err)
	}

	e := &Expression{
		text: text,
		tree: tree,
	}

	return e, nil
}

func (e *Expression) String() string {
	buf := bytes.Buffer{}
	printer.Fprint(&buf, token.NewFileSet(), e)
	return buf.String()
}

func (e *Expression) Eval() (int, error) {

	r, err := eval(e.tree)
	if err != nil {
		return 0, err
	}
	return r, nil
}
