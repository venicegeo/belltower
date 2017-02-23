package dsl

import "strconv"

type Id string

var currId int

const InvalidId Id = ""

func newId() Id {
	currId++
	return Id(strconv.Itoa(currId))
}

//---------------------------------------------------------------------------

type Dsl struct {
	TypeTables map[Id]*TypeTable
	Exprs      map[Id]ExprNode
}

func NewDsl() (*Dsl, error) {
	d := &Dsl{}
	d.TypeTables = map[Id]*TypeTable{}
	d.Exprs = map[Id]ExprNode{}
	return d, nil
}

func (d *Dsl) ParseDeclaration(decl string) (Id, error) {
	var err error

	tt, err := NewTypeTokenizer()
	if err != nil {
		return InvalidId, err
	}

	typeTable, err := tt.ParseJson(decl)
	if err != nil {
		return InvalidId, err
	}

	id := newId()
	d.TypeTables[id] = typeTable

	return id, nil
}

func (d *Dsl) ParseExpression(expr string) (Id, error) {
	var err error

	et := &ExprTokenizer{}
	toks, err := et.Tokenize(expr)
	if err != nil {
		return InvalidId, err
	}

	ep := &ExprParser{}
	node, err := ep.Parse(toks)
	if err != nil {
		return InvalidId, err
	}

	id := newId()
	d.Exprs[id] = node

	return id, nil
}

func (d *Dsl) Evaluate(exprId Id, typeId Id, env *Environment) (interface{}, error) {
	eval := &Eval{}
	result, err := eval.Evaluate(d.Exprs[exprId], env)
	if err != nil {
		return nil, err
	}

	return result, nil
}
