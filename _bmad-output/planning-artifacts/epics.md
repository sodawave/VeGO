---
stepsCompleted: ['step-01-validate-prerequisites', 'step-02-design-epics', 'step-03-create-stories', 'step-04-final-validation']
inputDocuments:
  - _bmad-output/planning-artifacts/prds/prd-VeGo-2026-09-12/prd.md
  - _bmad-output/planning-artifacts/prds/prd-VeGo-2026-09-12/addendum.md
  - _bmad-output/planning-artifacts/architecture/architecture-VeGo-2026-09-12/ARCHITECTURE-SPINE.md
  - _bmad-output/forge/vego/forged-idea.md
notes: 'Alpha = pre-release version (pre-versión). No UX doc. User authorized continuous Continue through workflow.'
validation: 'FR1-FR6 covered; no UX-DRs; stubs exist; forward deps OK; Alpha=pre-release'
---

# VeGo - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for VeGo, decomposing the requirements from the PRD and Architecture into implementable stories. **Alpha** is the pre-release version: fidelity-first, not BPE-%-gated.

## Requirements Inventory

### Functional Requirements

FR1: Host/LLM can persist agent-generated product source as `.vego` (not Go-primary).
FR2: Toolchain can Transpile `.vego` → Go for the Alpha grammar subset such that `go build` / `go run` succeed when IR is valid; invalid IR fails closed.
FR3: For Alpha grammar subset, Go → `.vego` → Go preserves `go/ast` semantic equivalence; tooling supports both directions; fixtures include hello HTTP + struct/`range`.
FR4: `cmd/vego` orchestrates build/run/fmt via Transpile then standard Go toolchain (no custom compiler).
FR5: `pkg/mcp` exposes Alpha tools `read_vego_context`, read expanded Go, and structural patch; all three succeed on fixtures.
FR6: Developer can expand `.vego` to `gofmt`-clean Go for audit (`vego fmt` or equivalent).

### NonFunctional Requirements

NFR1: Alpha outputs compile with standard Go toolchain pinned in `go.mod`.
NFR2: Alpha acceptance does not require any BPE token-reduction percentage.
NFR3: Mid-term milestone must support measuring token counts (tiktoken-go assumed); not Alpha-blocking.
NFR4: No BMAD under `src/`; CLI/MCP must not own CFG/symbol maps; packages stay under `src/pkg/*` and `src/cmd/vego`.

### Additional Requirements

- AD-1: Transpile-only path; never fork `cmd/compile`.
- AD-2: Bijective Go ↔ `.vego` for supported subset.
- AD-3: LLM emission contract — product source is `.vego`.
- AD-4: Package boundaries — `pkg/ast`, `pkg/transpiler`, `pkg/mcp`, `cmd/vego`.
- AD-5: Alpha = lossless fidelity; BPE % mid-term only.
- AD-6: Humans audit via expand-to-Go; not primary `.vego` authors.
- Alpha grammar subset locked: package, import, func, params/results, basic+named types, calls, selectors, literals, short decl, if/else, for+range, return, struct types. OUT: generics, concurrency, reflection, cgo, unsafe, aggressive identifier minify.
- Primary Host for Alpha = CLI; MCP minimal but present.
- Stubs already exist under `src/`; no separate greenfield starter template beyond Go module.

### UX Design Requirements

None — CLI/MCP product; no UX design contract.

### FR Coverage Map

- FR1: Epic 1 — Persist `.vego` as product source of record
- FR2: Epic 1 — `.vego` → Go builds/runs (fail closed)
- FR3: Epic 1 — Round-trip fidelity on Alpha grammar + fixtures
- FR4: Epic 2 — CLI build/run/fmt orchestration
- FR5: Epic 3 — MCP three-tool Alpha surface
- FR6: Epic 2 — Expand-to-Go audit path
- NFR1–NFR2, NFR4: Epic 1 (toolchain + boundaries); NFR3 deferred mid-term

## Epic List

### Epic 1: Lossless Compact IR (Alpha Core)
Builders get a reversible `.vego` IR for the Alpha grammar subset that expands to real Go and round-trips without semantic loss.
**FRs covered:** FR1, FR2, FR3 (NFR1, NFR2, NFR4)

### Epic 2: CLI Host Path
Through `cmd/vego`, builders build/run Alpha programs and audit via expand-to-Go—no custom compiler.
**FRs covered:** FR4, FR6

### Epic 3: MCP Agent Surface
Agent Hosts call minimal MCP tools to read `.vego`, read expanded Go, and structurally patch IR on Alpha fixtures.
**FRs covered:** FR5

