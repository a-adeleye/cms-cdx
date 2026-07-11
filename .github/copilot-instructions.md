# Repository Engineering Instructions

Before modifying this repository, read and follow:

1. `/AGENTS.md`
2. `/.ai/ORCHESTRATOR.md` — the mandatory operating procedure for substantial work (task intake, effort selection, delegation, verification, completion gates)
3. `/rules/engineering-rules.md` — rule levels, precedence, exceptions, Definition of Done, and the canonical Rule Files index; read every rule file whose trigger applies
4. `/rules/ai-agents.md` — rules governing AI agent conduct

Apply the relevant engineering rule IDs (e.g. `SEC-10`, `ARCH-15`) to implementation and review.

Preserve existing architecture and behaviour unless the task explicitly requires a change (AGT-01 to AGT-03, MNT-27).

Do not consider work complete until applicable tests pass and the real end-to-end path has been verified where one exists (GEN-13, GEN-15, AGT-10).

Do not commit unrelated changes, secrets, temporary debugging code, or generated artifacts that are not intentionally part of the task (MNT-05, SEC-14).
