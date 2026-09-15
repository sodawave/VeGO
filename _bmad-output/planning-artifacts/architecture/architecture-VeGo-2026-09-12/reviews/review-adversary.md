# Adversarial review — Architecture Spine VeGo (2026-09-12)

**Lens:** Construct two units one level down that each obey every AD to the letter yet still build incompatibly (clashing shared-data shapes, dual ownership, conflicting mutation paths). Each pair is a hole to close with a new or tightened AD.

**Inputs read:** `ARCHITECTURE-SPINE.md`, `_bmad-output/forge/vego/forged-idea.md` only.

**Verdict:** **fail** — multiple independent pairs can ship while citing AD-1…AD-6 and still disagree on IR shape, symbol ownership, write semantics, and the go/ast boundary.

---

## Method

Altitude is **initiative** → one level down ≈ packages / host surfaces (`pkg/ast`, `pkg/transpiler`, `pkg/mcp`, `cmd/vego`, plus any mid-term `pkg/bpe` named in the capability map). For each hole: Unit A and Unit B cite the same Rules; the clash is on something the spine never pins.

---

## Pair 1 — Clashing shared-data shape: on-disk `.vego` vs in-memory VeGo AST

### Units

| | Unit A — `pkg/ast` team | Unit B — `pkg/transpiler` team |
| --- | --- | --- |
| Claimed obedience | AD-2 (bijective w/ `go/ast`), AD-4 (`pkg/ast` = symbols + CFG) | AD-1, AD-2 (Go↔VeGo), AD-4 (`pkg/transpiler` = transform) |
| Choice | Canonical product is an **in-memory VeGo CFG** (node structs + symbol IDs). Disk `.vego` is “whatever the printer dumps” — e.g. Unicode-glyph stream with no version header, no package path encoding rule. | Canonical product is a **textual `.vego` file schema**: ASCII keyword abbreviations + optional BOM + `#vego-1` header; scanner owns lexemes; AST package is only a symbol table lookup. |

### Why every AD still holds

- AD-1: both expand to toolchain Go.
- AD-2: each can round-trip *its* IR ↔ `go/ast` for the subset it implements.
- AD-4: A owns “symbols + CFG”; B owns “Go↔VeGo” — neither Rule names **who owns the byte-level `.vego` grammar / versioning / encoding**.
- Consistency row “`.vego` = product source” names the *role*, not the *shape*.

### Incompatibility

A’s printer and B’s scanner disagree on tokens, headers, and identifier escaping. Agents emit files B cannot parse; A’s CFG never appears on disk in a form B accepts. Round-trip proofs in each package pass in isolation and fail across the package boundary.

### Hole → close

**New AD (IR text contract):** One owner for the on-disk `.vego` lexical/grammar contract (version header, encoding, symbol rendering). Prefer: `pkg/ast` owns alphabet + CFG + printer/parser for `.vego` text; `pkg/transpiler` only maps that IR ↔ `go/ast` / Go source. Or the inverse — but **one** owner, named, with a stable versioned schema.

---

## Pair 2 — Two owners of one entity: keyword ↔ Unicode / symbol map

### Units

| | Unit A — `pkg/ast` | Unit B — mid-term `pkg/bpe` (named in Capability → Architecture Map) / or transpiler “lex helpers” |
| --- | --- | --- |
| Claimed obedience | AD-4 (`pkg/ast` = symbols), AD-2 | AD-5 (BPE is mid-term measure, not Alpha gate); Deferred “Unicode↔keyword dictionary… mid-term” |
| Choice | Symbol map is the **grammar alphabet** — authoritative for parse/print; IDs are stable product API. | Maintains a **second** keyword↔glyph table for tiktoken experiments; ships “optimized” maps into agent prompts / emit hints without going through `pkg/ast`. |

### Why every AD still holds

- AD-4 binds symbols to `pkg/ast` but does not say the map is the **sole** product dictionary or that measurement packages must import it.
- AD-5 explicitly separates BPE work from Alpha — Unit B treats that as license to own a parallel dictionary “for metrics only,” then agents consume it for emission (AD-3 still “emits `.vego`”).
- Deferred item *names* the dictionary work but does not assign a single owner or forbid a second table.

### Incompatibility

Emit path (AD-3) and transpile path (AD-2) use different glyph sets. Bijective claims hold only against each unit’s private map. `pkg/bpe` appears in the capability map but not in Structural Seed or AD-4 — a silent second home for the same entity.

### Hole → close

