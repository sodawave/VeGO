---
name: review-tech-versions
type: architecture-spine-review
lens: tech-versions / brownfield reality-check
spine: architecture-VeGo-2026-09-12/ARCHITECTURE-SPINE.md
reviewed: '2026-09-15'
verdict: pass
sources_checked:
  - _bmad-output/planning-artifacts/architecture/architecture-VeGo-2026-09-12/ARCHITECTURE-SPINE.md
  - go.mod
  - go.sum (hashes for pinned modules)
  - src/ import graph (brief scan)
  - go proxy module version lists (live)
  - go.dev release notes / dl (Go 1.27)
---

# Tech-versions review — VeGo Architecture Spine

## Verdict

**pass** — Committed stack pins that are in force today match brownfield (`go.mod` + `src/` imports). Deferred ADR seeds are correctly labeled as deferred, still exist on the module proxy, and are not asserted as current implementation. No named technology is dead or inventively versioned. Residual issues are hygiene / doc drift, not false stack claims.

## Method

Reality-check only (no ADR deep rewrite):

1. Compare spine **Stack** table to `go.mod` and runtime `go version`.
2. Confirm each named third-party module still resolves and note latest proxy versions.
3. Brief `src/` import scan: what is actually linked vs what ADR historically proposed.
4. Cross-check Structural Seed / Capability map against packages on disk.

## Decision-by-decision check

| Decision / claim | Reality check | Status |
| --- | --- | --- |
| Go `1.27.0` (`go.mod`) | Runtime is `go1.27.0`; `go.mod` directive `go 1.27.0`. Go 1.27.0 released 2026-08-19; patch **1.27.1** is now on go.dev/dl as a newer stable patch. Language pin in module is still accurate for this tree. | OK (note patch lag) |
| stdlib `go/ast`, `go/parser`, `go/scanner`, `go/printer`, `go/token` | Transpile path uses `go/ast`, `go/parser`, `go/scanner`, `go/token`, and **`go/format`** (not a direct `go/printer` import). `go/format` is the idiomatic wrapper around printer — claim fits intent; table is slightly imprecise. | OK (minor imprecision) |
| `pkoukk/tiktoken-go` **v0.1.8** | Present in `go.mod` / `go.sum`. Proxy latest is still **v0.1.8**. Directly imported by `src/pkg/bpe` (and tests). Fits mid-term / Beta measurement role; not an Alpha gate — matches AD-5. | OK (see dep hygiene) |
| `alecthomas/participle/v2` deferred; hand-rolled scanner today | No participle imports under `src/`. Transpiler uses `go/scanner` + custom transform. Module still publishes (latest proxy **v2.1.4**). Deferred framing is correct. | OK |
| `spf13/cobra` deferred; `flag`-based CLI today | `cmd/vego` is a manual `switch` on `os.Args` (no cobra / no `flag` package either — spine says “`flag`-based” loosely for hand CLI). Module still publishes (latest **v1.10.2**). Deferred framing is correct. | OK (wording nit) |
| `mark3labs/mcp-go` deferred; minimal JSON-RPC stdio today | `pkg/mcp` is hand-rolled JSON-RPC (`initialize`, `tools/list`, `tools/call`); no mcp-go import. Module still exists and has matured to **v1.1.0**. Deferred framing is correct; any future adoption must re-pin, not copy ADR-era assumptions. | OK (re-pin on adopt) |
| Module path `github.com/sodawave/VeGO` | Matches `go.mod` and all internal imports. | OK |
| Layout `cmd/vego`, `pkg/ast`, `pkg/transpiler`, `pkg/mcp` | Present on disk; import edges match AD-4 (`cmd`→`mcp`/`transpiler`/`bpe`; `mcp`→`transpiler`; `transpiler`→`ast`/`bpe`). | OK |
| Capability map `src/pkg/bpe` + tiktoken | Package exists and is the real tiktoken consumer. | OK |
| AD-1 transpile-only / no `cmd/compile` fork | CLI shells out to `go build` / `go run`; no compiler fork in tree. | OK |

## Findings

### F1 — `tiktoken-go` is a direct dependency but marked `// indirect` in `go.mod` (hygiene)

- **Severity:** low (does not invalidate the spine pin)
- **Evidence:** `src/pkg/bpe/bpe.go` imports `github.com/pkoukk/tiktoken-go`; `go mod why` resolves through `src/pkg/bpe`. `go.mod` lists `v0.1.8 // indirect`.
- **Impact:** Version string in the spine is correct; the *classification* in `go.mod` is not. `go mod tidy` should promote it to a direct `require`.
- **Action:** Fix in module maintenance (outside spine text); optional one-line note in Stack that it is a **direct** Beta measurement dep.

