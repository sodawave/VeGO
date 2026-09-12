package ast_test

import (
	"testing"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sodawave/VeGO/src/pkg/ast"
)

func TestMapsBijective(t *testing.T) {
	for k, g := range ast.KeywordMap {
		if got := ast.ReverseMap[g]; got != k {
			t.Fatalf("keyword %q glyph %q reverse=%q", k, g, got)
		}
	}
	for k, g := range ast.PhraseMap {
		if got := ast.ReverseMap[g]; got != k {
			t.Fatalf("phrase %q glyph %q reverse=%q", k, g, got)
		}
	}
	for k, g := range ast.ImportMap {
		if got := ast.ReverseMap[g]; got != k {
			t.Fatalf("import %q glyph %q reverse=%q", k, g, got)
		}
	}
}

func TestDefaultSymbols(t *testing.T) {
	var s ast.DefaultSymbols
	if g, ok := s.Encode("func"); !ok || g != "æ" {
		t.Fatalf("Encode(func)=%q %v", g, ok)
	}
	if g, ok := s.EncodePhrase("fmt.Println"); !ok || g != "¥" {
		t.Fatalf("EncodePhrase=%q %v", g, ok)
	}
	if g, ok := s.EncodeImport(`"net/http"`); !ok || g != "¯" {
		t.Fatalf("EncodeImport=%q %v", g, ok)
	}
}

func TestFixedGlyphsOneToken(t *testing.T) {
	cl, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		t.Fatal(err)
	}
	o2, err := tiktoken.GetEncoding("o200k_base")
	if err != nil {
		t.Fatal(err)
	}
	check := func(label, g string) {
		t.Helper()
		a, b := len(cl.Encode(g, nil, nil)), len(o2.Encode(g, nil, nil))
		if a != 1 || b != 1 {
			t.Fatalf("%s %q tokens cl=%d o2=%d", label, g, a, b)
		}
	}
	for k, g := range ast.KeywordMap {
		check("kw:"+k, g)
	}
	for k, g := range ast.PhraseMap {
		check("ph:"+k, g)
	}
	for k, g := range ast.ImportMap {
		check("imp:"+k, g)
	}
	for _, g := range ast.StringPoolGlyphs {
		check("pool:"+g, g)
	}
}

func TestAlphaKindsNonEmpty(t *testing.T) {
	if len(ast.AlphaKinds) == 0 || len(ast.OutOfAlphaKinds) == 0 {
		t.Fatal("kind lists empty")
	}
}
