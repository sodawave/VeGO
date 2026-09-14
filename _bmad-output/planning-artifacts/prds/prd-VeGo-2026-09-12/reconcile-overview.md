# Reconcile: OVERVIEW → PRD + addendum

**Input:** `docs/OVERVIEW.md`  
**Against:** `prd.md`, `addendum.md` (same folder)

## Intentionally aligned (not gaps)

| OVERVIEW claim | PRD / addendum stance |
| --- | --- |
| Go `go/ast` base; transpile layer; no compiler fork | Vision, NFR-1, addendum Mechanism |
| Bidirectional lossless expand/emit path | FR-2, FR-3 |
| CLI `vego` build/run/fmt + MCP bridge | FR-4, FR-5, FR-6 |
| Human-illegible IR; audit via expand-to-Go | UJ-2, FR-6, Non-Users |
| 60–85% BPE token reduction as viability claim | **Demoted correctly** — Vision, NFR-2/3, Mid-term metrics, addendum Rejected + Pointers (“treat 60–85% as mid-term hypothesis, not Alpha gate”) |

## Capability gaps (OVERVIEW → PRD coverage)

BPE 60–85% as Alpha gate is **not** listed — demotion is correct.

2. **LLM calibration dataset + system prompt** — OVERVIEW §Integración treats a Go↔`.vego` example corpus and a standard system prompt as required to teach LLM emission and cut syntactic hallucinations. PRD/addendum omit it (spine Defers post-CFG). **Need:** Mid-term FR / Open Question so emission training is tracked when CFG stabilizes.

3. **Three compression strategies unnamed in product scope** — OVERVIEW’s core design is Symbolic Mapping (“token hacking”), Extreme Native Minification (whitespace strip + short idents), and Dimensional Structure (positional/Unicode scope vs `{}`). PRD only says “minify-style / compact IR”; addendum Mechanism does not enumerate the trio. **Need:** Alpha vs Mid-term scope lines for each strategy so fixtures and BPE work stay coherent.

4. **Mandatory contextual human documentation as semantic bridge** — OVERVIEW makes per-branch/PR human docs the primary bridge for an illegible store (not optional). PRD UJ-2/FR-6 cover expand-to-Go audit only. **Need:** process/NFR or Out-of-Alpha note if docs-as-gate stays methodology-only (BMAD), else a Mid-term review requirement.

5. **Benchmark suite beyond BPE %** — OVERVIEW calls for API cost, transpile latency, and compile/run success vs Go/TS/Python (and other token-reduction techniques). PRD Success Metrics = Alpha fidelity + Mid-term BPE reduction only (NFR-3 measurability). **Need:** Mid-term NFR row(s) for latency/cost/success, or explicitly defer comparative benches.

## Qualitative / non-blocking OVERVIEW leftovers

- Go vs TypeScript decision table — research rationale; product already locked on `go/ast`.
- Git `vego-diff` / AST merge driver — OVERVIEW critical path; PRD Scope Out of Alpha (aligned with spine Deferred); covered in ADR reconcile.
- Rob Pike AST “redesign” caveat — historical note; no PRD action.
- Pilot use-case / cultural-documentation recommendations — strategic advice, not FR surface.
