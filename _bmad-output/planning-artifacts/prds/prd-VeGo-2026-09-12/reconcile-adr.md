# Reconcile: ADR → PRD + addendum

**Input:** `docs/ADR.md`  
**Against:** `prd.md`, `addendum.md` (same folder)

## Intentionally aligned (not gaps)

| ADR claim | PRD / addendum stance |
| --- | --- |
| Transpile layer; never fork `cmd/compile` | NFR-1, addendum Mechanism, Scope Out |
| Stack: Participle / Cobra / mcp-go / tiktoken-go | addendum Mechanism; tiktoken mid-term only |
| MCP bridge for agents + CLI for humans/hosts | FR-4, FR-5 |
| Expand-to-Go for human review | FR-6, UJ-2 |
| ≥60% / &lt;40% BPE as MVP success | **Demoted correctly** — NFR-2, Mid-term metrics, addendum Rejected (forge kill) |

## Capability gaps (ADR → PRD timing / FR coverage)

Focus: CLI / MCP / git / diff / related surfaces. BPE % Alpha gate is **not** listed — demotion is correct.

1. **Git `diff` / textconv driver** — ADR Phase 4 + MVP criterion #3 (`*.vego` + `diff.vego.textconv` → readable Go). PRD Scope **Out of Alpha**; UJ-2 only says “textconv when available.” **Need:** explicit Mid-term FR (or confirm Alpha stays `vego fmt`-only for audit).

2. **MCP `search_symbol`** — ADR Fase 5 third tool (symbol lookup over repo AST). PRD FR-5 Alpha trio = read compressed / read expanded / structural patch only. **Need:** classify as Mid-term MCP FR or drop.

3. **CLI `vego tokens`** — ADR Phase 4 / Sprint Day 4–5 first-class command. PRD has NFR-3 mid-term measurability but no FR for a `tokens` (or equivalent) CLI surface. **Need:** Mid-term FR tying measurement to CLI/MCP, not only a library capability.

4. **Bidirectional CLI `fmt` (.go → `.vego`)** — ADR `fmt` is bi-directional (Go→VeGo serialize and VeGo→Go expand). PRD FR-6 is expand-to-Go only; Go→`.vego` sits in FR-3 round-trip tests / FR-1 host emit, not as a documented CLI reverse path. **Need:** Alpha FR for `fmt --to-vego` (or equivalent) if hosts/devs use CLI for IR emit, else mark Mid-term / host-only.

5. **Identifier minification / extreme minify in Alpha IR** — ADR core strategy (strip whitespace/comments + rename idents to short forms + dimensional `«»`). PRD “compact IR” does not say whether Alpha grammar includes identifier shortening vs keyword/structure compaction only. **Need:** Alpha scope line (in / out) so Transpile fidelity fixtures and Mid-term BPE work don’t fight.

## Qualitative / non-blocking ADR leftovers

- ADR “MCP as sole authorized agent write path” vs PRD dual CLI+MCP host surface — product choice; no FR conflict if CLI remains orchestration.
- ADR illustrative Unicode keyword map — research samples; Out of Alpha “guaranteed 1-token maps” already covers deferral.
- ADR `cmd/` vs repo `src/cmd` layout — superseded by spine/AD-4; addendum already corrects.
