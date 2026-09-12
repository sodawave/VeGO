# VeGo — forged idea

## Decisions (locked)

- VeGo = compact/minified textual IR over `go/ast`, reversible 1:1 (like JS minify for Go).
- LLM emits `.vego` (not Go). Go-only emit kills the layer’s advantage.
- Pipeline: NL → CLI/IDE host (vibe) → LLM → `.vego` → `vego` transpile → `go build`/`run` → behavior = original Go.
- Layer sits behind CLI/MCP/IDE; humans do not author `.vego` casually; audit via expand-to-Go.
- No fork of `cmd/compile`.
- Product root: `src/` (`pkg/ast`, `pkg/transpiler`, `pkg/mcp`, `cmd/vego`).
- **Alpha:** lossless compact + round-trip. **Mid-term:** measured BPE % as versioned goal (not alpha gate).

## Rejected

- ≥60% token reduction as v1/alpha acceptance criterion.
- LLM emits Go and VeGo only minifies on save as the primary path.
- VeGo as a human-readable language or compiler fork.

## Why it holds

Without `.vego` as the emitted source, the minify layer has no product surface. Execution fidelity comes from transpile-to-Go, not a new runtime.
