# Maintainability Rules

Read this file when writing or refactoring any application code, or when adding, upgrading, or removing dependencies. Maintainable code is easy to understand, change, test, review, and delete; optimize for the next developer as well as the first delivery.

## Code Structure

- **MNT-01.** Files, modules, classes, functions, and components MUST have one responsibility, describable in one sentence without "and" joining unrelated behaviours; the reviewer decides. A unit that must change for two unrelated kinds of requirement violates this rule.
- **MNT-02.** Functions SHOULD NOT exceed 30 lines or three levels of nesting (project-configurable defaults). Exceeding a limit requires either splitting into named operations or a one-line justification in the PR.
- **MNT-03.** Business and domain logic MUST be testable without framework, UI, transport, or persistence dependencies (see ARCH-10 for the definitive test).
- **MNT-04.** Duplicate logic SHOULD be removed when the same business rule exists in more than one place. Reusable abstractions MUST be created only when they remove real duplication or clarify a stable concept — not for coincidentally similar code.
- **MNT-05.** A change MUST NOT introduce dead code, unused branches, commented-out code, debug logs, or unreferenced TODOs. Dead code, obsolete flags, and outdated comments in files touched by a change MUST be removed or corrected in that change; instances found elsewhere SHOULD be recorded as a tracked issue, not fixed inline. This is the canonical dead-code rule; other files reference it.
- **MNT-06.** A reviewer unfamiliar with a feature SHOULD be able to trace from the feature entry point to each business rule using static code navigation alone (no runtime debugging, no global text search); a failed attempt is grounds for requesting change.

## Naming

- **MNT-07.** Names for domain concepts MUST use the business domain's terminology; technical code MUST use established technical terms. Conventional short names (loop indices, lambda parameters) MAY be used within scopes of a few lines.
- **MNT-08.** Vague names such as `data`, `item`, `temp`, `obj`, `result`, `manager`, or `helper` SHOULD NOT be used unless the local context makes the meaning obvious.
- **MNT-09.** Boolean names SHOULD read as true/false statements: `isActive`, `hasPermission`, `canEdit`, `shouldRetry`.
- **MNT-10.** Public APIs, events, commands, and database fields MUST use consistent terminology.
- **MNT-11.** Names MUST NOT hide side effects: a function named like a query (`get*`, `is*`, `find*`) MUST NOT mutate observable state (internal memoization is exempt).

## Comments and Documentation

- **MNT-12.** Self-explanatory code SHOULD be preferred over comments that restate the implementation; comments SHOULD explain why a non-obvious decision exists.
- **MNT-13.** Comments made outdated by a change MUST be corrected or removed in that change (scope per MNT-05).
- **MNT-14.** Every public module MUST document its purpose and usage where the code lives. Every workflow spanning more than one service or involving manual operational steps MUST have a runbook.
- **MNT-15.** Every deployable application or service MUST have a README covering purpose, ownership, setup, and the run and test commands. Public APIs SHOULD be documented from a source-of-truth schema (e.g. OpenAPI).
- **MNT-16.** Workarounds MUST include the reason, a linked issue when available, and the removal condition.
- **MNT-17.** A TODO MUST reference a tracked issue; TODOs without one are removable on sight.

## Dependencies

This section owns the dependency lifecycle. Supply-chain security controls (scanning, install scripts, pipeline credentials) live in security.md §Dependencies and Supply Chain.

- **MNT-18.** Dependencies SHOULD NOT be added for tasks the language, platform, or an existing dependency already handles well.
- **MNT-19.** A new direct dependency MUST have: its purpose stated in the PR that introduces it (Gate: Evidence); a license on the organization's approved list (absent a list, a mainstream OSI license — flag the missing list per AGT-06; Gate: Review); no known unpatched critical/high vulnerabilities per the CI scanner (SEC-31; Gate: CI); and either evidence of active maintenance (a release or commit within the last 12 months, project-configurable) or a written justification that it is stable and complete (Gate: Evidence).
- **MNT-20.** Every direct dependency MUST have an owner; ownership defaults to the team owning the consuming module. Transitive dependencies are governed via the lockfile and scanner (MNT-21, SEC-31), not individual ownership.
- **MNT-21.** Dependency versions MUST be pinned via a lockfile (or ecosystem equivalent) committed to the repository; lockfile changes MUST be included in review, and builds MUST be reproducible from the committed manifest and lockfile. Gate: CI.
- **MNT-22.** Unused dependencies MUST be removed.
- **MNT-23.** Dependency upgrades SHOULD be small, frequent, and tested.
- **MNT-24.** A dependency that adds a new runtime service, build step, license, network or privileged access, or a major-version framework change MUST be approved by the owning team's senior reviewer or architecture owner, recorded in the PR or an ADR. Gate: Evidence.

## Refactoring

- **MNT-25.** Code that has become hard to understand, test, change, or operate SHOULD be refactored.
- **MNT-26.** Large refactors MUST NOT be mixed with unrelated feature changes in one PR.
- **MNT-27.** Existing behaviour MUST be preserved unless the task explicitly requires changing it; behaviour-preserving refactors SHOULD be covered by tests before the change is made.
- **MNT-28.** A refactor spanning multiple modules, requiring more than one PR, or changing a published interface MUST have a written plan (ADR or tracked issue) covering sequencing, compatibility during transition, and rollback — before the first PR merges. Gate: Evidence.

## Technical Debt

- **MNT-29.** A known-deficient implementation merged deliberately (accepted debt) MUST be registered as debt with an owner and revisit date, using the GEN-05 exception fields. A one-off SHOULD deviation needs only the GEN-06 justification.
- **MNT-30.** Debt items SHOULD be reviewed on a cadence the owning team defines, and debt remediation SHOULD ship as separate small changes, not bundled with features.

## Consistency

- **MNT-31.** Formatting MUST be automated, and formatting, linting, and type checks MUST be enforced in CI (canonical rule: TST-30); additional static analysis SHOULD be enforced in CI. Gate: CI.
- **MNT-32.** Existing project conventions SHOULD be followed; deviations MUST cite a documented reason in the PR. Following a repo convention never justifies violating a MUST in another rule file (GEN-03); apply the MUST and note the divergence.
- **MNT-33.** Introducing a new pattern for a problem class the codebase already solves MUST be called out in the PR and approved by a code owner, with the pattern documented where the project's conventions live. Gate: Evidence.

## Configuration

- **MNT-34.** All required configuration MUST be validated at startup against a typed schema, and every configuration key MUST have a documented purpose, type, and default (in the schema definition or a config reference).
- **MNT-35.** Environment-specific values MUST NOT be hardcoded into application logic.
- **MNT-36.** Missing or invalid required configuration MUST fail fast with a safe, actionable error.
- **MNT-37.** Defaults MUST be safe for production or clearly marked as development-only.
