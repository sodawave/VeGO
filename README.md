# BMAD workspace template

Clean-slate repository scaffolded for **BMad Method** multi-agent development. Product code and domain names are intentionally absent — clone or copy this tree, rename the project, then invent architecture through BMAD (debate → spine → spec → implement).

## Layout

| Path | Role |
|------|------|
| `src/` | Product source (stubs only until a product is defined) |
| `docs/` | Long-term project knowledge (`bmm.project_knowledge`) |
| `_bmad/` | BMad Method modules (BMM, CIS, TEA) |
| `_bmad-output/` | Planning, implementation, and test artifacts (English) |
| `.agents/skills/` | Cursor / OpenCode / Pi / Antigravity CLI skills |
| `.claude/skills/` | Claude Code skills |
| `.agent/skills/` | Google Antigravity skills |
| `.opencode/commands/` | OpenCode command bindings |
| `AGENTS.md` | Process rules for coding agents |

## Prerequisites

- Node.js 20.12+
- [uv](https://docs.astral.sh/uv/) (required by BMAD Python skills)
- A supported AI coding tool (Cursor, Claude Code, OpenCode, Antigravity, or Pi)

## Get started

1. Open this folder in your AI tool.
2. Invoke the `bmad-help` skill and ask what to do next.
3. Prefer: **debate → architecture spine + memlog → spec or epics/stories → implement → evidence**.

Chat may be in Spanish; formal BMAD documents are English.

Authorship is human-only (Sodawave). After clone, enable the trailer-strip hook:

```bash
git config core.hooksPath .githooks
```

## Refresh BMAD

```bash
npx bmad-method install --directory . --yes
```

Use `--modules` / `--tools` when adding modules or IDE bindings. See [BMAD install docs](https://docs.bmad-method.org/start/install-bmad/).
