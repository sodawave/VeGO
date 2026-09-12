package ast_test

import (
	"testing"

	"github.com/sodawave/VeGO/src/pkg/ast"
)

func TestKeywordMapBijective(t *testing.T) {
	for k, g := range ast.KeywordMap {
		got, ok := ast.ReverseMap[g]
		if !ok || got != k {
			t.Fatalf("glyph %q for %q reverse=%q ok=%v", g, k, got, ok)
		}
	}
}

func TestDefaultSymbols(t *testing.T) {
	var s ast.DefaultSymbols
	enc, ok := s.Encode("func")
	if !ok || enc != "ƒ" {
		t.Fatalf("Encode(func)=%q %v", enc, ok)
	}
	dec, ok := s.Decode("ƒ")
	if !ok || dec != "func" {
		t.Fatalf("Decode(ƒ)=%q %v", dec, ok)
	}
}

func TestAlphaKindsNonEmpty(t *testing.T) {
	if len(ast.AlphaKinds) == 0 || len(ast.OutOfAlphaKinds) == 0 {
		t.Fatal("Alpha kind lists must be populated")
	}
}
