# Gemini Repository Instructions

Read and follow `AGENTS.md` before performing repository work.

You MUST also read:

1. `rules/engineering-rules.md` — rule levels, precedence, enforcement gates, the exception process, shared definitions, and the Definition of Done. Its **Rule Files index** is the canonical routing table: read every rule file whose "read when" trigger applies to your task.
2. `rules/ai-agents.md` — rules governing your own conduct as an autonomous agent.
3. `.ai/ORCHESTRATOR.md` — the mandatory operating procedure (SOP) for all substantial work.

Treat these files as authoritative repository instructions.

## Role

- When acting as the **orchestrator**, follow the complete orchestration workflow in `.ai/ORCHESTRATOR.md`, including rule mapping by ID, delegation contracts, rule-compliance review, and the completion rule report.
- When acting as a **delegated specialist**, follow only the assigned scope, the applicable engineering rule IDs, acceptance criteria, and verification requirements. Report conflicts or blockers instead of silently changing scope.
- When assigned **exploratory work**, return evidence, alternatives, risks, edge cases, and recommendations. Do not modify unrelated files or treat exploratory conclusions as verified implementation.

## Mandatory Behavior

- If a rule conflicts with your default behavior, the rule wins. If two rules conflict, apply the precedence order in `engineering-rules.md` (GEN-03); if the conflict survives, stop and escalate (GEN-04) — never self-resolve.
- Prevent violations before writing code rather than fixing them afterwards. If a violation is nevertheless discovered after generation, report it and correct it immediately — never suppress it.
- If a task cannot be completed without violating a MUST rule, stop before writing code. State the rule ID, the conflict, and the available options, then wait for explicit human direction. Never implement the violating approach or silently substitute a different requirement.
- Do not mark a task complete based only on generated code, compilation, or unit tests when the actual user flow can be exercised.

## If a Rule Is Violated

- Do not silently proceed.
- Identify the specific rule by ID.
- Correct the implementation.
- Briefly state which rule applied and what was changed.
