// Package ast defines VeGo symbol mapping and grammar surface (v0.1β+).
package ast

import "fmt"

// KeywordMap maps Go keywords to Latin-1 / letterlike glyphs that are exactly
// one BPE token in both cl100k_base and o200k_base (validated in tests).
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
var PhraseMap = map[string]string{
	"fmt.Println":         "¥",
	"fmt.Fprintln":        "£",
	"http.HandleFunc":     "À",
	"http.ListenAndServe": "Á",
	"http.ResponseWriter": "Â",
	"http.Request":        "Ã",
	"strconv.Atoi":        "Ä",
	"os.Args":             "Í",
}

// ImportMap maps quoted import-path string literals to one glyph.
// `"net/http"` is 3 cl100k tokens; a Latin-1 glyph is 1.
var ImportMap = map[string]string{
	`"fmt"`:        "©",
	`"os"`:         "ª",
	`"strconv"`:    "«",
	`"net/http"`:   "¯",
	`"net"`:        "Ç",
	`"io"`:         "É",
	`"time"`:       "Î",
	`"strings"`:    "Ð",
	`"bytes"`:      "Ñ",
	`"context"`:    "Ó",
	`"encoding/json"`: "Ö",
	`"log"`:        "Ú",
	`"errors"`:     "Ü",
	`"sync"`:       "à",
	`"path/filepath"`: "á",
}

// StringPoolGlyphs are free 1-token glyphs reserved for per-file string tables (Σ…).
// Must not overlap KeywordMap / PhraseMap / ImportMap (checked in init).
var StringPoolGlyphs = []string{
	"³", "´", "·", "¹", "º", "¼", "½", "¾", "â", "ã", "ä", "å", "è", "ê", "ë",
	"ì", "ï", "ñ", "ò", "ô", "õ", "ö", "ù", "û",
}

// ReverseMap is glyph → keyword, phrase, or quoted import/string.
var ReverseMap map[string]string

func init() {
	ReverseMap = make(map[string]string, len(KeywordMap)+len(PhraseMap)+len(ImportMap))
	add := func(k, v string) {
		if prev, ok := ReverseMap[v]; ok {
			panic(fmt.Sprintf("glyph collision: %q maps to both %q and %q", v, prev, k))
		}
		ReverseMap[v] = k
	}
	for k, v := range KeywordMap {
		add(k, v)
	}
	for k, v := range PhraseMap {
		add(k, v)
	}
	for k, v := range ImportMap {
		add(k, v)
	}
	// Pool glyphs must not collide with fixed maps.
	for _, g := range StringPoolGlyphs {
		if _, ok := ReverseMap[g]; ok {
			panic(fmt.Sprintf("string pool glyph %q collides with fixed map", g))
		}
	}
}

// SymbolMap maps Go keywords/constructs to glyphs.
type SymbolMap interface {
	Encode(keyword string) (symbol string, ok bool)
	Decode(symbol string) (keyword string, ok bool)
	EncodePhrase(phrase string) (symbol string, ok bool)
	EncodeImport(quoted string) (symbol string, ok bool)
}

// DefaultSymbols is the Beta keyword + phrase + import dictionary.
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

func (DefaultSymbols) EncodeImport(quoted string) (string, bool) {
	s, ok := ImportMap[quoted]
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
