# Claude Code

Follow **[`AGENTS.md`](AGENTS.md)** for this workspace.

## Process (generic)

1. **BMAD required** at repo root (`_bmad/`, `.claude/skills/` and `.agents/skills/`). Start planned work with `bmad-help`. Never install BMAD under product source trees.
2. Cycle: **debate → architecture spine + memlog → spec or epics/stories → implement → evidence**.
3. Before inventing architecture: read `_bmad-output/planning-artifacts/**/ARCHITECTURE-SPINE.md` and `.memlog.md`. Obey **Binds / Prevents / Deferred / Rejected**. Do not reopen closed paths without a new adopted spine decision.
4. Chat in **Spanish** when the human writes Spanish; BMAD formal artifacts in **English**.

## Skills

BMAD skills live under `.claude/skills/` (installed via `npx bmad-method install`). Prefer those over ad-hoc process inventions.

## Authorship (canonical)

Follow **Authorship (canonical)** in [`AGENTS.md`](AGENTS.md). Never add Cursor/`cursoragent` co-author or Made-with trailers. Project + human own conception and architecture.
