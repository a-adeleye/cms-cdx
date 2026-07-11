# Engineering Rules

The minimum engineering standard for production systems in this organization. It governs code written by humans and by autonomous AI coding agents. Read this file first, always; then read every rule file whose trigger applies (see the Rule Files index at the end).

## Rule Levels

Every rule carries exactly one uppercase level keyword and a stable ID (e.g. `SEC-12`). Lowercase modals ("must", "should") and statements without a level keyword are non-normative context. Cite the rule ID when reporting a violation, requesting an exception, or blocking a review.

- **MUST / MUST NOT**: Required. A violation blocks merge or release unless an approved exception exists (see Exceptions).
- **SHOULD / SHOULD NOT**: Required by default. Deviations require a written justification in the pull request or design record.
- **MAY**: Optional guidance.

## Scope

- **GEN-01.** These rules MUST be applied to all code merged into a production repository, including application code, tests, CI pipelines, infrastructure code, and generated code.
- **GEN-02.** Prototypes and throwaway scripts MAY be exempted from all rules except those in security.md and reliability.md, and only while they cannot reach production data, systems, or users.

## Precedence

- **GEN-03.** When rules appear to conflict, resolve in this order: (1) MUST outranks SHOULD outranks MAY; (2) at equal level, the more specific rule file (e.g. angular.md) overrides the more general one within its scope; (3) a repo-local convention MAY override a SHOULD when the deviation is noted in the pull request, and MUST NOT override a MUST.
- **GEN-04.** A conflict that survives GEN-03 is a defect in this rule set. The author or agent MUST stop, state both rules and the conflict, and obtain a human decision recorded via the exception process. Autonomous agents MUST NOT self-resolve rule conflicts.

## Enforcement

Every MUST maps to a named enforcement gate; a rule bundling several clauses names one gate per clause. Where a rule does not name its gate, the gate is Review.

- **CI**: a named automated check (formatter, linter, compiler, boundary rule, scanner, test gate) fails the build on violation.
- **Review**: the reviewer verifies the rule's stated objective criterion from the diff; approval asserts it passed.
- **Evidence**: the pull request or release record attaches a named artifact (test output, ADR, rollback note, migration plan); the reviewer verifies the artifact exists and addresses the rule.

## Exceptions

- **GEN-05.** An exception to a MUST rule requires a written record with five fields — reason, risk, owner, expiry date, mitigation — approved by the code owner or tech lead (for urgent changes, the on-call lead MAY approve), linked from the pull request, and stored in the exceptions log (default `docs/exceptions/`; a project MAY relocate it if the README says where).
- **GEN-06.** A deviation from a SHOULD rule requires only a justification in the pull request or design record.
- **GEN-07.** Autonomous agents MUST NOT approve exceptions. An agent MAY draft the record (reason, risk, mitigation) with owner and expiry marked "requires human assignment", and MUST NOT merge or declare the task complete until a human approves.
- **GEN-08.** Waiver language elsewhere in this rule set ("risk accepted", "explicitly approved", "documented reason") MUST be interpreted as an exception under GEN-05 or GEN-06. There is no third process.

## Definitions

These terms are used across all rule files.

- **Module**: a declared unit in the repository's module artifact (workspace/package config, Nx project graph, or a MODULES.md registry) with a name, owner, purpose, and public entry point. See architecture.md.
- **Critical code**: code handling money or payments, authentication, authorization or permissions, tenant isolation, or irreversible data changes — plus paths the owning team designates in a repo-recorded list. The code owner decides borderline cases in review.
- **Critical workflow / critical dependency**: one whose failure causes data loss, financial impact, security exposure, or user-facing outage; the service owner records the designation.
- **Risky change**: a change that touches critical code, migrates or transforms data, or alters concurrency, caching, or security behaviour.
- **Sensitive data**: credentials, tokens, personal data, payment data, and data covered by regulation.

## Core Principle (non-normative)

Build simple, secure, reliable, maintainable systems that can be operated in production. Prefer proven technology over unnecessary complexity — the enforceable form of this preference is the ADR rule (GEN-11) and the dependency rules in maintainability.md.

## Enterprise Baseline

- **GEN-09.** Every application MUST have: an owning team resolvable from CODEOWNERS or the service registry; a purpose statement in its README; and a documented support or escalation contact. Gate: Evidence.
- **GEN-10.** Every production change MUST: merge via a reviewed pull request (code-review.md §Merge Requirements); pass the required CI checks (testing.md §CI Quality Gates); ship through the standard deployment pipeline (delivery.md); emit the telemetry baseline (reliability.md §Observability); and have a rollback or remediation plan (delivery.md §Releases).
- **GEN-11.** A decision SHOULD be recorded as a lightweight ADR (default location `docs/adr/`) when it is expensive to reverse, crosses module or service boundaries, introduces a new technology or dependency category, or changes a public contract. A reviewer MAY require an ADR for other decisions.
- **GEN-12.** Generated code, prototypes, scripts, and internal tools MUST follow security.md and reliability.md whenever they touch production data or systems.

