package transpiler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

// normalizeAssignOps rewrites x = x + y into x += y (and -, *, /, %).
// Equivalent Go; denser under tiktoken.
func normalizeAssignOps(src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	ast.Inspect(f, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || as.Tok != token.ASSIGN || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
			return true
		}
		bin, ok := as.Rhs[0].(*ast.BinaryExpr)
		if !ok || !exprEq(as.Lhs[0], bin.X) {
			return true
		}
		var tok token.Token
		switch bin.Op {
		case token.ADD:
			tok = token.ADD_ASSIGN
		case token.SUB:
			tok = token.SUB_ASSIGN
		case token.MUL:
			tok = token.MUL_ASSIGN
		case token.QUO:
			tok = token.QUO_ASSIGN
		case token.REM:
			tok = token.REM_ASSIGN
		default:
			return true
		}
		as.Tok = tok
		as.Rhs = []ast.Expr{bin.Y}
		return true
	})
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return nil, fmt.Errorf("normalize: %w", err)
	}
	return buf.Bytes(), nil
}

func exprEq(a, b ast.Expr) bool {
	ia, oka := a.(*ast.Ident)
	ib, okb := b.(*ast.Ident)
	return oka && okb && ia.Name == ib.Name
}
