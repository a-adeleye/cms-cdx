# Code Review, Commits, and Merge Rules

Read this file when preparing, committing, or reviewing a change, or when branching and merging. Code review is a quality gate for correctness, security, reliability, maintainability, and operational readiness; review should improve the change, not merely approve it.

## Commits and Branches

- **REV-01.** Commit messages MUST state what changed and why, in imperative form.
- **REV-02.** Each commit SHOULD build and pass tests independently, and branches SHOULD be short-lived and follow the project's naming convention.
- **REV-03.** Force-pushing or rewriting history on shared branches MUST NOT occur.
- **REV-04.** A secret ever committed MUST be rotated and purged per SEC-17 — rotation first, history cleanup second.

## Before Submitting

- **REV-05.** Formatting, linting, type checks, and builds MUST pass before requesting review, and the test suites of every changed module MUST pass. A failing required check means the change is not ready for review.
- **REV-06.** The change MUST contain no dead code or debug artifacts per MNT-05, and no debug-only or developer-only toggles. Release feature flags for phased rollout are exempt but MUST follow REL-30 (owner and removal condition).
- **REV-07.** One pull request SHOULD do one main thing; formatting-only, refactoring-only, dependency-only, and feature changes SHOULD be separate unless tightly related — "tightly related" meaning the parts cannot merge independently without breaking build or behaviour.
- **REV-08.** Pull requests over roughly 400 changed lines (project-configurable; excluding generated files and lockfiles) SHOULD be split, and the reviewer MAY require a split.
- **REV-09.** The description MUST explain what changed, why, how it was tested, and any risks; known limitations and follow-up work MUST be disclosed.
- **REV-10.** Supporting artifacts SHOULD be included when they make review clearer — in particular: before/after screenshots or recordings for user-visible UI changes, example requests and responses for API changes, and migration and rollback notes for schema or migration changes.
- **REV-11.** Public API, schema, migration, security, and operational changes MUST be called out explicitly in the description, and the change MUST be understandable without reconstructing context from private conversations.

## Review Checklist (reviewer attention prompts, SHOULD-level)

Reviewers SHOULD check: Does the change solve the correct problem, as simply as the problem allows? Is it readable and maintainable? Are edge cases and failure modes handled? Are authentication, authorization, privacy, and security handled? Are errors handled safely and observably? Are tests added at the right level? Is there unnecessary complexity or coupling? Will it scale for the stated capacity assumptions? Are migrations, rollbacks, and compatibility handled? Will another developer understand this later?

## Blocking Issues

- **REV-12.** Reviewers MUST block a change that introduces a known security vulnerability; the blocking comment MUST identify the weakness and cite its source — a scanner finding, a CVE, a violated security.md rule, or a concretely described exploit scenario.
- **REV-13.** Reviewers MUST block a change that bypasses required authorization or validation, or that risks data loss, corruption, or tenant/user data exposure.
- **REV-14.** Reviewers MUST block a change that breaks a public contract without following ARCH-15.
- **REV-15.** Reviewers MUST block a change that adds untested critical code (definition in engineering-rules.md §Definitions).
- **REV-16.** Reviewers MUST block a change that adds required tests violating the determinism rules in testing.md (TST-25 to TST-27), citing the violated rule.
- **REV-17.** Reviewers MUST block a change that adds dependencies, frameworks, or architecture changes without the justification required by MNT-19/MNT-24 or the ADR required by ARCH-27.
- **REV-18.** Reviewers MUST block a change that makes production operation worse without mitigation — including removing or degrading monitoring or alerting, breaking rollback or deploy safety, adding an unmonitored failure mode, or comparable degradations; the blocking comment states the specific degradation and the author provides mitigation evidence in the PR.

## Comment Taxonomy

- **REV-19.** Review comments use three labels: **blocking** (violates a MUST rule or a Blocking Issue — must be resolved before merge, rule ID cited), **suggestion** (SHOULD-level or judgment — the author MAY decline with a written rationale), and **nit** (style or preference not backed by any rule — MUST NOT block).
- **REV-20.** Review comments SHOULD identify the concrete risk or improvement, and reviewers SHOULD verify tests or evidence for critical code rather than taking the description's word.

## Author Expectations

- **REV-21.** Authors MUST respond to every blocking comment by making the change or documenting why it does not apply, and MUST NOT resolve substantive comments without doing one of the two.
- **REV-22.** Authors SHOULD keep the pull request easy to re-review after changes (append commits during review rather than rewriting).
- **REV-23.** If author and reviewer cannot agree after one round of discussion, either party MUST escalate to the code owner or tech lead, whose decision is recorded in the pull request.

## Merge Requirements

- **REV-24.** At least one approval from a qualified reviewer who is not the author MUST be obtained; authors MUST NOT approve their own pull requests, and autonomous agents MUST NOT approve or merge changes (AGT-09).
- **REV-25.** All required CI checks MUST pass at merge time (TST-30); changes MUST NOT merge on a red pipeline.
- **REV-26.** Changes to designated sensitive areas (authentication, payments, migrations, secrets — the org-configured CODEOWNERS list) MUST be approved by the code owner.
- **REV-27.** Substantive commits pushed after approval invalidate that approval and MUST be re-reviewed; trivial fixups (typos, comments, rebases with no content change) MAY merge on the existing approval. Enable dismiss-stale-approval enforcement where the platform supports it. Gate: CI.
- **REV-28.** Urgent changes follow the same standards; risk acceptance MUST be recorded in the pull request using the exception format from GEN-05, accepted by the code owner or on-call lead. Gate: Evidence.
