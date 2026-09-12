# VeGo (Vector Go)

BPE token-optimized, non-human-readable programming surface: a bi-directional, lossless AST serialization layer on standard Go (`go/ast`). Goal: **60–85%** token reduction with **0%** syntactic hallucination on round-trip.

See [`docs/OVERVIEW.md`](docs/OVERVIEW.md) and [`docs/ADR.md`](docs/ADR.md).

## Layout

| Path | Role |
|------|------|
| `src/` | Product source root |
| `src/pkg/ast` | Symbol mapping & CFG grammar |
| `src/pkg/transpiler` | Go ↔ VeGo transform |
| `src/pkg/mcp` | MCP tool server |
| `src/cmd/vego` | CLI entry point |
| `docs/` | Project knowledge (`bmm.project_knowledge`) |
| `_bmad/` | BMad Method modules |
| `_bmad-output/` | Planning / implementation / test artifacts |
| `.agents/skills/` | Pi / OpenCode / Antigravity CLI skills |
| `.claude/skills/` | Claude Code skills |
| `.agent/skills/` | Google Antigravity skills |
| `.opencode/commands/` | OpenCode command bindings |
| `AGENTS.md` | Process rules for coding agents |

## Prerequisites

- Go 1.22+
- Node.js 20.12+
- [uv](https://docs.astral.sh/uv/) (BMAD Python skills)
- Supported AI tool: Claude Code, OpenCode, Antigravity, or Pi (Cursor also supported)

## Get started

1. Invoke `bmad-help` and ask what to do next.
2. Prefer: **debate → architecture spine + memlog → spec or epics/stories → implement → evidence**.
3. Verify stubs: `go build ./...` and `go run ./src/cmd/vego`

Chat may be in Spanish; formal BMAD documents are English.

Authorship is human-only (Sodawave). After clone:

```bash
git config core.hooksPath .githooks
```

## Refresh BMAD

```bash
npx bmad-method install --directory . --yes \
  --tools antigravity,claude-code,pi,opencode,antigravity-cli \
  --modules bmm,cis,tea
```