**Tighten AD-4 (or new AD-7 — Sole symbol map):** The keyword↔compact-token dictionary has exactly one owner (`pkg/ast`). Measurement (`pkg/bpe` / tiktoken) **reads** that map; it must not define or publish a competing alphabet. Agent prompts / emission contracts reference the same map version as the parser.

---

## Pair 3 — Conflicting state-mutation paths: who may write `.vego` / `.go` and where

### Units

| | Unit A — `cmd/vego` | Unit B — `pkg/mcp` |
| --- | --- | --- |
| Claimed obedience | AD-3, AD-4 (CLI orchestration), AD-6 (`vego fmt` expand-to-Go) | AD-3, AD-4 (protocol surface only), AD-6 (humans review Go) |
| Choice | CLI is the **only writer**: `transpile` writes sibling `.go`; `fmt` overwrites audit Go; agents shell out to CLI. Fail closed = non-zero exit, no partial file. | MCP tools **write repo files** directly (`write_vego`, `expand`, `build`); returns Go in tool payloads *and* writes `*_gen.go` under a different layout. Errors returned as JSON-RPC success with `isError` content; partial files left on disk. |

### Why every AD still holds

- AD-4: both call `pkg/transpiler`; neither “owns grammar.” MCP is “protocol surface only” — Unit B argues write tools are protocol, not grammar.
- AD-3: both require agent product source to be `.vego` (Unit B may also return expanded Go for “audit,” citing AD-6).
- Consistency: “`.go` = ephemeral **or** audit artifact unless explicitly checked in” — each unit picks a different branch of the OR.
- “Fail closed” is a convention, not an AD; error transport (exit code vs MCP content) is unbound.

### Incompatibility

Two writers, two path conventions, two error/partial-write policies. Host integrations (Cursor vs CLI scripts) corrupt each other’s artifacts; CI cannot rely on a single mutation story. Dual owners of “repo product state.”

### Hole → close

**New AD (mutation authority):** Define (1) which surfaces may create/update/delete `.vego` and expanded `.go`, (2) canonical path layout for expand artifacts (ephemeral-only vs sidecar vs explicit check-in), (3) fail-closed write semantics (no partial files; error shape). Likely: `pkg/transpiler` is pure (no FS policy); **one** orchestration policy shared by CLI and MCP (same writer module), MCP must not invent a second layout.

---

## Pair 4 — Dual ownership of the go/ast boundary type

### Units

| | Unit A — `pkg/transpiler` | Unit B — `pkg/ast` |
| --- | --- | --- |
| Claimed obedience | AD-4 Rule: transpiler = Go↔VeGo; diagram `xf → go_ast_parser_printer` | AD-4 Rule: `pkg/ast` = symbols + CFG / VeGo AST; AD-2 binds both packages |
| Choice | Public API is `func ParseVego([]byte) (*goast.File, error)` and `func PrintGo(*goast.File) ([]byte, error)` — stdlib `go/ast` is the **only** shared IR; `pkg/ast` is a data package of maps. | Public API is `type File struct {…}` (VeGo CFG); transpiler must accept/return **VeGo AST**; conversion to `go/ast` is a helper **inside** `pkg/ast`. |

### Why every AD still holds

- AD-2 requires semantic equivalence with `go/ast`, not that stdlib nodes are the package boundary type.
- AD-4 lists responsibilities in slogans without naming the **shared type** crossing `pkg/ast` ↔ `pkg/transpiler` ↔ hosts.
- Mermaid shows both `xf → ast` and `xf → stdlib`, which licenses either reading.

### Incompatibility

MCP/CLI authors import different root types; tests assert different round-trip surfaces; circular or duplicated convert layers appear. Two owners of “the” IR entity at the API seam.

### Hole → close

**Tighten AD-2/AD-4:** Name the canonical in-process IR type and which package exports it. Example Rule: hosts and MCP speak only to `pkg/transpiler`; `pkg/transpiler` is the sole package that imports `go/ast`/`go/parser`/`go/printer`; `pkg/ast` exports VeGo CFG + symbol map only — **or** the reverse, but pick one and forbid the other package from re-exporting a competing root type.

---

## Pair 5 — “Supported grammar subset” unbounded → fidelity gates diverge

### Units

| | Unit A — Alpha implementer (`pkg/ast` + transpiler) | Unit B — TEA / fidelity gate author (AD-5) |
| --- | --- | --- |
| Claimed obedience | AD-2 (“for the supported grammar subset”), Deferred “full grammar… wait until round-trip core” | AD-5 (alpha = lossless compact IR + transpile + `go build`/`run`) |
| Choice | Alpha subset = package clause + funcs + basic stmts; rejects generics/`go`/`select` at parse. | Alpha gate corpus includes concurrency and generics because “lossless + `go build`” with no subset annex. |

