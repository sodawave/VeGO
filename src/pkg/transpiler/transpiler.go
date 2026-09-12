// Package transpiler provides the Go <-> VeGo bi-directional transform.
package transpiler

import (
	"bytes"
	"fmt"
	"go/format"
	"go/scanner"
	"go/token"
	"sort"
	"unicode"
	"unicode/utf8"

	"github.com/sodawave/VeGO/src/pkg/ast"
)

// Direction selects transform orientation.
type Direction int

const (
	// ToGo expands .vego into standard Go source.
	ToGo Direction = iota
	// ToVeGo compresses Go source into .vego.
	ToVeGo
)

// Transformer is a lossless, deterministic Go <-> VeGo converter.
type Transformer interface {
	Transform(src []byte, dir Direction) ([]byte, error)
}

// Engine performs token-level keyword/phrase glyph substitution using go/scanner.
type Engine struct {
	Symbols ast.SymbolMap
	glyphs  []string // longest-first glyph list
}

type phrased interface {
	EncodePhrase(phrase string) (symbol string, ok bool)
}

// New returns an Engine with Beta default symbols when Symbols is nil.
func New(symbols ast.SymbolMap) *Engine {
	if symbols == nil {
		symbols = ast.DefaultSymbols{}
	}
	glyphs := make([]string, 0, len(ast.ReverseMap))
	for g := range ast.ReverseMap {
		glyphs = append(glyphs, g)
	}
	sort.Slice(glyphs, func(i, j int) bool { return len(glyphs[i]) > len(glyphs[j]) })
	return &Engine{Symbols: symbols, glyphs: glyphs}
}

// Transform converts between Go and .vego.
func (e *Engine) Transform(src []byte, dir Direction) ([]byte, error) {
	switch dir {
	case ToVeGo:
		return e.toVeGo(src)
	case ToGo:
		return e.toGo(src)
	default:
		return nil, fmt.Errorf("unknown direction %d", dir)
	}
}

type scanned struct {
	tok token.Token
	lit string
}

func (e *Engine) toVeGo(src []byte) ([]byte, error) {
	if _, err := format.Source(src); err != nil {
		return nil, fmt.Errorf("vego: input is not valid Go: %w", err)
	}
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	var sc scanner.Scanner
	sc.Init(file, src, nil, 0)

	var toks []scanned
	for {
		_, tok, lit := sc.Scan()
		if tok == token.EOF {
			break
		}
		text := lit
		if text == "" {
			text = tok.String()
		}
		toks = append(toks, scanned{tok: tok, lit: text})
	}

	var out bytes.Buffer
	var prevKind kind
	for i := 0; i < len(toks); i++ {
		// Aggressive phrase fold: IDENT "." IDENT → one glyph when mapped.
		if toks[i].tok == token.IDENT && i+2 < len(toks) &&
			toks[i+1].tok == token.PERIOD && toks[i+2].tok == token.IDENT {
			phrase := toks[i].lit + "." + toks[i+2].lit
			if sym, ok := e.encodePhrase(phrase); ok {
				if needsSep(prevKind, kindGlyph) {
					out.WriteByte(' ')
				}
				out.WriteString(sym)
				prevKind = kindGlyph
				i += 2
				continue
			}
		}

		tok := toks[i].tok
		text := toks[i].lit

		// Emit ';' for ASI newlines and explicit semicolons → single-line IR.
		if tok == token.SEMICOLON {
			out.WriteByte(';')
			prevKind = kindOther
			continue
		}

		k := kindOther
		write := text
		if tok.IsKeyword() {
			if sym, ok := e.Symbols.Encode(text); ok {
				write = sym
				k = kindGlyph
			} else {
				k = kindIdent
			}
		} else {
			switch tok {
			case token.IDENT:
				k = kindIdent
			case token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING:
				k = kindLit
			default:
				write = tok.String()
				k = kindOther
			}
		}
		if needsSep(prevKind, k) {
			out.WriteByte(' ')
		}
		out.WriteString(write)
		prevKind = k
	}
	return out.Bytes(), nil
}

func (e *Engine) encodePhrase(phrase string) (string, bool) {
	if p, ok := e.Symbols.(phrased); ok {
		return p.EncodePhrase(phrase)
	}
	return ast.PhraseMap[phrase], ast.PhraseMap[phrase] != ""
}

type kind int

const (
	kindOther kind = iota
	kindGlyph
	kindIdent
	kindLit
)

func needsSep(prev, next kind) bool {
	if prev == kindOther || next == kindOther {
		if prev == kindLit && next == kindLit {
			return true
		}
		return false
	}
	return true
}

func (e *Engine) toGo(src []byte) ([]byte, error) {
	replaced, err := e.replaceGlyphs(src)
	if err != nil {
		return nil, err
	}
	formatted, err := format.Source(replaced)
	if err != nil {
		return nil, fmt.Errorf("vego: expanded source is not valid Go: %w\n---\n%s", err, replaced)
	}
	return formatted, nil
}

func (e *Engine) replaceGlyphs(src []byte) ([]byte, error) {
	var out bytes.Buffer
	i := 0
	for i < len(src) {
		if g, kw, n := e.matchGlyph(src[i:]); n > 0 {
			if identEnd(out.Bytes()) {
				out.WriteByte(' ')
			}
			out.WriteString(kw)
			i += n
			if i < len(src) && identStart(src[i:]) {
				out.WriteByte(' ')
			}
			_ = g
			continue
		}
		r, size := utf8.DecodeRune(src[i:])
		if r == utf8.RuneError && size == 1 {
			return nil, fmt.Errorf("vego: invalid UTF-8 at byte %d", i)
		}
		out.Write(src[i : i+size])
		i += size
	}
	return out.Bytes(), nil
}

func (e *Engine) matchGlyph(rest []byte) (glyph, keyword string, n int) {
	for _, g := range e.glyphs {
		gb := []byte(g)
		if len(rest) < len(gb) || !bytes.Equal(rest[:len(gb)], gb) {
			continue
		}
		kw, ok := e.Symbols.Decode(g)
		if !ok {
			continue
		}
		return g, kw, len(gb)
	}
	return "", "", 0
}

func identStart(rest []byte) bool {
	r, _ := utf8.DecodeRune(rest)
	return unicode.IsLetter(r) || r == '_'
}

func identEnd(out []byte) bool {
	if len(out) == 0 {
		return false
	}
	r, _ := utf8.DecodeLastRune(out)
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// Stub retains the prior name; delegates to Engine.
type Stub struct {
	Symbols ast.SymbolMap
}

func (s *Stub) Transform(src []byte, dir Direction) ([]byte, error) {
	return New(s.Symbols).Transform(src, dir)
}
