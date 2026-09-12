# VeGo (Vector Go)

**BPE-oriented, bi-directional, lossless compact IR over Go AST** — so AI coding agents can emit denser code without forking the Go compiler.

> **Alpha = pre-release.** The gate is **fidelity** (round-trip Go ↔ `.vego`), not compression %. Token savings are a **mid-term** product goal.

## Idea in one line

LLM emits **`.vego`** (not Go) → VeGo transpiles → `go build` / `go run`. Like a minify layer for Go IR.

```
NL → CLI/IDE Host → LLM → .vego → transpile → go tool → Go behavior
```

## Quick start

```bash
git clone https://github.com/sodawave/VeGO.git
cd VeGO
go test ./src/...
go run ./src/cmd/vego fmt path/to/file.go
go run ./src/cmd/vego build path/to/file.vego
go run ./src/cmd/vego run path/to/file.vego
```

MCP (stdio JSON-RPC):

```bash
go run ./src/cmd/vego mcp stdio
```

## Docs

| Doc | Purpose |
|-----|---------|
| [docs/OVERVIEW.md](docs/OVERVIEW.md) | Product vision |
| [docs/ADR.md](docs/ADR.md) | Early decisions |
| [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) | BMAD trail, CLI, tests |
| [docs/ALPHA-STATUS.md](docs/ALPHA-STATUS.md) | What Alpha shipped |
| [web/](web/) | Public showcase site |

Planning artifacts (English): `_bmad-output/planning-artifacts/`  
Forge debate: `_bmad-output/forge/vego/`

## Layout

```
src/cmd/vego     CLI + MCP entry
src/pkg/ast      Keyword↔glyph map
src/pkg/transpiler  Round-trip engine
src/pkg/mcp      MCP tools
src/pkg/bpe      Rough token estimator (mid-term)
docs/            Knowledge
web/             Static public site
_bmad-output/    Specs, spine, sprint
```

## Authorship

Conception and architecture: **Sodawave / SodaWave**. Agents assist only — no AI co-author trailers on commits/PRs.

## License

See repository license file when published.
