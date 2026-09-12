---
title: "VeGo PRD"
status: draft
created: 2026-09-12
updated: 2026-09-12
sources:
  - _bmad-output/forge/vego/forged-idea.md
  - _bmad-output/planning-artifacts/architecture/architecture-VeGo-2026-09-12/ARCHITECTURE-SPINE.md
  - docs/OVERVIEW.md
  - docs/ADR.md
---

# PRD: VeGo (Vector Go)

## 0. Document Purpose

This PRD is for PM, architecture, TEA, and implementation agents building VeGo. Vocabulary is Glossary-anchored; features carry globally numbered FRs. Product locks come from the forge (`forged-idea.md`); invariants from `ARCHITECTURE-SPINE.md` (AD-1…AD-6). Technical mechanism detail lives in `addendum.md`, not here.

## 1. Vision

VeGo is a compact, reversible source layer for Go: a minify-style IR over `go/ast` that agentic CLI/IDE hosts use when the user says “build an app…”. The LLM emits `.vego`, not Go. The toolchain expands `.vego` to standard Go and runs it with the ordinary Go compiler so behavior matches original Go.

Alpha proves lossless compact IR and a working emit→transpile→build/run path behind the host. Measured BPE token reduction is a mid-term versioned goal, not an alpha gate.

## 2. Target User

### 2.1 Jobs To Be Done

- As a builder using vibe-style CLI/IDE agents, I want generated Go projects stored in a compact IR so the agent workflow has a dedicated source surface—not raw verbose Go as the only artifact.
- As a reviewer, I want to expand `.vego` to readable Go so I can audit without forking the compiler.
- As the project owner, I want alpha success defined as “runs like Go,” with compression % deferred to a later milestone.

### 2.2 Non-Users (v1)

- Humans expecting to hand-write `.vego` as a daily language.
- Teams needing a custom Go runtime or `cmd/compile` fork.
- Buyers who treat ≥60% token savings as a v1 contractual SLA.

### 2.3 Key User Journeys

- **UJ-1. Rafa ships a tiny Go app through the agent CLI.** Rafa asks the host CLI to “make a hello HTTP server.” The host’s LLM writes `.vego` into the product tree under `src/` conventions for the app being built. Rafa runs `vego` build/run (or host-equivalent); the binary behaves as ordinary Go. He expands once with `vego fmt` to skim Go before merge.

- **UJ-2. Rafa audits a change without reading glyphs.** Rafa reviews a PR of `.vego` files via expand-to-Go (fmt / textconv when available) and confirms logic in Go form; he does not edit `.vego` by hand.

## 3. Glossary

- **VeGo** — The compact reversible IR and its toolchain (transpile, CLI, MCP), not a Go compiler fork.
- **`.vego`** — Product source file format emitted by the LLM / written by the host pipeline.
- **Host** — CLI or IDE agent environment (e.g. Cursor, Codex, vibe-style tools) that takes NL orders and drives VeGo.
- **Transpile** — Lossless transform between `.vego` and Go (`go/ast`-equivalent) for the supported grammar subset.
- **Alpha** — First shippable milestone: compact + reversible + build/run fidelity; no BPE % gate.
- **Mid-term BPE goal** — Later versioned objective to measure token reduction; not an Alpha acceptance criterion.

## 4. Features

### 4.1 Compact IR and round-trip

**Description:** Define and implement the `.vego` surface and bijective Transpile for an Alpha grammar subset. Realizes UJ-1. Governed by AD-2, AD-4.

**Functional Requirements:**

#### FR-1: Emit and store `.vego`

The Host/LLM can persist agent-generated product source as `.vego` (not Go-primary). Realizes UJ-1.

**Consequences (testable):**
- A successful generation path produces at least one `.vego` artifact as the product source of record.
- Go-only emission without `.vego` is out of Alpha’s happy path.

#### FR-2: Lossless Transpile to Go

The toolchain can Transpile `.vego` → Go for the Alpha grammar subset such that `go build` / `go run` succeed when the IR is valid. Realizes UJ-1. Governed by AD-1, AD-2.

