// Package ast defines VeGo symbol mapping and Alpha grammar surface.
package ast

import "fmt"

// KeywordMap maps Go keywords/operators (token text) to Alpha Unicode glyphs.
// Glyphs are chosen for compactness; BPE 1-token validation is mid-term (AD-5).
var KeywordMap = map[string]string{
	"package":   "ð",
	"import":    "î",
	"func":      "ƒ",
	"return":    "®",
	"if":        "¿",
	"else":      "¬",
	"for":       "∀",
	"range":     "∈",
	"var":       "∂",
	"const":     "ç",
	"type":      "†",
	"struct":    "§",
	"map":       "µ",
	"chan":      "⊳",
	"go":        "⇒",
	"defer":     "↺",
	"select":    "⇆",
	"case":      "◈",
	"default":   "⊘",
	"break":     "⎋",
	"continue":  "↻",
	"switch":    "⇌",
	"interface": "ℑ",}

// ReverseMap is keyword lookup from glyph.
var ReverseMap map[string]string

func init() {
	ReverseMap = make(map[string]string, len(KeywordMap))
	for k, v := range KeywordMap {
		if prev, ok := ReverseMap[v]; ok {
			panic(fmt.Sprintf("glyph collision: %q maps to both %q and %q", v, prev, k))
		}
		ReverseMap[v] = k
	}
}

// SymbolMap maps Go keywords/constructs to glyphs.
type SymbolMap interface {
	Encode(keyword string) (symbol string, ok bool)
	Decode(symbol string) (keyword string, ok bool)
}

// DefaultSymbols is the Alpha keyword dictionary.
type DefaultSymbols struct{}

func (DefaultSymbols) Encode(keyword string) (string, bool) {
	s, ok := KeywordMap[keyword]
	return s, ok
}

func (DefaultSymbols) Decode(symbol string) (string, bool) {
	k, ok := ReverseMap[symbol]
	return k, ok
}

// AlphaKinds lists construct kinds covered by the Alpha grammar subset.
var AlphaKinds = []string{
	"package", "import", "func", "params", "results", "basic_type", "named_type",
	"call", "selector", "literal", "short_decl", "if", "else", "for", "range",
	"return", "struct_type",
}

// OutOfAlphaKinds are explicitly unsupported in Alpha.
var OutOfAlphaKinds = []string{
	"generics", "goroutine", "chan", "select", "reflect", "cgo", "unsafe", "ident_minify",
}

// Grammar is the VeGo CFG over the serialized AST surface.
type Grammar interface {
	Validate(src []byte) error
}

// Node is a VeGo AST node kind tag.
type Node interface {
	Kind() string
}

// KindTag is a simple Node for Alpha surface documentation/tests.
type KindTag string

func (k KindTag) Kind() string { return string(k) }
