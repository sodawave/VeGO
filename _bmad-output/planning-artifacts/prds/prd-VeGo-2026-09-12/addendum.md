# VeGo PRD Addendum

Mechanism and architecture pointers — not product requirements. Spine is authoritative for invariants.

## Mechanism (from docs/ADR + spine)

- Source-to-source transpile over `go/ast`; never fork `cmd/compile` (AD-1).
- Planned libs: Participle v2.1.4 (CFG), Cobra v1.10.2 (CLI), mark3labs/mcp-go v1.0.0 (MCP), tiktoken-go mid-term only.
- Layout: `src/pkg/ast`, `src/pkg/transpiler`, `src/pkg/mcp`, `src/cmd/vego` (AD-4).

## Rejected alternatives (forge)

- Go-primary LLM emit with minify-on-save as the main path.
- ≥60% BPE as Alpha acceptance gate (demoted to mid-term versioned goal).
- Human-first `.vego` language / compiler-fork framing.

## Pointers

- Forge: `_bmad-output/forge/vego/forged-idea.md`
- Spine: `_bmad-output/planning-artifacts/architecture/architecture-VeGo-2026-09-12/ARCHITECTURE-SPINE.md`
- Research: `docs/OVERVIEW.md`, `docs/ADR.md` (treat 60–85% claims as mid-term hypothesis, not Alpha gate)
