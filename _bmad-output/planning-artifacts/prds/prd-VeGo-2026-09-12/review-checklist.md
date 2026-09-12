# PRD Quality Review — VeGo (prd-VeGo-2026-09-12)

## Overall verdict

**concerns.** The PRD earns its thesis: Alpha is fidelity-first compact IR for agent hosts, BPE % is demoted, and Non-Users / Out of Alpha / rejected alternatives (addendum) make the trade-offs honest. What blocks a clean pass is done-ness: FR-1–FR-3 and Alpha success metrics depend on an undefined Alpha grammar subset and fixture set, so engineering cannot yet know “done” without inventing scope. Fix that one open question (or park it as an explicit pre-epic gate) and tighten FR-5 / Host-target notes, and this is decision-ready for story breakdown.

## Decision-readiness — adequate

Trade-offs are stated as decisions, not smoothed: Vision and NFR-2 lock Alpha without a BPE gate; §7 Out of Alpha and Non-Users (§2.2) name contractual ≥60% SLA, compiler fork, and human-written `.vego` as excluded; the addendum’s Rejected alternatives repeats the forge demotions with what was given up. AD references (AD-1…AD-6) connect product locks to the spine without re-litigating mechanism in the PRD body.

Gaps that keep this from *strong*: §8’s first open question (“Exact Alpha grammar subset”) is load-bearing for every fidelity FR and the Alpha success table, yet it is framed as “revisit at first epic breakdown” without a `[NOTE FOR PM]` that shipping green-light is blocked until that subset exists. The Host integration target (Cursor MCP vs plain CLI) is correctly open, but there is no PM callout at the tension between FR-5 (MCP in Alpha) and “CLI MVP first” revisit timing — a decision-maker could read FR-5 as requiring MCP parity with CLI in the same milestone.

### Findings

- **high** Alpha grammar subset left as deferred epic detail (§8, FR-1–FR-3, §6) — Decision-makers cannot approve a fixed Alpha scope; fidelity metrics are circular until fixtures are named. *Fix:* Promote to a pre-epic gate with an owner deliverable (node allowlist + fixture list), or mark Alpha scope as “subset TBD; no FR sign-off until Architect publishes allowlist.”
- **medium** Host-first target vs MCP-in-Alpha tension undermarked (§8 Q2, FR-4, FR-5) — No `[NOTE FOR PM]` at whether Alpha ships CLI-only with MCP stubs or requires callable MCP tools. *Fix:* Add a PM note: Alpha MCP = stub list vs must-work tools; defer Host brand choice without deferring the capability bar.

## Substance over theater — strong

Lean document: two UJs, one named protagonist (Rafa), JTBD that map to FRs, Vision that would not swap into a generic “AI coding tool” PRD (“LLM emits `.vego`,” “expand… ordinary Go compiler,” “BPE… not an alpha gate”). NFRs are product-specific (compiler pin, BPE exclusion, BMAD isolation) rather than scalable/secure/reliable boilerplate. Innovation claims stay mechanistic and pointed at the addendum/spine; no differentiation theater section. Personas are Jobs-shaped, not a cast of four archetypes.

### Findings

- *(none material)* JTBD lines are short but drive Non-Users and FR grouping; leave as-is.

## Strategic coherence — strong

Thesis is explicit: agent hosts need a dedicated compact reversible source surface; Alpha proves lossless IR + emit→transpile→build/run; compression is mid-term. Feature groups (4.1 IR, 4.2 host surface, 4.3 audit bridge) follow that arc. Success Metrics table aligns phases with the thesis and includes a counter-metric against premature tokenizer tuning (AD-5). MVP kind is problem-solving / capability (fidelity path), not a backlog dump. Addendum Rejected alternatives keep the forge demotions visible so epic planning does not re-open Go-primary emit or BPE-as-Alpha-gate.

### Findings

- **low** Counter-metric wording is slightly awkward (§6) — “Time-to-first-green… Prefer shipping fidelity…” reads as a preference, not a measurable counter. *Fix:* Phrase as “Do not delay Alpha fixture green path for tokenizer experiments” with a tracking rule (e.g. no tiktoken work in Alpha milestone tickets).

## Done-ness clarity — thin

FR-1–FR-4 and FR-6 carry testable consequences (artifact of record, fail-closed invalid IR, round-trip tests, CLI invokes expand + `go`, `vego fmt` → gofmt-clean). That is the right pattern. The weakness is shared dependency on undefined inputs: “Alpha grammar subset,” “fixture programs,” and “Alpha fixture set” appear in FR-2, FR-3, FR-6, and §6 without enumeration or pointer to a companion list. An engineer cannot close FR-2/FR-3 without inventing which `go/ast` nodes count.

FR-5 consequences stop at “MCP server starts and lists tools; each tool returns structured success/error” — that accepts a hollow tool list and does not require the assumed trio (`read_vego_context`, expanded Go, structural patch) to actually mutate or return IR. NFR-1 requires “`go` version in `go.mod`” but the PRD never states the pin (acceptable if spine/module is SOLE, but then say so). NFR-3 correctly gates mid-term measurability behind assumptions.