## Definition of Done

Every code change MUST satisfy every item below before merge.

- **GEN-13.** The intended behaviour MUST be demonstrated by passing automated tests or, where automation is impractical, by manual verification steps documented in the pull request. Gate: Evidence.
- **GEN-14.** Edge cases and failure paths relevant to the change MUST be handled, with errors following reliability.md §Error Handling.
- **GEN-15.** All required checks — formatting, lint, type checks, build, and the test suites of every changed module — MUST have been executed and passed, with the author reporting what was run and the result. Declaring a change done without executed evidence is itself a MUST violation. Gate: Evidence.
- **GEN-16.** Tests MUST be added or updated wherever the change affects behaviour (testing.md).
- **GEN-17.** The change MUST NOT introduce dead code or debug artifacts (canonical list: MNT-05 in maintainability.md).
- **GEN-18.** Naming MUST be clear and consistent with the domain (maintainability.md §Naming).
- **GEN-19.** For changes touching authentication, authorization, input handling, personal data, schema, or deployment, the pull request MUST state the security/privacy impact and the operational impact (monitoring, migration, rollback), or state "none" with a reason. Gate: Evidence.
- **GEN-20.** Documentation invalidated by the change (README, module docs, runbooks) MUST be updated in the same change.
- **GEN-21.** The pull request MUST meet the pre-submission checklist in code-review.md (one logical change; description states what, why, how it was tested, and risks).

## Release Readiness

Detailed release, deployment, and migration rules live in delivery.md; migration mechanics live in reliability.md §Data Integrity; feature-flag lifecycle lives in reliability.md §Feature Flags. Non-negotiables:

- **GEN-22.** Production releases MUST have a documented rollback plan; where rollback is technically impossible (e.g. a destructive migration), a written remediation plan approved by the owning team before release. An unwritten "we will fix forward" does not satisfy this rule. Gate: Evidence.
- **GEN-23.** Breaking API, schema, or behaviour changes MUST follow architecture.md §Contracts (documented, communicated to consumers, and versioned or safely migrated).
- **GEN-24.** Known risks MUST be listed in the pull request or release record before release, not discovered from production incidents. Gate: Evidence.

## Rule Set Governance

- **GEN-25.** Changes to `rules/` MUST be approved by the standards owner (a named role or group the organization designates once).
- **GEN-26.** Adding, renaming, or removing a rule file or the operating procedure MUST update the Rule Files index below and the agent instruction files (`.claude/CLAUDE.md`, `AGENTS.md`, `.gemini/GEMINI.md`, `.github/copilot-instructions.md`) in the same change.
- **GEN-27.** Rule IDs MUST be stable: IDs of deleted rules are retired, never reused, and rules are never renumbered.

## Rule Files (canonical index)

This table is the single source of truth for rule-file routing. The agent instruction files reference it and must not carry a diverging copy.

| File                 | Read when                                                                                       |
| -------------------- | ----------------------------------------------------------------------------------------------- |
| `engineering-rules.md` | Always — levels, precedence, enforcement, exceptions, definitions, Definition of Done         |
| `ai-agents.md`       | Always, when the work is performed by an autonomous coding agent                                 |
| `../.ai/ORCHESTRATOR.md` | Always, before orchestrating or performing substantial work — the operating procedure (SOP): task intake, effort selection, delegation, verification, completion gates (AGT-16) |
| `architecture.md`    | Creating or modifying modules, services, layers, folder structure, or designing/changing APIs, events, or webhooks |
| `angular.md`         | Writing any Angular component, service, directive, pipe, or template                             |
| `maintainability.md` | Writing or refactoring any application code, or adding/upgrading/removing dependencies          |
| `reliability.md`     | Handling errors, external services, validation, state, data integrity, migrations, or feature flags |
| `scalability.md`     | Writing database queries, caching, background jobs, or performance-sensitive code               |
| `security.md`        | Handling authentication, authorization, secrets, input, user data, or reviewing supply-chain risk |
| `testing.md`         | Writing or reviewing tests of any kind                                                          |
| `code-review.md`     | Preparing, committing, or reviewing a change; branching and merge mechanics                     |
| `delivery.md`        | Changing CI/CD pipelines, build or deploy configuration, environments, infrastructure, or preparing a release |
