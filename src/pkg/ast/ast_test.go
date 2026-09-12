package ast_test

import (
	"testing"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sodawave/VeGO/src/pkg/ast"
)

func TestKeywordMapBijective(t *testing.T) {
	for k, g := range ast.KeywordMap {
		got, ok := ast.ReverseMap[g]
		if !ok || got != k {
			t.Fatalf("glyph %q for %q reverse=%q ok=%v", g, k, got, ok)
		}
	}
	for k, g := range ast.PhraseMap {
		got, ok := ast.ReverseMap[g]
		if !ok || got != k {
			t.Fatalf("phrase glyph %q for %q reverse=%q ok=%v", g, k, got, ok)
		}
	}
}

func TestDefaultSymbols(t *testing.T) {
	var s ast.DefaultSymbols
	enc, ok := s.Encode("func")
	if !ok || enc != "æ" {
		t.Fatalf("Encode(func)=%q %v want æ", enc, ok)
	}
	dec, ok := s.Decode("æ")
	if !ok || dec != "func" {
		t.Fatalf("Decode(æ)=%q %v", dec, ok)
	}
	ph, ok := s.EncodePhrase("fmt.Println")
	if !ok || ph != "¥" {
		t.Fatalf("EncodePhrase(fmt.Println)=%q %v", ph, ok)
	}
}

func TestGlyphsOneTokenBothEncodings(t *testing.T) {
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
		a := len(cl.Encode(g, nil, nil))
		b := len(o2.Encode(g, nil, nil))
		if a != 1 || b != 1 {
			t.Fatalf("%s glyph %q tokens cl=%d o2=%d (want 1/1)", label, g, a, b)
		}
	}
	for k, g := range ast.KeywordMap {
		check("keyword:"+k, g)
	}
	for k, g := range ast.PhraseMap {
		check("phrase:"+k, g)
	}
}

func TestAlphaKindsNonEmpty(t *testing.T) {
	if len(ast.AlphaKinds) == 0 || len(ast.OutOfAlphaKinds) == 0 {
		t.Fatal("Alpha kind lists must be populated")
	}
}