## Epic 1: Lossless Compact IR (Alpha Core)

Builders get a reversible `.vego` IR for the Alpha grammar subset that expands to real Go and round-trips without semantic loss.

### Story 1.1: Alpha Grammar Skeleton and Symbol Surface

As a VeGo implementer,
I want `pkg/ast` to define the Alpha grammar subset and symbol/CFG surface,
So that Transpile and fixtures share one locked IR contract.

**Acceptance Criteria:**

**Given** the Alpha grammar subset locked in the PRD  
**When** I inspect `src/pkg/ast`  
**Then** types/APIs exist covering package, import, func, params/results, basic+named types, calls, selectors, literals, short decl, if/else, for+range, return, struct types  
**And** unsupported constructs (generics, concurrency, cgo, unsafe, aggressive id minify) are explicitly rejected or undocumented as out of Alpha  
**And** no BMAD files exist under `src/` (NFR4)

### Story 1.2: `.vego` → Go Expand for Hello Fixture

As a builder,
I want valid Alpha `.vego` expanded to Go that `go build`/`go run` can execute,
So that compact IR is a real executable source path (FR2, NFR1).

**Acceptance Criteria:**

**Given** a hello-HTTP `.vego` fixture within the Alpha subset  
**When** I expand via `pkg/transpiler` (or CLI wrapper later)  
**Then** the output is valid Go that `go build` and `go run` succeed against  
**And** invalid `.vego` fails closed with an explicit error (no silent bad Go)

### Story 1.3: Go → `.vego` → Go Round-Trip

As a builder,
I want Go fixtures compacted to `.vego` and expanded back with semantic equivalence,
So that the IR is lossless for Alpha (FR3, AD-2).

**Acceptance Criteria:**

**Given** hello-HTTP and a second fixture covering struct types + `range`  
**When** I run Go → `.vego` → Go  
**Then** automated tests assert `go/ast` semantic equivalence (or agreed normalized equality)  
**And** both fixtures remain green end-to-end

### Story 1.4: Persist `.vego` as Source of Record

As an agent Host pipeline,
I want `.vego` files treated as the product source artifact on disk,
So that Go-primary emission is not the Alpha happy path (FR1, AD-3).

**Acceptance Criteria:**

**Given** an Alpha generation/write path in the toolchain  
**When** product source is written  
**Then** at least one `.vego` artifact is the source of record for the fixture/app under test  
**And** documentation/tests treat Go-only happy path as out of Alpha

## Epic 2: CLI Host Path

Through `cmd/vego`, builders build/run Alpha programs and audit via expand-to-Go—no custom compiler.

### Story 2.1: `vego build` / `vego run`

As a developer using the CLI Host,
I want `vego build` and `vego run` to expand `.vego` then invoke the standard Go toolchain,
So that I can execute Alpha programs without a custom compiler (FR4, AD-1).

**Acceptance Criteria:**

**Given** Alpha fixtures from Epic 1  
**When** I run `vego build` / `vego run` on them  
**Then** commands orchestrate expand + `go` tools only  
**And** programs behave as expected for the fixtures  
**And** no custom compiler binary is introduced

### Story 2.2: `vego fmt` Expand-to-Go Audit

As a reviewer,
I want `vego fmt` (or equivalent) to print/write `gofmt`-clean Go from `.vego`,
So that I can audit without reading glyphs (FR6, AD-6).

**Acceptance Criteria:**

**Given** each Alpha fixture `.vego`  
**When** I run `vego fmt`  
**Then** output is valid `gofmt`-clean Go  
**And** the command does not require editing `.vego` by hand

## Epic 3: MCP Agent Surface

Agent Hosts call minimal MCP tools to read `.vego`, read expanded Go, and structurally patch IR on Alpha fixtures.

### Story 3.1: MCP Server with Three Alpha Tools

As an agent Host,
I want an MCP server exposing `read_vego_context`, read-expanded-Go, and structural patch,
So that agents can operate on `.vego` without pasting full Go as the only mode (FR5, AD-3, AD-4).

**Acceptance Criteria:**

**Given** `pkg/mcp` and Alpha fixtures  
**When** the MCP server starts  
**Then** it lists the three Alpha tools (plus optional non-mutating health/list if required by SDK)  
**And** `read_vego_context` returns compact IR for a fixture  
**And** expanded-Go read returns valid Go for that fixture  
**And** structural patch changes a node without full-file rewrite  
**And** MCP does not own CFG/symbol maps (delegates to `pkg/ast` / `pkg/transpiler`)
