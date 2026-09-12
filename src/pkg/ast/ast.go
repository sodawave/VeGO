// Package ast defines VeGo symbol mapping and CFG grammar stubs.
// Implementation maps bijectively onto go/ast nodes (see docs/ADR.md).
package ast

// SymbolMap maps Go keywords/constructs to BPE-validated Unicode glyphs.
type SymbolMap interface {
	Encode(keyword string) (symbol string, ok bool)
	Decode(symbol string) (keyword string, ok bool)
}

// Grammar is the VeGo CFG over the serialized AST surface.
type Grammar interface {
	Validate(src []byte) error
}

// Node is a placeholder for a VeGo AST node pending CFG adoption.
type Node interface {
	Kind() string
}
