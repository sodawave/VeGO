# Agent instructions — VeGo (Vector Go)

## What this repo is

- **Product:** VeGo — BPE-optimized, bi-directional, lossless AST serialization over standard Go (`go/ast`). See `docs/OVERVIEW.md` and `docs/ADR.md`.
- **Method:** BMad Method at the workspace root (`_bmad/`, skills under `.agents/skills/`, `.claude/skills/`, `.agent/skills/` including `bmad-help`).
- BMAD is the **principal development process**. Do not invent a parallel methodology.

## Branches

Default flow: **work branch → test with evidence → human OK → merge to `dev`**. Do not develop directly on `dev` or `main` unless the human explicitly agrees (e.g. trivial hotfix).

| Branch | Role |
|--------|------|
| `feat/*`, `fix/*`, `chore/*`, `docs/*`, `spec/*`, … | Day-to-day work — one branch per unit of change |
| `dev` | Integration only — receives merges after develop → test → OK |
| `main` | Production / deploy — merge from `dev` only when the human explicitly asks |

- Prefer branch names aligned with conventional commits (`feat/…` + `feat: …`, `fix/…` + `fix: …`, etc.).
- When work is driven by a formal spec, a `spec/YYYY-MM-DD-<slug>` branch (or a `feat/…` linked to that spec) is the default option — not the only branch type.
- Do not merge to `main`, release, or deploy without an explicit human request.

## Language

- Chat with the human in **Spanish** when they write in Spanish.
- Write BMAD formal artifacts (PRD, epics, stories, architecture, specs) in **English**.

## Working tree

- Implement product code under `src/` (`src/pkg/…`, `src/cmd/vego` — documented in the managed block below).
- Keep methodology files at repo root (`_bmad/`, `_bmad-output/`, `docs/`, `AGENTS.md`).
- Do not install or reinstall BMAD inside product source trees.

## Quality

- Follow spec → plan → implement → verify with evidence.
- Prefer tests for behavior changes; otherwise run an executable verification checklist.
- Double-confirm large/irreversible steps with the human.
- Never claim success, existence, or test results without verifying in the current session.

## Process

- Prefer BMAD Cursor skills (start with `bmad-help`, then agents and workflows) for product work.
- Superpowers and other local skills are optional helpers; they do not replace BMAD for planned product delivery.

## BMAD process (required, generic)

- **BMAD must be installed** at the workspace root (`_bmad/`, skills under `.agent/skills/`). If missing, install/adopt via BMAD (`bmad-help` / `bmad-project-context`). Do **not** invent a parallel methodology; do **not** install BMAD under product source trees.
- Planned product work follows: **debate → architecture spine + memlog → spec or epics/stories → implement → verify with evidence**. Do not skip to large architectural changes without that record.
- Before inventing or changing architecture, read the relevant `_bmad-output/planning-artifacts/**/ARCHITECTURE-SPINE.md` and `.memlog.md`.
- Honor each decision’s **Binds / Prevents / Rule** and any **Deferred / Rejected** list. Those negatives exist so agents do **not** reopen closed design paths on their own. Overturning them requires a new adopted spine decision — not a drive-by refactor.

<!-- bmad:context -->
<!-- Bootstrap from docs/OVERVIEW.md + docs/ADR.md. Refresh via bmad-project-context after first architecture spine. -->

## project (VeGo)

VeGo (Vector Go): source-to-source transpiler and MCP surface for AI-agent-optimized Go AST serialization. Stack: Go (`go/ast`, planned Participle/Cobra/mcp-go/tiktoken-go per ADR). Knowledge in `docs/`; cycle artifacts in `_bmad-output/`.

## Policy

- Never invent product behavior that contradicts `src/` — verify before asserting.
- Never fork `cmd/compile`; VeGo is a transpile layer only.
- Never invent keyword↔Unicode maps without BPE/tiktoken validation.
- Never put BMAD method files under product source trees (`src/`).
- Never commit secrets (`.env`, signing keys); honor `.gitignore` and `.env.example`.

## Authorship (canonical)

- **Owner:** Sodawave / SodaWave. Conception, architecture, and authorship of this project are the human's — not Cursor, not any model/agent.
- **Forbidden on commits and PRs:** `Co-authored-by: Cursor`, `cursoragent@cursor.com`, `Made-with: Cursor`, or any AI/tool authorship trailer or agent author/committer identity.
- **Project overrides:** `.cursor/cli.json` keeps `attribution.attributeCommitsToAgent` and `attributePRsToAgent` at `false`. Do not flip them on.
- **Hook:** after clone, run `git config core.hooksPath .githooks` so `.githooks/prepare-commit-msg` strips those trailers if a harness re-injects them.
- Agents must never claim co-authorship, ownership, or license of this work.

## Where things are

- Product root: `src/` — stubs in `src/pkg/ast`, `src/pkg/transpiler`, `src/pkg/mcp`, `src/cmd/vego`
- Project knowledge: `docs/` (`OVERVIEW.md`, `ADR.md`)
- Planning / implementation / test artifacts: `_bmad-output/`
- BMAD skills: `.agents/skills/`, `.claude/skills/`, `.agent/skills/`, `.opencode/commands/`

## Running and verifying

```bash
go build ./...
go run ./src/cmd/vego
```

Document additional install/run/test commands here after the first non-stub implementation lands.

<!-- /bmad:context -->