### F2 — Structural Seed omits `src/pkg/bpe` while Capability map includes it (brownfield doc drift)

- **Severity:** low–medium (spine internal consistency)
- **Evidence:** Seed tree lists only `cmd/vego`, `pkg/ast`, `pkg/transpiler`, `pkg/mcp`. Disk and Capability → Architecture Map both include `pkg/bpe`.
- **Impact:** Implementers following only the Seed can miss the measurement package that AD-5 / Stack already acknowledge.
- **Action:** Add `pkg/bpe/` to Structural Seed on next spine edit.

### F3 — Stack lists `go/printer`; brownfield imports `go/format` (imprecision)

- **Severity:** low
- **Evidence:** `rg` finds no `go/printer` under `src/`; transpiler uses `go/format` (+ `go/ast`, `go/parser`, `go/scanner`, `go/token`).
- **Impact:** Not a false technology — printer remains the formatting substrate — but the table overstates a direct dependency.
- **Action:** Prefer wording `go/format` (via printer) in Stack.

### F4 — Deferred libraries still fit; versions must be re-checked on adoption (esp. mcp-go)

- **Severity:** informational (spine correctly deferred)
- **Evidence (proxy, 2026-09-15):**
  - `participle/v2` → latest **v2.1.4** (still active)
  - `cobra` → latest **v1.10.2** (still de-facto CLI toolkit)
  - `mark3labs/mcp-go` → latest **v1.1.0** (large jump from early 0.x ADR-era seeds)
- **Impact:** Spine does **not** pin stale versions for deferred items — good. ADR knowledge docs still name these as primary stack without versions; spine correctly overrides with “deferred seed / hand-rolled today.” Future adoption of mcp-go should treat **v1.x** APIs as the baseline, not early 0.x snippets.
- **Action:** When un-deferring, pin concrete versions into Stack + `go.mod` after a fresh proxy check; do not resurrect ADR sample import paths blindly.

### F5 — Go patch line: module/runtime on `1.27.0` while `1.27.1` exists (freshness note)

- **Severity:** informational
- **Evidence:** Spine + `go.mod` + environment `go1.27.0`; go.dev/dl lists **go1.27.1** as a current stable patch.
- **Impact:** `go 1.27.0` language version remains valid. Not out-of-date in the sense of claiming a non-existent release. Optional toolchain bump is product hygiene, not an architecture fail.
- **Action:** Optional: bump CI/toolchain to 1.27.1; keep `go 1.27.0` (or `1.27`) directive unless intentionally raising minimum.

### F6 — CLI wording “`flag`-based” vs actual hand `os.Args` switch (nit)

- **Severity:** nit
- **Evidence:** `cmd/vego/main.go` switches on `args[0]`; no `flag` or cobra import.
- **Impact:** Brownfield “not cobra” claim is true; “flag-based” is slightly overstated.
- **Action:** Prefer “hand `os.Args` dispatcher” in Stack / prose.

## What was *not* found

- No hallucinated library names.
- No pinned version that does not exist on the module proxy.
- No brownfield contradiction of AD-1–AD-6 package boundaries or transpile-only rule.
- No accidental import of deferred ADR libraries that would make “deferred” a lie.

## Summary table — named technologies

| Name | Spine claim | Exists | Fits role | Current / noted latest | Brownfield match |
| --- | --- | --- | --- | --- | --- |
| Go | 1.27.0 | yes | yes | 1.27.1 patch available | yes |
| go/ast family | stdlib transpile | yes | yes | stdlib | yes (`format` not `printer`) |
| pkoukk/tiktoken-go | v0.1.8 measure | yes | yes | v0.1.8 (= latest) | yes (mark direct) |
| participle/v2 | deferred | yes | yes if adopted | v2.1.4 | deferred correctly |
| cobra | deferred | yes | yes if adopted | v1.10.2 | deferred correctly |
| mark3labs/mcp-go | deferred | yes | yes if adopted | v1.1.0 | deferred correctly |

## Recommended spine deltas (non-blocking)

1. Structural Seed: add `pkg/bpe/`.
2. Stack: `go/format` instead of / in addition to `go/printer`; clarify CLI as hand `os.Args`.
3. Optionally note deferred latests as “re-pin on adopt” without treating them as committed pins.
4. Separate `go.mod` tidy for direct `tiktoken-go` (code change, not spine-only).

## Reviewer conclusion

Committed decisions in this spine were **reality-checked against the tree and live module/toolchain data**, not merely restated from ADR training-style aspirational stack. Verdict **pass** with the hygiene and doc-drift findings above; nothing in the Stack table is obsolete enough to block finalize.
