package expr

// expression parser started life from https://thorstenball.com/blog/2016/11/16/putting-eval-in-go/

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

func eval(expr ast.Expr) (int, error) {
	switch typedExpr := expr.(type) {
	case *ast.BinaryExpr:
		return evalBinaryExpr(typedExpr)
	case *ast.BasicLit:
		switch typedExpr.Kind {
		case token.INT:
			return strconv.Atoi(typedExpr.Value)
		default:
			return 0, fmt.Errorf("ast.BasicLit kind not supported: %s [%s]", typedExpr.Kind.String(), typedExpr.Value)
		}
	default:
		return 0, fmt.Errorf("ast.Expr type not supported: %s [%s]", expr, typedExpr)
	}

	// notreached
}

func evalBinaryExpr(expr *ast.BinaryExpr) (int, error) {
	left, err := eval(expr.X)
	if err != nil {
		return 0, err
	}
	right, err := eval(expr.Y)
	if err != nil {
		return 0, err
	}

	switch expr.Op {
	case token.ADD:
		return left + right, nil
	case token.SUB:
		return left - right, nil
	case token.MUL:
		return left * right, nil
	case token.QUO:
		return left / right, nil
	default:
		return 0, fmt.Errorf("token Op type not supported: %s", expr.Op.String())
	}

	// notreached
}
