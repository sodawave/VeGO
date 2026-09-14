# VeGo — Development Guide

Human-facing documentation of what was built, how BMAD was followed, and how to verify Alpha.

## What VeGo is

VeGo is a **lossless compact IR** for Go: LLM/host emits `.vego`, toolchain expands to standard Go, `go build`/`run` executes. Not a compiler fork.

**Alpha = pre-release (pre-versión):** fidelity and round-trip first. BPE % savings are a **mid-term** metric (`vego tokens`), not an Alpha gate.

## BMAD trail (this repo)

| Stage | Artifact |
|-------|----------|
| Forge (HARDENED) | `_bmad-output/forge/vego/` |
| Architecture spine | `_bmad-output/planning-artifacts/architecture/architecture-VeGo-2026-09-12/` |
| PRD (final) | `_bmad-output/planning-artifacts/prds/prd-VeGo-2026-09-12/prd.md` |
| Epics + stories | `_bmad-output/planning-artifacts/epics.md` |
| Sprint status | `_bmad-output/implementation-artifacts/sprint-status.yaml` |
| Alpha status | [`docs/ALPHA-STATUS.md`](./ALPHA-STATUS.md) |
| Public showcase | [`web/`](../web/) |

Cycle used: **forge → spine → PRD → epics/stories → sprint → build → evidence (tests)**. Next public layer: this guide + `web/`.

## Packages

| Path | Role |
|------|------|
| `src/pkg/ast` | Keyword↔glyph map, Alpha grammar surface |
| `src/pkg/transpiler` | Go ↔ `.vego` transform |
| `src/pkg/mcp` | MCP JSON-RPC tools (`vego mcp stdio`) |
| `src/pkg/bpe` | Mid-term token estimate heuristic |
| `src/cmd/vego` | CLI: `fmt`, `build`, `run`, `tokens`, `mcp` |
| `src/examples` | Runnable Alpha samples (`.go` + `.vego`) — see [`src/examples/README.md`](../src/examples/README.md) |

## CLI

```bash
go test ./src/...
go build -o vego ./src/cmd/vego

./vego fmt path/to/file.go > file.vego   # .vego is single-line (ASI → ';')
./vego fmt file.vego          # expand to Go on stdout
./vego build file.vego
./vego run file.vego
./vego tokens file.go         # JSON estimate (not an Alpha gate)
./vego mcp stdio              # JSON-RPC tool server
```

`.vego` output has **no newlines**: scanner ASI newlines are emitted as `;` so the IR stays one line. Expand still round-trips (gofmt may reflow braces).

## Tests (TEA evidence)

```bash
go test ./src/... -cover
```

Coverage targets today: `ast`, `transpiler`, `mcp`, `bpe`, and CLI e2e (`fmt` / `build` / `tokens`).

## Website

Static showcase in [`web/`](../web/index.html). Open locally:

```bash
cd web && python3 -m http.server 8080
# → http://localhost:8080
```

Deployable to any static host (GitHub Pages, Netlify, Cloudflare Pages).

## Authorship

Human owner: Sodawave / SodaWave. No AI co-author trailers on commits/PRs.
