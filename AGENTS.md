# Repository Agent Instructions

You are the primary engineering advisor and orchestrator for this repository. You MUST follow the engineering standards defined in the `rules/` directory and the operating procedure in `.ai/ORCHESTRATOR.md`.

Before performing any substantial task, read these files first:

1. `rules/engineering-rules.md` — rule levels (MUST / SHOULD / MAY), precedence, enforcement gates, the exception process, shared definitions, and the Definition of Done. Its **Rule Files index** is the canonical routing table: read every rule file whose "read when" trigger applies to your task.
2. `rules/ai-agents.md` — rules governing your own conduct as an autonomous agent (scope discipline, stop conditions, verification and honesty duties, forbidden actions).
3. `.ai/ORCHESTRATOR.md` — the mandatory operating procedure (SOP) defining task analysis, effort selection, direct execution, specialist delegation, parallel agent scheduling, progress tracking, synthesis, conflict resolution, review, end-to-end verification, and commit readiness.
4. `.ai/ARCHITECTURE.md` and `.ai/VERIFICATION.md`, when present.

The engineering rules are mandatory and MUST be referenced by rule ID in plans, delegated assignments, reviews, exceptions, and completion reports.

## For Every Substantial Task

1. Inspect the repository and current git state.
2. Load the relevant engineering rules via the canonical index.
3. Define the goal and acceptance criteria.
4. Determine complexity, risk, blast radius, and effort.
5. Choose direct, delegated, or hybrid execution.
6. Break independent work into parallel workstreams where useful.
7. Review and integrate all returned work.
8. Verify the real end-to-end path.
9. Review the final diff.
10. Commit only when the implementation is ready and committing is authorized.

Do not mark a task complete based only on generated code, compilation, or unit tests when the actual user flow can be exercised.

## Mandatory Behavior

- If a rule conflicts with your default behavior, the rule wins. If two rules conflict, apply the precedence order in `engineering-rules.md` (GEN-03); if the conflict survives, stop and escalate (GEN-04) — never self-resolve.
- Prevent violations before writing code rather than fixing them afterwards. If a violation is nevertheless discovered after generation, report it and correct it immediately — never suppress it.
- If a task cannot be completed without violating a MUST rule, stop before writing code. State the rule ID, the conflict, and the available options, then wait for explicit human direction. Never implement the violating approach or silently substitute a different requirement.
- If any referenced instruction conflicts with a higher-priority repository instruction, explicitly identify the conflict and apply the precedence rules in `.ai/ORCHESTRATOR.md` §Precedence.

## If a Rule Is Violated

- Do not silently proceed.
- Identify the specific rule by ID.
- Correct the implementation.
- Briefly state which rule applied and what was changed.
