// Package transpiler provides the Go <-> VeGo bi-directional transform.
package transpiler

import (
	"bytes"
	"fmt"
	"go/format"
	"go/scanner"
	"go/token"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/sodawave/VeGO/src/pkg/ast"
	"github.com/sodawave/VeGO/src/pkg/bpe"
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

// Engine performs token-level keyword/phrase/import/string glyph substitution.
type Engine struct {
	Symbols ast.SymbolMap
	glyphs  []string
}

type phrased interface {
	EncodePhrase(phrase string) (symbol string, ok bool)
}

type imported interface {
	EncodeImport(quoted string) (symbol string, ok bool)
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

	// Dynamic string table only for repeated multi-token literals.
	// Single-occurrence strings + preamble are a net token loss; ImportMap
	// stays preamble-free (codec dictionary) and still wins.
	freq := map[string]int{}
	for _, t := range toks {
		if t.tok == token.STRING {
			freq[t.lit]++
		}
	}
	strGlyph := map[string]string{}
	poolIdx := 0
	for lit, n := range freq {
		if n < 2 {
			continue
		}
		if _, ok := e.encodeImport(lit); ok {
			continue
		}
		tokn, err := bpe.Count([]byte(lit), bpe.DefaultEncoding)
		if err != nil || tokn <= 1 {
			continue
		}
		if poolIdx >= len(ast.StringPoolGlyphs) {
			break
		}
		strGlyph[lit] = ast.StringPoolGlyphs[poolIdx]
		poolIdx++
	}

	var out bytes.Buffer
	if len(strGlyph) > 0 {
		out.WriteRune('Σ')
		// Stable order for determinism.
		keys := make([]string, 0, len(strGlyph))
		for k := range strGlyph {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, lit := range keys {
			out.WriteString(strGlyph[lit])
			out.WriteString(lit) // includes quotes
		}
		out.WriteByte(';')
	}

	var prevKind kind
	for i := 0; i < len(toks); i++ {
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
			case token.STRING:
				if sym, ok := e.encodeImport(text); ok {
					write = sym
					k = kindGlyph
				} else if sym, ok := strGlyph[text]; ok {
					write = sym
					k = kindGlyph
				} else {
					k = kindLit
				}
			case token.INT, token.FLOAT, token.IMAG, token.CHAR:
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

func (e *Engine) encodeImport(quoted string) (string, bool) {
	if p, ok := e.Symbols.(imported); ok {
		return p.EncodeImport(quoted)
	}
	return ast.ImportMap[quoted], ast.ImportMap[quoted] != ""
}

type kind int

const (
	kindOther kind = iota
	kindGlyph
	kindIdent
	kindLit
)

func needsSep(prev, next kind) bool {
	// Densify: glyphs glue to idents/lits; expand re-inserts required spaces.
	if prev == kindGlyph && (next == kindIdent || next == kindLit) {
		return false
	}
	if prev == kindOther || next == kindOther {
		if prev == kindLit && next == kindLit {
			return true
		}
		return false
	}
	return true
}

func (e *Engine) toGo(src []byte) ([]byte, error) {
	body, local, err := splitStringTable(src)
	if err != nil {
		return nil, err
	}
	replaced, err := e.replaceGlyphs(body, local)
	if err != nil {
		return nil, err
	}
	formatted, err := format.Source(replaced)
	if err != nil {
		return nil, fmt.Errorf("vego: expanded source is not valid Go: %w\n---\n%s", err, replaced)
	}
	return formatted, nil
}

// splitStringTable parses optional Σglyph"str"glyph"str"; preamble.
func splitStringTable(src []byte) (body []byte, local map[string]string, err error) {
	local = map[string]string{}
	if !bytes.HasPrefix(src, []byte("Σ")) {
		return src, local, nil
	}
	i := len("Σ")
	for i < len(src) {
		if src[i] == ';' {
			return src[i+1:], local, nil
		}
		g, gsz := utf8.DecodeRune(src[i:])
		if g == utf8.RuneError && gsz == 1 {
			return nil, nil, fmt.Errorf("vego: invalid UTF-8 in string table")
		}
		glyph := string(g)
		i += gsz
		if i >= len(src) || src[i] != '"' {
			return nil, nil, fmt.Errorf("vego: string table entry %q missing quoted literal", glyph)
		}
		j := i + 1
		for j < len(src) {
			if src[j] == '\\' {
				j += 2
				continue
			}
			if src[j] == '"' {
				j++
				break
			}
			j++
		}
		if j > len(src) {
			return nil, nil, fmt.Errorf("vego: unterminated string in Σ table")
		}
		local[glyph] = string(src[i:j])
		i = j
	}
	return nil, nil, fmt.Errorf("vego: Σ string table not terminated with ';'")
}

func (e *Engine) replaceGlyphs(src []byte, local map[string]string) ([]byte, error) {
	var out bytes.Buffer
	i := 0
	for i < len(src) {
		if g, kw, n := e.matchGlyph(src[i:], local); n > 0 {
			if needsSpaceBeforeExpand(out.Bytes(), kw) {
				out.WriteByte(' ')
			}
			out.WriteString(kw)
			i += n
			if i < len(src) && needsSpaceAfterExpand(kw, src[i:]) {
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

func (e *Engine) matchGlyph(rest []byte, local map[string]string) (glyph, keyword string, n int) {
	// Prefer local string-table glyphs (also 1 rune).
	if r, size := utf8.DecodeRune(rest); size > 0 {
		g := string(r)
		if lit, ok := local[g]; ok {
			return g, lit, size
		}
	}
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

func needsSpaceAfterExpand(kw string, rest []byte) bool {
	if len(rest) == 0 {
		return false
	}
	// Phrases/imports that already include dots or quotes don't need trailing space before '('.
	if strings.ContainsAny(kw, ".\"") {
		// import path / phrase: space before ident only
		return identStart(rest)
	}
	r, _ := utf8.DecodeRune(rest)
	if r == '"' || r == '\'' || r == '`' {
		return true // return"x" invalid
	}
	return identStart(rest)
}

func needsSpaceBeforeExpand(out []byte, kw string) bool {
	if len(out) == 0 {
		return false
	}
	if !identEnd(out) {
		return false
	}
	// Avoid gluing ident + expanded keyword/phrase.
	if kw == "" {
		return false
	}
	r, _ := utf8.DecodeRune([]byte(kw))
	return unicode.IsLetter(r) || r == '_' || r == '"'
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
