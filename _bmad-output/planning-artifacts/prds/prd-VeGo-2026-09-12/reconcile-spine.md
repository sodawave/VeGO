---
input: ARCHITECTURE-SPINE.md
slug: spine
source: _bmad-output/planning-artifacts/architecture/architecture-VeGo-2026-09-12/ARCHITECTURE-SPINE.md
compared_to:
  - prd.md
  - addendum.md
created: '2026-09-12'
---

# Reconcile extract — Architecture Spine → PRD

Spine is authoritative for invariants (per addendum). This extract checks what the spine carries that `prd.md` / `addendum.md` under-capture or drop—especially structural ideas the FR list may silence.

## Covered (no action)

| Spine item | Where in PRD / addendum |
| --- | --- |
| AD-1 transpile-only / no `cmd/compile` fork | Vision, Non-users, FR-2, NFR-1, addendum Mechanism |
| AD-2 bijective round-trip | FR-2, FR-3, Success Metrics Alpha |
| AD-3 LLM emits `.vego` (Go expansion-only) | Vision, FR-1, Rejected alternatives |
| AD-4 package boundaries | FR-4/5 descriptions, NFR-4, addendum Layout |
| AD-5 Alpha = fidelity; BPE % mid-term | Vision, NFR-2/3, Success Metrics, Out of Alpha |
| AD-6 human audit via expand-to-Go; textconv deferred | UJ-2, FR-6, Out of Alpha (git textconv) |
| Stack pins (Participle/Cobra/mcp-go; tiktoken mid-term) | Assumptions Index, addendum Mechanism |
| Fail-closed transform errors | FR-2 consequences |
| Deferred: full grammar, textconv, Unicode↔keyword maps | Scope Boundaries Out of Alpha |

## Gaps (spine → PRD drop / under-spec)

### Gap 1 — Design paradigm name absent

Spine opens with **pipes-and-filters source-to-source IR** as the named paradigm. PRD Vision/UJ narrate the flow (NL → LLM → `.vego` → transpile → `go`) but never name or bind the paradigm. Downstream epics may invent a different shape without that lock.

**Suggest:** One sentence in Vision or addendum Mechanism naming pipes-and-filters and pointing at the spine mermaid.

### Gap 2 — Ephemeral Go artifact convention

Spine Consistency: expanded `.go` is **ephemeral or audit artifact** unless explicitly checked in for bridge workflows. PRD FR-6 / UJ-2 require expand-to-Go for review but do not state whether expanded Go is checked in, gitignored, or host-local-only.

**Suggest:** Clarify in FR-6 consequences or Glossary / Data & formats note.

### Gap 3 — Deferred LLM calibration / system prompt

Spine Deferred: **LLM calibration dataset / system prompt for `.vego` emission — after CFG stabilizes.** PRD Open Questions cover Alpha grammar subset and first Host, but not calibration/prompt as an explicit out-of-Alpha / mid-term workstream. AD-3 emission contract has no product FR for “how hosts are taught to emit `.vego`.”

**Suggest:** Add to Out of Alpha or Open Questions with owner + revisit (post-CFG).

### Gap 4 — MCP operational deploy topology

Spine Deferred: **Operational deploy topology for MCP hosting — not owned at this altitude.** FR-5 assumes MCP server starts and lists tools; Scope Boundaries do not mark hosting/ops as out of Alpha. Risk: stories invent deploy topology without a spine decision.

**Suggest:** Explicit Out of Alpha (or NFR note): Alpha = in-process / local MCP stub; hosting topology deferred.

### Gap 5 — AD-4 Prevent: circular deps / grammar ownership sharpness

Spine AD-4 **Prevents:** MCP/CLI owning grammar or symbol maps; **circular deps into methodology trees.** PRD says grammar lives out of MCP/CLI (feature 4.2) and NFR-4 bans BMAD under `src/`, but never states the circular-deps prevent or that `pkg/ast` alone owns CFG/symbol maps as a product rule (only addendum Layout lists packages).

**Suggest:** One Prevents-style line under NFR-4 or FR-4/5: CLI/MCP must not own CFG/maps; no deps from product packages into `_bmad*` / methodology trees.
