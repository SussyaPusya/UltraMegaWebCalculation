package pkg

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
)

func validate(node ast.Node) error {
	switch n := node.(type) {
	case *ast.ParenExpr:
		return validate(n.X)
	case *ast.BasicLit:
		return nil
	case *ast.BinaryExpr:
		if n.Op != token.ADD && n.Op != token.SUB && n.Op != token.MUL && n.Op != token.QUO {
			return errors.New("invalid binary operator")
		}
		if err := validate(n.X); err != nil {
			return err
		}
		if err := validate(n.Y); err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid expression")
	}
}

func ParseAst(expression string) (ast.Node, error) {
	node, err := parser.ParseExpr(expression)
	if err != nil {
		return nil, err
	}
	if err := validate(node); err != nil {
		return nil, err
	}
	return node, nil
}