**Consequences (testable):**
- For fixture programs in the Alpha subset, expand-then-`go test`/`go run` matches expected behavior.
- Invalid `.vego` fails closed with an explicit error (no silent bad Go).

#### FR-3: Round-trip fidelity

For the Alpha subset, Go → `.vego` → Go preserves `go/ast` semantic equivalence. Governed by AD-2.

**Consequences (testable):**
- Automated round-trip tests pass on the Alpha fixture set.

### 4.2 Host-facing CLI and MCP

**Description:** Expose Transpile and related operations to Hosts via CLI and MCP without putting grammar ownership in those layers. Realizes UJ-1, UJ-2. Governed by AD-3, AD-4.

**Functional Requirements:**

#### FR-4: CLI orchestration

`cmd/vego` can build/run/fmt paths that orchestrate Transpile then the standard Go toolchain. Realizes UJ-1.

**Consequences (testable):**
- Documented commands invoke expand + `go` tools; no custom compiler binary.

#### FR-5: MCP tool surface

`pkg/mcp` exposes Host-callable tools for reading/writing/patching `.vego` context without requiring the LLM to paste full Go sources as the only mode. Realizes UJ-1. `[ASSUMPTION: Alpha MCP tool set is read_vego_context, read expanded Go, and structural patch — exact names may match ADR examples.]`

**Consequences (testable):**
- MCP server starts and lists tools; each tool returns structured success/error.

### 4.3 Human audit bridge

**Description:** Humans inspect via expanded Go, not glyph editing. Realizes UJ-2. Governed by AD-6.

**Functional Requirements:**

#### FR-6: Expand-to-Go for review

A developer can expand `.vego` to formatted Go for audit. Realizes UJ-2.

**Consequences (testable):**
- `vego fmt` (or equivalent) writes or prints valid `gofmt`-clean Go for Alpha fixtures.

## 5. Non-Functional Requirements

#### NFR-1: Compiler compatibility

Alpha outputs must compile with the standard Go toolchain pinned in the module (`go` version in `go.mod`). Governed by AD-1.

#### NFR-2: Alpha gate excludes BPE %

Alpha acceptance does **not** require ≥60% (or any) BPE token reduction. Governed by AD-5.

#### NFR-3: Mid-term BPE measurability

When the mid-term milestone opens, the system must support measuring token counts for `.vego` vs Go on agreed tokenizer(s). `[ASSUMPTION: tokenizer list chosen at that milestone; tiktoken-go is the measurement library.]`

#### NFR-4: Methodology isolation

No BMAD method files under `src/`. Product code stays under `src/pkg/*` and `src/cmd/vego`.

## 6. Success Metrics

| Phase | Metric | Target |
| --- | --- | --- |
| Alpha | Round-trip + `go build`/`run` on Alpha fixtures | 100% fixtures green |
| Alpha | Syntactic / semantic fail-closed rate on invalid IR | No silent bad binaries |
| Mid-term | BPE reduction vs Go on sample corpus | Versioned goal (not Alpha); track in release notes |
| Counter-metric | Time-to-first-green Alpha fixture path | Prefer shipping fidelity over premature tokenizer tuning |

## 7. Scope Boundaries

**In Alpha:** IR subset, Transpile, CLI orchestration, MCP stubs→minimal tools, expand-to-Go, fixtures/tests for fidelity.

**Out of Alpha:** Full Go grammar, git textconv driver, guaranteed multi-tokenizer 1-token maps, contractual token-savings SLA, human `.vego` IDE language service.

## 8. Open Questions

- Exact Alpha grammar subset (which `go/ast` nodes first)? Owner: Architect/Dev; revisit at first epic breakdown.
- Which Host is the first integration target (Cursor MCP vs plain CLI)? Owner: PM; revisit after CLI MVP.
- Mid-term tokenizer contract list? Owner: TEA/Architect; revisit when AD-5 mid-term milestone starts.

## 9. Assumptions Index

- `[ASSUMPTION]` MCP Alpha tool trio mirrors ADR naming (FR-5).
- `[ASSUMPTION]` tiktoken-go for mid-term measurement (NFR-3).
- Stack pins in spine seed (Participle/Cobra/mcp-go) apply at first dependency pull.
