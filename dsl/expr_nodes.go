package dsl

import (
	"fmt"
	"strconv"
)

//---------------------------------------------------------------------

type ExprValue struct {
	Type  BaseType
	Value interface{}
}

type ExprNode interface {
	String() string
	Eval(*Environment) ExprValue
}

//---------------------------------------------------------------------

type ExprNodeMultiply struct {
	left  ExprNode
	right ExprNode
}

func NewExprNodeMultiply(left ExprNode, right ExprNode) *ExprNodeMultiply {
	n := &ExprNodeMultiply{
		left:  left,
		right: right,
	}
	return n
}

func (n *ExprNodeMultiply) Eval(env *Environment) ExprValue {
	l := n.left.Eval(env)
	r := n.right.Eval(env)

	if l.Type != r.Type {
		panic(19)
	}
	switch l.Type {
	case IntType:
		return ExprValue{Type: IntType, Value: l.Value.(int) * r.Value.(int)}
	case FloatType:
		return ExprValue{Type: FloatType, Value: l.Value.(float64) * r.Value.(float64)}
	}
	panic(19)
}

func (n *ExprNodeMultiply) String() string {
	return fmt.Sprintf("(* %v %v)", n.left, n.right)
}

//---------------------------------------------------------------------

type ExprNodeAdd struct {
	left  ExprNode
	right ExprNode
}

func NewExprNodeAdd(left ExprNode, right ExprNode) *ExprNodeAdd {
	n := &ExprNodeAdd{
		left:  left,
		right: right,
	}
	return n
}

func (n *ExprNodeAdd) Eval(env *Environment) ExprValue {
	l := n.left.Eval(env)
	r := n.right.Eval(env)

	if l.Type != r.Type {
		panic(19)
	}
	switch l.Type {
	case IntType:
		return ExprValue{Type: IntType, Value: l.Value.(int) + r.Value.(int)}
	case FloatType:
		return ExprValue{Type: FloatType, Value: l.Value.(float64) + r.Value.(float64)}
	}
	panic(19)
}

func (n *ExprNodeAdd) String() string {
	return fmt.Sprintf("(+ %v %v)", n.left, n.right)
}

//---------------------------------------------------------------------------

type ExprNodeSymbolRef struct {
	name string
}

func NewExprNodeSymbolRef(name string) *ExprNodeSymbolRef {
	n := &ExprNodeSymbolRef{
		name: name,
	}
	return n
}

func (n *ExprNodeSymbolRef) Eval(env *Environment) ExprValue {
	return env.get(n.name)
}

func (n *ExprNodeSymbolRef) String() string {
	return fmt.Sprintf("SYMBOLREF(%s)", n.name)
}

//---------------------------------------------------------------------------

type ExprNodeIntConstant struct {
	value ExprValue
}

func NewExprNodeIntConstant(value int) *ExprNodeIntConstant {
	n := &ExprNodeIntConstant{
		value: ExprValue{Type: IntType, Value: value},
	}
	return n
}

func (n *ExprNodeIntConstant) Eval(env *Environment) ExprValue {
	return n.value
}

func (n *ExprNodeIntConstant) String() string {
	return fmt.Sprintf("INTCONSTANT(%d)", n.value)
}

//---------------------------------------------------------------------------

type ExprNodeFloatConstant struct {
	value ExprValue
}

func NewExprNodeFloatConstant(value float64) *ExprNodeFloatConstant {
	n := &ExprNodeFloatConstant{
		value: ExprValue{Type: FloatType, Value: value},
	}
	return n
}

func (n *ExprNodeFloatConstant) Eval(env *Environment) ExprValue {
	return n.value
}

func (n *ExprNodeFloatConstant) String() string {
	return fmt.Sprintf("FLOATCONSTANT(%f)", n.value.Value.(float64))
}

//---------------------------------------------------------------------------

type ExprNodeBoolConstant struct {
	value ExprValue
}

func NewExprNodeBoolConstant(value bool) *ExprNodeBoolConstant {
	n := &ExprNodeBoolConstant{
		value: ExprValue{Type: BoolType, Value: value},
	}
	return n
}

func (n *ExprNodeBoolConstant) Eval(env *Environment) ExprValue {
	return n.value
}

func (n *ExprNodeBoolConstant) String() string {
	return fmt.Sprintf("BOOLCONSTANT(%t)", n.value.Value.(bool))
}

//---------------------------------------------------------------------------

type ExprNodeStringConstant struct {
	value ExprValue
}

func NewExprNodeStringConstant(value string) *ExprNodeStringConstant {
	n := &ExprNodeStringConstant{
		value: ExprValue{Type: StringType, Value: value},
	}
	return n
}

func (n *ExprNodeStringConstant) Eval(env *Environment) ExprValue {
	return n.value
}

func (n *ExprNodeStringConstant) String() string {
	return fmt.Sprintf("STRINGCONSTANT(%s)", n.value.Value.(string))
}

//---------------------------------------------------------------------------

// "message.timestamp"
type ExprNodeStructRef struct {
	symbol string
	field  string
}

func NewExprNodeStructRef(symbol string, field string) *ExprNodeStructRef {
	n := &ExprNodeStructRef{
		symbol: symbol,
		field:  field,
	}
	return n
}

func (n *ExprNodeStructRef) Eval(env *Environment) ExprValue {
	return env.get(n.symbol + "." + n.field)
}

func (n *ExprNodeStructRef) String() string {
	return fmt.Sprintf("STRUCTREF(%s.%s)", n.symbol, n.field)
}

//---------------------------------------------------------------------------

type ExprNodeArrayRef struct {
	symbol string
	index  ExprNode
}

func NewExprNodeArrayRef(symbol string, index ExprNode) *ExprNodeArrayRef {
	n := &ExprNodeArrayRef{
		symbol: symbol,
		index:  index,
	}
	return n
}

func (n *ExprNodeArrayRef) Eval(env *Environment) ExprValue {
	idx := n.index.Eval(env)
	if idx.Type != IntType {
		panic(19)
	}
	idxstr := strconv.Itoa(idx.Value.(int))
	return env.get(n.symbol + "#" + idxstr)
}

func (n *ExprNodeArrayRef) String() string {
	return fmt.Sprintf("ARRAYREF(%s,%v)", n.symbol, n.index)
}

//---------------------------------------------------------------------------

type ExprNodeMapRef struct {
	symbol string
	key    ExprNode
}

func NewExprNodeMapRef(symbol string, key ExprNode) *ExprNodeMapRef {
	n := &ExprNodeMapRef{
		key: key,
	}
	return n
}

func (n *ExprNodeMapRef) String() string {
	return fmt.Sprintf("MAPREF(%s,%v)", n.symbol, n.key)
}

//---------------------------------------------------------------------------
