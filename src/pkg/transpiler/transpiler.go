// Package transpiler provides the Go <-> VeGo bi-directional transform stubs.
package transpiler

import "github.com/sodawave/VeGO/src/pkg/ast"

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

// Stub is a no-op Transformer until the AST mapping is implemented.
type Stub struct {
	Symbols ast.SymbolMap
}

// Transform returns src unchanged (clean-slate stub).
func (s *Stub) Transform(src []byte, _ Direction) ([]byte, error) {
	out := make([]byte, len(src))
	copy(out, src)
	return out, nil
}
