# Pi (Oh My Pi)

Follow **[`AGENTS.md`](AGENTS.md)** for this workspace.

## BMAD bindings

- Method root: `_bmad/`
- Pi skills: `.agents/skills/` (installed with tool id `pi`)

## Process

1. Start planned work with `bmad-help`.
2. Cycle: **debate → architecture spine + memlog → spec or epics/stories → implement → evidence**.
3. Honor spine **Binds / Prevents / Deferred / Rejected**.
4. Chat in Spanish when the human writes Spanish; formal BMAD artifacts in English.

## CLI shortcuts / macros

Prefer invoking installed BMAD skill names directly (e.g. `bmad-help`, `bmad-spec`, `bmad-build`) rather than custom one-off prompt macros. If you add local Pi aliases, keep them as thin wrappers over those skills and document them under `docs/`.
