# Engineering Rules

These are the minimum engineering standards for enterprise-grade applications. Each topic is covered in a dedicated rule file.

## Rule Files

- [architecture.md](architecture.md) - module design, layers, dependency rules, boundaries, contracts
- [maintainability.md](maintainability.md) - code structure, naming, comments, dependencies, refactoring
- [reliability.md](reliability.md) - error handling, validation, state, external services, observability, recovery
- [scalability.md](scalability.md) - database, performance, background work, capacity, growth
- [security.md](security.md) - authentication, authorization, secrets, input safety, data protection, abuse prevention
- [testing.md](testing.md) - unit, integration, contract, UI tests, test doubles, test quality
- [code-review.md](code-review.md) - pre-submission checklist, review standards, pull request rules

## Rule Levels

- MUST: Required for production code. A violation blocks merge or release unless an approved exception exists.
- SHOULD: Required by default. Exceptions must be justified in the pull request or design record.
- MAY: Optional guidance that can be used when appropriate.

## Core Principle

Build simple, secure, reliable, maintainable systems that can be operated in production. Prefer proven technology over unnecessary complexity. Security, reliability, scalability, and maintainability are not optional afterthoughts.

## Enterprise Baseline

- Every application MUST have clear ownership, a documented purpose, and a support path.
- Every production change MUST be reviewable, testable, deployable, observable, and reversible.
- Critical decisions SHOULD be recorded in lightweight architecture decision records.
- Exceptions to these rules MUST include the reason, risk, owner, expiry date, and mitigation.
- Generated code, prototypes, scripts, and internal tools MUST still follow security and reliability requirements when they touch production data or systems.

## Definition of Done

A task is not complete unless:

- The code works as expected for the intended behaviour.
- Edge cases and failure paths are handled.
- Errors are handled safely and logged with enough context to debug.
- Tests are added or updated where the change affects behaviour.
- Formatting, linting, type checks, builds, and relevant tests pass.
- No unused code, dead paths, debug logs, or unnecessary dependencies were introduced.
- Naming is clear and consistent with the domain.
- Security, authorization, privacy, and data handling impacts were considered.
- Operational impact is understood, including monitoring, migration, rollback, and support needs.
- The change can be reviewed without confusion.

## Release Readiness

- Production releases MUST have a rollback or remediation path.
- Database migrations MUST be safe for existing data and active deployments.
- Feature flags SHOULD be used for risky or phased changes.
- Breaking API, schema, or behaviour changes MUST be documented and coordinated with consumers.
- Known risks MUST be captured before release, not discovered from production incidents.
