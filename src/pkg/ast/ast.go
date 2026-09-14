// Package ast defines VeGo symbol mapping and grammar surface (v0.1β+ hard pass).
package ast

import "fmt"

// KeywordMap maps Go keywords to 1-token glyphs (cl100k + o200k).
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
	"fmt.Sprintf":         "â",
	"fmt.Fprintf":         "ã",
	"fmt.Errorf":          "ä",
	"http.HandleFunc":     "À",
	"http.ListenAndServe": "Á",
	"http.ResponseWriter": "Â",
	"http.Request":        "Ã",
	"http.StatusOK":       "å",
	"strconv.Atoi":        "Ä",
	"strconv.Itoa":        "è",
	"os.Args":             "Í",
	"os.Exit":             "ê",
	"log.Fatal":           "ë",
	"log.Println":         "ì",
	"errors.New":          "ï",
	"json.Marshal":        "ñ",
	"json.Unmarshal":      "ò",
	"strings.Contains":    "ó",
	"strings.Join":        "ô",
	"strings.Split":       "õ",
	"time.Now":            "ö",
	"context.Background":  "ù",
	"bytes.Buffer":        "ú",
}

// ImportMap maps quoted import-path literals to one glyph.
var ImportMap = map[string]string{
	`"fmt"`:           "©",
	`"os"`:            "ª",
	`"strconv"`:       "«",
	`"net/http"`:      "¯",
	`"net"`:           "Ç",
	`"io"`:            "É",
	`"time"`:          "Î",
	`"strings"`:       "Ð",
	`"bytes"`:         "Ñ",
	`"context"`:       "Ó",
	`"encoding/json"`: "Ö",
	`"log"`:           "Ú",
	`"errors"`:        "Ü",
	`"sync"`:          "à",
	`"path/filepath"`: "á",
}

// LitMap maps common multi-token quoted literals to one glyph (codec dictionary).
var LitMap = map[string]string{
	`"/"`:  "¹",
	`"--"`: "º",
	`"\n"`: "·",
	`" "`:  "´",
}

// CompositeMap maps keyword/ident combos to one glyph.
var CompositeMap = map[string]string{
	"package main": "¾",
	"func main":    "½",
	"else if":      "¼",
}

// StringPoolGlyphs are free 1-token glyphs for repeated-literal Σ tables.
var StringPoolGlyphs = []string{
	"³", "²", "û", "ü",
}

// ReverseMap is glyph → expansion text.
var ReverseMap map[string]string

func init() {
	ReverseMap = make(map[string]string)
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
	for k, v := range LitMap {
		add(k, v)
	}
	for k, v := range CompositeMap {
		add(k, v)
	}
	for _, g := range StringPoolGlyphs {
		if _, ok := ReverseMap[g]; ok {
			panic(fmt.Sprintf("string pool glyph %q collides", g))
		}
	}
}

// SymbolMap maps Go constructs to glyphs.
type SymbolMap interface {
	Encode(keyword string) (symbol string, ok bool)
	Decode(symbol string) (text string, ok bool)
	EncodePhrase(phrase string) (symbol string, ok bool)
	EncodeImport(quoted string) (symbol string, ok bool)
	EncodeLit(quoted string) (symbol string, ok bool)
}

// DefaultSymbols is the Beta dictionary.
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
func (DefaultSymbols) EncodeLit(quoted string) (string, bool) {
	s, ok := LitMap[quoted]
	return s, ok
}

// AlphaKinds lists construct kinds covered by the Alpha/Beta grammar subset.
var AlphaKinds = []string{
	"package", "import", "func", "params", "results", "basic_type", "named_type",
	"call", "selector", "literal", "short_decl", "if", "else", "for", "range",
	"return", "struct_type",
}

// OutOfAlphaKinds are explicitly unsupported.
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
