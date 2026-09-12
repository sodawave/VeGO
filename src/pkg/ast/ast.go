// Package ast defines VeGo symbol mapping and grammar surface (v0.1β).
package ast

import "fmt"

// KeywordMap maps Go keywords to Latin-1 / letterlike glyphs that are exactly
// one BPE token in both cl100k_base and o200k_base (validated in tests).
// Aggressive rule: never emit a glyph that costs more tokens than the keyword.
var KeywordMap = map[string]string{
	"package":   "ð",
	"import":    "î",
	"func":      "æ",
	"return":    "®",
	"if":        "¿",
	"else":      "¬",
	"for":       "ø",
	"range":     "×",
	"var":       "ý",
	"const":     "ç",
	"type":      "†",
	"struct":    "§",
	"map":       "µ",
	"chan":      "¢",
	"go":        "»",
	"defer":     "ß",
	"select":    "±",
	"case":      "°",
	"default":   "¤",
	"break":     "¡",
	"continue":  "¨",
	"switch":    "¦",
	"interface": "¶",
}

// PhraseMap maps high-frequency multi-token selectors to a single glyph.
// Applied as IDENT "." IDENT so `fmt.Println` (2 tokens) → one glyph (1 token).
// Only phrases that save tokens under tiktoken belong here.
var PhraseMap = map[string]string{
	"fmt.Println":           "¥",
	"fmt.Fprintln":          "£",
	"http.HandleFunc":       "À",
	"http.ListenAndServe":   "Á",
	"http.ResponseWriter":   "Â",
	"http.Request":          "Ã",
	"strconv.Atoi":          "Ä",
}

// ReverseMap is glyph → keyword or phrase (e.g. "fmt.Println").
var ReverseMap map[string]string

func init() {
	ReverseMap = make(map[string]string, len(KeywordMap)+len(PhraseMap))
	for k, v := range KeywordMap {
		if prev, ok := ReverseMap[v]; ok {
			panic(fmt.Sprintf("glyph collision: %q maps to both %q and %q", v, prev, k))
		}
		ReverseMap[v] = k
	}
	for k, v := range PhraseMap {
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
	EncodePhrase(phrase string) (symbol string, ok bool)
}

// DefaultSymbols is the Beta keyword + phrase dictionary.
type DefaultSymbols struct{}

func (DefaultSymbols) Encode(keyword string) (string, bool) {
	s, ok := KeywordMap[keyword]
	return s, ok
}

func (DefaultSymbols) Decode(symbol string) (string, bool) {
	k, ok := ReverseMap[symbol]
	return k, ok
}

func (DefaultSymbols) EncodePhrase(phrase string) (string, bool) {
	s, ok := PhraseMap[phrase]
	return s, ok
}

// AlphaKinds lists construct kinds covered by the Alpha/Beta grammar subset.
var AlphaKinds = []string{
	"package", "import", "func", "params", "results", "basic_type", "named_type",
	"call", "selector", "literal", "short_decl", "if", "else", "for", "range",
	"return", "struct_type",
}

// OutOfAlphaKinds are explicitly unsupported in Alpha/Beta grammar subset.
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

// KindTag is a simple Node for surface documentation/tests.
type KindTag string

func (k KindTag) Kind() string { return string(k) }
