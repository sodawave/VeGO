# VeGo PRD Addendum

Mechanism and architecture pointers — not product requirements. Spine is authoritative for invariants.

## Mechanism (from docs/ADR + spine)

- Source-to-source transpile over `go/ast`; never fork `cmd/compile` (AD-1).
- Planned libs: Participle v2.1.4 (CFG), Cobra v1.10.2 (CLI), mark3labs/mcp-go v1.0.0 (MCP), tiktoken-go mid-term only.
- Layout: `src/pkg/ast`, `src/pkg/transpiler`, `src/pkg/mcp`, `src/cmd/vego` (AD-4).
- Pipeline one-liner: NL → Host CLI → LLM emits `.vego` → transpile → `go build`/`run`.
- Analogy: JS minify for Go IR — compact storage/agent surface, not a new runtime.

## Rejected alternatives (forge)

- Go-primary LLM emit with minify-on-save as the main path.
- ≥60% BPE as Alpha acceptance gate (demoted to mid-term versioned goal).
- Human-first `.vego` language / compiler-fork framing.

## Deferred to mid-term / later (reconcile)

- BPE dictionary validation per tokenizer; `vego tokens`.
- MCP `search_symbol`; deep IDE Host integration beyond minimal MCP.
- Git `textconv` / diff driver; LLM calibration dataset / system prompt as productized assets.
- Aggressive identifier minification maps; OVERVIEW’s multi-strategy compression taxonomy as explicit product modes.

## Pointers

- Forge: `_bmad-output/forge/vego/forged-idea.md`
- Spine: `_bmad-output/planning-artifacts/architecture/architecture-VeGo-2026-09-12/ARCHITECTURE-SPINE.md`
- Research: `docs/OVERVIEW.md`, `docs/ADR.md` (treat 60–85% claims as mid-term hypothesis, not Alpha gate)
- Finalize reviews: `review-checklist.md`, `reconcile-*.md` in this folder
