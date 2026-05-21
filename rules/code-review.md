# Code Review Rules

Code review is a quality gate for correctness, security, reliability, maintainability, and operational readiness. Review should improve the change, not merely approve it.

## Before Submitting

- Formatting, linting, type checks, builds, and relevant tests MUST be run before requesting review.
- Unused files, debug logs, commented-out code, temporary flags, and dead code MUST be removed.
- Pull requests SHOULD be focused on one main change.
- The description MUST explain what changed, why it changed, how it was tested, and any risks.
- Screenshots, recordings, API examples, migration notes, or operational notes SHOULD be included when they make review clearer.
- Known limitations or follow-up work MUST be disclosed.

## Review Checklist

Reviewers SHOULD check:

- Does the change solve the correct problem?
- Is the solution as simple as the problem allows?
- Is the code readable and maintainable?
- Are edge cases and failure modes handled?
- Are authentication, authorization, privacy, and security handled?
- Are errors handled safely and observably?
- Are tests added or updated at the right level?
- Is there unnecessary complexity or coupling?
- Will this scale reasonably for expected usage?
- Are migrations, rollbacks, and compatibility handled?
- Will another developer understand this later?

## Blocking Issues

Reviewers MUST block changes that:

- Introduce known security vulnerabilities.
- Bypass required authorization or validation.
- Risk data loss, data corruption, or tenant/user data exposure.
- Break public contracts without migration or coordination.
- Add untested critical business logic.
- Add unreliable or non-deterministic required tests.
- Add unexplained large dependencies, frameworks, or architecture changes.
- Make production operation materially worse without mitigation.

## Pull Request Standards

- One pull request SHOULD do one main thing.
- Formatting-only, refactoring-only, dependency-only, and feature changes SHOULD be separate unless tightly related.
- Large changes SHOULD be split into smaller reviewable parts.
- Public API, schema, migration, security, or operational changes MUST be called out explicitly.
- Reviewers SHOULD be able to understand the change without reconstructing context from private conversations.
- If a change is urgent, review standards still apply; risk acceptance must be explicit.

## Reviewer Expectations

- Review comments SHOULD identify the concrete risk or improvement.
- Nitpicks SHOULD be marked clearly and should not block unless they violate agreed standards.
- Reviewers SHOULD distinguish required fixes from optional suggestions.
- Reviewers SHOULD verify tests or evidence for critical paths.
- Approval means the reviewer believes the change is safe enough to merge under the project standards.

## Author Expectations

- Authors MUST respond to review comments or make the requested change.
- Authors SHOULD explain tradeoffs when declining a suggestion.
- Authors MUST not resolve substantive review comments without addressing the issue or explaining why it is not applicable.
- Authors SHOULD keep the pull request updated and easy to re-review after changes.