### Findings

- **critical** Fidelity FRs and Alpha SMs are untestable until subset/fixtures exist (FR-2, FR-3, FR-6, §6) — “100% fixtures green” with zero listed fixtures is theater of precision. *Fix:* Add an Alpha Fixture Contract (even draft): e.g. hello main, package with func, one import, fail-closed bad glyph — or explicitly “fixtures = TBD by Architect before story ready.”
- **high** FR-5 acceptance too weak (§4.2) — Listing tools ≠ Host-callable read/patch of `.vego`. *Fix:* Consequences must require each assumed Alpha tool to succeed on a fixture (read IR, return expanded Go, apply one structural patch) with structured error on invalid input.
- **medium** “semantic equivalence” / “matches expected behavior” without oracle (§ FR-2, FR-3) — Round-trip via `go/ast` is the right bar but needs a stated check (e.g. `go/types` or AST deep-equal after normalize). *Fix:* One sentence naming the equivalence oracle for Alpha tests.

## Scope honesty — adequate

§2.2 Non-Users, §7 Out of Alpha, and addendum Rejected alternatives do real work: full grammar, textconv, multi-tokenizer 1-token maps, SLA, human IDE LS, Go-primary path, compiler fork. Assumptions Index (§9) exists; FR-5 and NFR-3 are tagged inline. Open Questions are mostly real (not rhetorical). Document Purpose correctly pushes mechanism to addendum.

Open-item density is moderate for a greenfield capability PRD: three OQs, three assumptions, zero `[NOTE FOR PM]`. For a draft feeding epics that is acceptable; for a “build Alpha now” green light, the grammar OQ density is a blocker (see Decision-readiness / Done-ness). Missing: explicit `[NON-GOAL for MVP]` tags on items that could creep (git textconv, IDE LS) — they appear in Out of Alpha but not callout-tagged beside FRs that might tempt them (FR-6 audit bridge).

### Findings

- **medium** No `[NOTE FOR PM]` on deferred blockers (§8) — High-stakes deferrals look like ordinary backlog. *Fix:* Tag grammar subset and Host/MCP Alpha bar with `[NOTE FOR PM]`.
- **low** Stack pins assumption only in Index (§9) — “Participle/Cobra/mcp-go… first dependency pull” has no inline `[ASSUMPTION]` in body. *Fix:* Tag once under FR-4/FR-5 or §5, or drop from Index if spine-only.

## Downstream usability — adequate

Glossary covers VeGo, `.vego`, Host, Transpile, Alpha, Mid-term BPE goal; FR/UJ IDs are contiguous (UJ-1–2, FR-1–6, NFR-1–4). UJs name Rafa inline. Cross-refs to AD-* and UJs resolve conceptually. Shape is pullable for architecture (already has spine) and stories, except the grammar/fixture hole will force Architect/Dev to invent acceptance before story Ready.

Weak spots for extractors: “Host/LLM” vs Glossary “Host”; “minify-style IR” (Vision) is not Glossary-defined (harmless if addendum-owned); FR-5 tool names live only in an assumption, not Glossary or a named Alpha MCP surface list.

### Findings

- **medium** Alpha MCP tool surface not glossary- or FR-stable (§ FR-5, §9) — Downstream MCP stories will bikeshed names. *Fix:* Freeze Alpha tool names in Glossary or a short “Alpha MCP surface” bullet under FR-5.
- **low** Vision “minify-style IR” vs Glossary “compact reversible IR” — Synonym drift for extractors. *Fix:* Prefer Glossary wording in Vision.

## Shape fit — strong

VeGo is a technical capability / host-toolchain product, not a multi-stakeholder consumer UX. Light UJ density (two paths, one protagonist) and operational Alpha metrics fit; over-formalized journey maps would be wrong here. Chain-top intent (PM → architecture/TEA → implement) is declared in §0; mechanism correctly lives in addendum + spine. Brownfield stubs under `src/` are not contradicted; PRD does not invent compiler-fork behavior. Stakes match a draft capability PRD: rigor on invariants and Non-Goals, lighter on persona theater.

### Findings

- *(none material)*

## Mechanical notes

- **Assumptions Index roundtrip:** FR-5 and NFR-3 inline tags match Index bullets 1–2. Index bullet 3 (stack pins) has **no** matching inline `[ASSUMPTION]` in the body.
- **ID continuity:** UJ-1–2, FR-1–6, NFR-1–4 contiguous; no duplicates; UJ/FR cross-refs resolve.
- **Glossary drift:** “Host/LLM,” “minify-style IR,” “glyphs” (UJ-2 / FR-6 area) vs Glossary terms — minor.
- **UJ protagonists:** Both UJs name Rafa with enough context.
- **Required sections:** Vision, users/UJs, Glossary, FRs, NFRs, SMs, scope, OQs, Assumptions present; addendum carries rejected alternatives and mechanism — appropriate split.
- **Output path note:** Rubric template names `review-rubric.md`; this review is written to `review-checklist.md` per task request.
