# Alpha status — VeGo

**Date:** 2026-09-12  
**Gate meaning:** Alpha = **pre-release**. Success = fidelity / round-trip, not BPE %.

## Done

| Area | Evidence |
|------|----------|
| Keyword↔glyph bijective map | `src/pkg/ast` — 25 keywords, RoundTripKeywords |
| Alpha grammar fixture | `src/pkg/transpiler/testdata/*.go` + `*.vego` |
| Go → `.vego` → Go | `TestRoundTripFixture`, `TestRoundTripIdempotent` |
| CLI `fmt` / `build` / `run` | `src/cmd/vego` + CLI tests |
| MCP three tools | `src/pkg/mcp` + JSON-RPC stdio `ServeStdio` |
| Mid-term rough tokens | `src/pkg/bpe` + `vego tokens` (not Alpha gate) |
| Public examples | `src/examples/` — hello, HTTP, struct/range, control |
| Architecture spine | `_bmad-output/planning-artifacts/architecture/architecture-VeGo-2026-09-12/` |
| PRD final | `_bmad-output/planning-artifacts/prds/prd-VeGo-2026-09-12/` |
| Epics / sprint | `epics.md`, `sprint-status.yaml` (Alpha stories done) |

## Deferred (mid-term)

- Real tiktoken / BPE corpus validation of glyph map
- Full AST-aware edits (beyond keyword substitution)
- Generics / concurrency in Alpha grammar
- Git textconv / IDE Host packaging
- Public Pages deploy of `web/` (artifact ready)

## How to verify

```bash
go test ./src/...
go run ./src/cmd/vego fmt src/pkg/transpiler/testdata/hello_http.go
```

See [DEVELOPMENT.md](./DEVELOPMENT.md).
