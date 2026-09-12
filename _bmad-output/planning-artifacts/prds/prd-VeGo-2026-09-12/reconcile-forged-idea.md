# Reconcile: forged-idea → PRD + addendum

**Input:** `_bmad-output/forge/vego/forged-idea.md`  
**Against:** `prd.md`, `addendum.md` (same folder)

## Landed in PRD / addendum

| Forge lock | Where |
| --- | --- |
| Compact reversible IR over `go/ast` (minify-style) | PRD Vision, Glossary, FR-1–3 |
| LLM emits `.vego` (not Go-primary) | Vision, FR-1, Non-Users, addendum Rejected |
| Host pipeline: NL → host → LLM → `.vego` → transpile → `go build`/`run` | Vision, UJ-1, FR-2/4 |
| Behind CLI/MCP/IDE; humans don’t casual-author; audit via expand-to-Go | FR-4–6, Non-Users, UJ-2 |
| No `cmd/compile` fork | NFR-1, addendum Mechanism, Scope Out |
| Product root `src/` (`pkg/ast`, `transpiler`, `mcp`, `cmd/vego`) | NFR-4, addendum Layout |
| Alpha = lossless + round-trip; BPE % mid-term only | Vision, NFR-2/3, Success Metrics, addendum Rejected |
| Rejected: ≥60% alpha gate; Go-emit+minify-on-save primary; human language / compiler fork | §2.2 Non-Users, NFR-2, Scope Out, addendum Rejected |

## Gaps (forge → PRD/addendum)

1. **Linear pipeline formula** — forge’s one-line `NL → CLI/IDE (vibe) → LLM → .vego → vego → go build/run → behavior = Go` is scattered across Vision/UJ/FR; not restated as a single canonical pipeline.
2. **“JS minify for Go” teaching analogy** — only soft “minify-style”; the JS parallel is not named.
3. **“Why it holds” product-surface argument** — forge’s claim that without `.vego` emit the layer has no product surface is implied by FR-1, not carried as explicit rationale.
4. **“Not a new runtime” fidelity framing** — execution-via-transpile (not custom runtime) is assumed via AD-1/NFR-1; forge’s explicit contrast is missing from PRD prose.
5. **Casual-authoring rule as positive constraint** — covered via Non-Users / Out of Alpha IDE language service; forge’s “humans do not author `.vego` casually” is not a standalone FR/NFR.

## Qualitative ideas dropped

- Formal “why it holds” narrative (product surface + transpile-not-runtime).
- Explicit JS-minify metaphor as onboarding language.
- Single equality slogan `behavior = original Go` as a locked one-liner (fidelity is metric/FR-driven instead).
