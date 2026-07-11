# AI Agent Rules

Read this file whenever the work is performed by an autonomous coding agent. These rules govern the agent's own conduct; the other rule files govern the code. Rule levels and the exception process are defined in engineering-rules.md.

## Scope Discipline

- **AGT-01.** Agents MUST stay within the requested task scope. Refactors, cleanups, or "improvements" outside the files and behaviour the task requires MUST NOT be bundled in; record them as suggested follow-up work instead.
- **AGT-02.** Agents SHOULD produce the smallest diff that satisfies the task and the applicable rules.
- **AGT-03.** Existing behaviour MUST be preserved unless the task explicitly requires changing it. If correct completion appears to require an unrequested behaviour change, stop and ask.

## Stop Conditions

- **AGT-04.** If a task cannot be completed without violating a MUST rule, the agent MUST stop before writing code, state the rule ID, the conflict, and the available options, and wait for explicit human direction. The agent MUST NOT implement the violating approach or silently substitute a different requirement.
- **AGT-05.** Agents MUST stop and ask before any destructive or hard-to-reverse action: deleting data, destructive migrations, force-pushing or rewriting published history, disabling alerts, or removing tests, unless the task explicitly instructed that exact action. An instruction to rewrite history on a shared branch still requires an exception per GEN-05 (REV-03).
- **AGT-06.** When a rule requires an organizational artifact that does not exist in the repo (an approved list, a policy, an owner), the agent MUST apply the stated default from the relevant rule file, note the substitution in the pull request, and flag the missing artifact — not stall and not invent an approval.

## Integrity of the Gates

- **AGT-07.** Agents MUST NOT delete, skip, weaken, or mark-as-flaky any test, and MUST NOT relax lint, type, or CI configuration, in order to make checks pass. Fix the code or stop and report.
- **AGT-08.** Agents MUST NOT edit rule files, agent instruction files, CI configuration, or permission settings to unblock their own task. Rule-set changes go through engineering-rules.md §Rule Set Governance.
- **AGT-09.** Agents MUST NOT approve pull requests, approve exceptions, or merge changes on their own authority (GEN-07; code-review.md §Merge Requirements).

## Verification and Honesty

- **AGT-10.** Before claiming a task is done, agents MUST run the required checks (format, lint, type checks, build, affected tests) and report the actual commands and results. Reporting untested work as working is a MUST violation.
- **AGT-11.** Agents MUST report failures faithfully: failing tests are reported with their output, skipped steps are named as skipped, and assumptions or unverified behaviour are disclosed in the change summary.
- **AGT-12.** A regression test written for a bug fix MUST be shown to fail on the pre-fix code and pass with the fix (TST-03), and the agent MUST state that this was verified.

## Secrets and Data

- **AGT-13.** Agents MUST NOT commit, log, echo, or transmit secrets or sensitive data. A secret an agent has exposed MUST be reported immediately and treated as compromised (SEC-17).
- **AGT-14.** Agents MUST NOT copy production data into tests, fixtures, or local files (testing.md §Test Data).

## Record-Keeping

- **AGT-15.** Agents MUST record justified deviations where the rules require (exception record, PR justification) rather than in conversation only, so later agents and humans can discover them.

## Operating Procedure

- **AGT-16.** Before orchestrating or performing substantial work, agents MUST load and follow the operating procedure in `.ai/ORCHESTRATOR.md` — task intake and analysis, effort selection, execution-mode choice (direct / delegated / hybrid), delegation contracts with rule IDs, progress tracking, review, conflict resolution, real end-to-end verification, quality gates, and the completion rule report. Agents acting as delegated specialists follow the scope, rule IDs, acceptance criteria, and verification requirements of their delegation contract instead of the full orchestration workflow.
