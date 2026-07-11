# Engineering Rules Enforcement

You MUST follow the engineering standards defined in the `rules/` directory and the operating procedure in `.ai/ORCHESTRATOR.md`.

Before writing code, reviewing, or preparing a change, read these files first:

1. `rules/engineering-rules.md` — rule levels (MUST / SHOULD / MAY), precedence, enforcement gates, the exception process, shared definitions, and the Definition of Done. Its **Rule Files index** is the canonical routing table: read every rule file whose "read when" trigger applies to your task.
2. `rules/ai-agents.md` — rules governing your own conduct as an autonomous agent (scope discipline, stop conditions, verification and honesty duties, forbidden actions).
3. `.ai/ORCHESTRATOR.md` — the mandatory operating procedure (SOP) for all substantial work: task intake and analysis, effort selection, direct vs delegated vs hybrid execution, specialist selection, parallel agent scheduling, delegation contracts, progress tracking, review, conflict resolution, synthesis, real end-to-end verification, quality gates, and commit readiness.

All code changes and implementation work MUST satisfy the Definition of Done in `rules/engineering-rules.md` before being considered complete.

## Role

- When acting as the **orchestrator** (owning a goal end to end), follow the complete orchestration workflow in `.ai/ORCHESTRATOR.md`, including rule mapping by ID, delegation contracts with an Engineering Rules section, rule-compliance review, and the completion rule report.
- When acting as a **delegated specialist**, follow only the assigned scope, the applicable engineering rule IDs, acceptance criteria, and verification requirements from your delegation contract. Report conflicts or blockers instead of silently changing scope.

## Mandatory Behavior

- Read the relevant rule files **before** writing code, reviewing, or preparing a change. When in doubt, read the file. The rules take precedence over your defaults.
- If a rule conflicts with your default behavior, the rule wins. If two rules conflict, apply the precedence order in `engineering-rules.md` (GEN-03); if the conflict survives, stop and escalate (GEN-04) — never self-resolve.
- Prevent violations before writing code rather than fixing them afterwards. If a violation is nevertheless discovered after generation, report it and correct it immediately — never suppress it.
- If a task cannot be completed without violating a MUST rule, stop before writing code. State the rule ID, the conflict, and the available options, then wait for explicit human direction. Never implement the violating approach or silently substitute a different requirement.
- Do not mark a task complete based only on generated code, compilation, or unit tests when the actual user flow can be exercised — verify the real end-to-end path per `.ai/ORCHESTRATOR.md` §Real End-to-End Verification.

## If a Rule Is Violated

- Do not silently proceed.
- Identify the specific rule by ID.
- Correct the implementation.
- Briefly state which rule applied and what was changed.