### Why every AD still holds

- AD-2 explicitly scopes bijection to an undefined “supported grammar subset.”
- Deferred acknowledges incomplete coverage but does not bind a **subset contract owner** or artifact.
- AD-5’s alpha Rule does not say “subset S as published in X.”

### Incompatibility

A ships “alpha done”; B’s gates fail forever (or B lowers gates and A expands ad hoc). Not a data-shape clash alone — it is an unbound shared contract that lets two units define product readiness incompatibly.

### Hole → close

**New AD or annex under AD-2/AD-5:** Alpha grammar subset is a versioned, single-owned artifact (e.g. `pkg/ast` + documented subset list); fidelity gates may only require constructs inside that list; expansion of the subset is an explicit spine/spec change, not a silent package choice.

---

## Pair 6 (secondary) — Agent emission channel vs human expand channel

### Units

| | Unit A — MCP “generate” tool | Unit B — CLI human-bridge path |
| --- | --- | --- |
| Claimed obedience | AD-3 (agents emit `.vego`) | AD-6 (humans review via expand-to-Go) |
| Choice | Tool accepts NL and **writes `.vego`**, optionally also writes Go “for convenience.” | `vego fmt` is expand-only; any Go in the repo is treated as human-authored product (contradicting AD-3 if agents also wrote Go). |

### Why it matters

Forged idea Rejected: “LLM emits Go and VeGo only minifies on save.” AD-3 Prevents “Go-primary generation,” but neither AD forbids **dual emit** (`.vego` + `.go` together) nor says which file wins when both exist and drift. Unit A can claim AD-3 (`.vego` was generated) while shipping Go-primary drift.

### Hole → close

**Tighten AD-3:** Product source of truth in the repo is `.vego` only; expanded Go must be non-authoritative (derived, marked, or out-of-tree). Dual-write that treats Go as editable product source is forbidden. Conflict resolution when both present: `.vego` wins; Go is regenerated.

---

## Cross-check against forged idea

| Locked forge decision | Covered by spine? | Adversarial gap |
| --- | --- | --- |
| Compact IR, reversible 1:1 | AD-2 | IR **text** shape / owner unset (Pair 1, 4) |
| LLM emits `.vego` not Go | AD-3 | Dual-write / SoT when both exist (Pair 6) |
| Pipeline to `go build`/`run` | AD-1 | Mutation/write policy across CLI vs MCP (Pair 3) |
| Humans audit via expand | AD-6 | Expand artifact layout unbound (Pair 3) |
| No `cmd/compile` fork | AD-1 | OK — hard to violate while obeying Rule |
| `src/` package layout | AD-4 | Symbol map / `pkg/bpe` / go/ast seam under-specified (Pairs 2, 4) |
| Alpha = lossless; BPE mid-term | AD-5 | Subset + dictionary ownership still open (Pairs 2, 5) |

---

## Recommended AD closures (priority)

1. **P0 — Sole IR text + symbol map owner** (Pairs 1–2): close shared-data and dual-map holes.
2. **P0 — Canonical in-process IR / `go/ast` import boundary** (Pair 4): close API dual ownership.
3. **P0 — Repo mutation & artifact layout policy for CLI+MCP** (Pair 3): close conflicting writers.
4. **P1 — Alpha grammar subset contract** (Pair 5): bind AD-2/AD-5 to a named subset artifact.
5. **P1 — Source-of-truth / no dual-authoritative Go** (Pair 6): tighten AD-3.

Until those land, the spine is **not** a sufficient consistency contract for independent package/host teams one level down.

---

## Summary for gate

| ID | Severity | Class | Finding |
| --- | --- | --- | --- |
| ADV-1 | critical | shared-data shape | `.vego` on-disk schema owner unset between `pkg/ast` and `pkg/transpiler` |
| ADV-2 | critical | dual ownership | Symbol/Unicode map can be owned by `pkg/ast` and `pkg/bpe`/emit path |
| ADV-3 | critical | mutation path | CLI and MCP can write different layouts with different fail-closed behavior |
| ADV-4 | high | dual ownership | Stdlib `go/ast` vs VeGo CFG as public boundary type both legal |
| ADV-5 | high | unbound contract | “Supported grammar subset” has no owner or annex |
| ADV-6 | medium | SoT drift | AD-3 allows dual emit unless Go is marked non-authoritative |

**Gate verdict from this lens:** fail.
