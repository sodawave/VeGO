package transpiler_test

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sodawave/VeGO/src/pkg/transpiler"
)

func TestRoundTripFixtures(t *testing.T) {
	eng := transpiler.New(nil)
	files := []string{"hello_http.go", "struct_range.go"}
	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			goSrc, err := os.ReadFile(filepath.Join("testdata", name))
			if err != nil {
				t.Fatal(err)
			}
			canon, err := format.Source(goSrc)
			if err != nil {
				t.Fatal(err)
			}
			vego, err := eng.Transform(canon, transpiler.ToVeGo)
			if err != nil {
				t.Fatalf("ToVeGo: %v", err)
			}
			if bytes.Equal(vego, canon) {
				t.Fatal("expected .vego to differ from Go (keywords should be glyphs)")
			}
			if !bytes.Contains(vego, []byte("¾")) && !bytes.Contains(vego, []byte("æ")) && !bytes.Contains(vego, []byte("†")) {
				t.Fatalf("expected glyph keywords in .vego, got: %s", vego)
			}
			if bytes.Contains(vego, []byte{'\n'}) {
				t.Fatalf(".vego must be single-line (no newlines), got: %q", vego)
			}
			back, err := eng.Transform(vego, transpiler.ToGo)
			if err != nil {
				t.Fatalf("ToGo: %v\nvego=%q", err, vego)
			}
			// Formatting may differ after single-line expand; compare syntax trees.
			want, err := astFingerprint(canon)
			if err != nil {
				t.Fatal(err)
			}
			got, err := astFingerprint(back)
			if err != nil {
				t.Fatal(err)
			}
			if want != got {
				t.Fatalf("AST mismatch\n--- go ---\n%s\n--- back ---\n%s\n--- vego ---\n%q\n--- fp go ---\n%s\n--- fp back ---\n%s",
					canon, back, vego, want, got)
			}
		})
	}
}

func TestInvalidVeGoFailsClosed(t *testing.T) {
	eng := transpiler.New(nil)
	_, err := eng.Transform([]byte("ð not valid {{{"), transpiler.ToGo)
	if err == nil {
		t.Fatal("expected error for invalid .vego")
	}
}

func TestInvalidGoFailsClosed(t *testing.T) {
	eng := transpiler.New(nil)
	_, err := eng.Transform([]byte("package !!!"), transpiler.ToVeGo)
	if err == nil {
		t.Fatal("expected error for invalid Go")
	}
}

// astFingerprint is a position-independent syntax dump for round-trip checks.
func astFingerprint(src []byte) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", src, 0)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case nil:
			return false
		case *ast.Ident:
			fmt.Fprintf(&b, "IDENT %s;", x.Name)
		case *ast.BasicLit:
			fmt.Fprintf(&b, "LIT %s;", x.Value)
		case *ast.BinaryExpr:
			fmt.Fprintf(&b, "BIN %v;", x.Op)
		case *ast.UnaryExpr:
			fmt.Fprintf(&b, "UN %v;", x.Op)
		case *ast.AssignStmt:
			fmt.Fprintf(&b, "ASSIGN %v;", x.Tok)
		case *ast.GenDecl:
			fmt.Fprintf(&b, "GEN %v;", x.Tok)
		case *ast.BranchStmt:
			fmt.Fprintf(&b, "BRANCH %v;", x.Tok)
		case *ast.IncDecStmt:
			fmt.Fprintf(&b, "INCDEC %v;", x.Tok)
		case *ast.Comment, *ast.CommentGroup:
			return true
		default:
			fmt.Fprintf(&b, "%T;", n)
		}
		return true
	})
	return b.String(), nil
}
