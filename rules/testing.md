# Testing Rules

Read this file when writing or reviewing tests of any kind. Tests must prove important behaviour, prevent regressions, and support safe change; tests are part of the product, not disposable scaffolding.

**Critical code** and **risky change** are defined in engineering-rules.md §Definitions. Use those terms; this file does not redefine them.

## General

- **TST-01.** Critical code MUST be tested, covering normal cases, edge cases, and failure cases.
- **TST-02.** Risky changes MUST include tests at every level the change affects (unit, integration, contract, UI, end-to-end).
- **TST-03.** Bug fixes MUST include a regression test unless the defect is not practically testable; the exception is documented in the PR description with the reason. The regression test MUST fail on the pre-fix code and pass with the fix; state this verification in the PR. Gate: Evidence.
- **TST-04.** Tests MUST NOT access private members, private state, or internal call sequences of the unit under test. Outbound calls to boundary dependencies (e.g. verifying a payment API was not retried) count as observable behaviour.

## Structure and Naming

- **TST-05.** Tests SHOULD follow Arrange-Act-Assert or an equivalent clear structure, with setup explicit enough for the reader to understand the scenario.
- **TST-06.** One test SHOULD verify one behaviour, and a test MUST NOT assert multiple unrelated behaviours (behaviours that can fail independently) in a single case.
- **TST-07.** Shared helpers MUST NOT contain assertions invisible at the call site; any value a test asserts on MUST be visible in the test body or in the helper's name or arguments.
- **TST-08.** Test names MUST describe the scenario and expected outcome.

Bad:

```text
test_user()
it works
testCase1
```

Good:

```text
returns 401 when token is missing
calculates total with discount applied
does not retry non-idempotent payment requests
```

## Unit Tests

- **TST-09.** Unit tests SHOULD cover pure logic, utilities, validators, calculations, state transitions, and domain rules.
- **TST-10.** Unit tests MUST NOT perform network I/O or touch any real database, external service, or the file system. In-process substitutes (in-memory fakes, embedded in-memory databases) are permitted; file-system utilities MAY use isolated temporary storage.
- **TST-11.** Unit tests SHOULD complete fast enough to run on every local build.

## Integration Tests

- **TST-12.** Integration tests MUST cover important flows across modules, services, persistence, queues, or framework boundaries for critical code.
- **TST-13.** Real implementations SHOULD be used unless the dependency is external, slow, unstable, costly, or unsafe; integration tests SHOULD mock only at system boundaries (external APIs, third-party services, email, payments, identity providers, infrastructure that cannot run safely in tests). This is the canonical mock-boundary rule for all test levels.
- **TST-14.** Database integration tests MUST use isolated test data.

## Contract and API Tests

- **TST-15.** Externally consumed APIs, events, webhooks, and shared schemas MUST have an automated compatibility check in CI (contract test, schema diff, or consumer-driven contract) that fails on breaking changes; interfaces lacking one require an exception per GEN-05. Gate: CI.
- **TST-16.** Consumer-impacting changes SHOULD include compatibility tests or migration evidence.

## UI and End-to-End Tests

- **TST-17.** Important user journeys through critical code SHOULD be tested end to end, covering loading, empty, success, validation, authorization, and error states — not every visual detail.
- **TST-18.** Snapshot updates MUST be reviewed as code changes, not bulk-approved, and dynamic values (timestamps, IDs, generated class names) MUST be normalized or masked before snapshotting. Snapshots SHOULD target a specific component or element; full-page snapshots SHOULD NOT be required gates unless the tooling provides reviewable per-change visual diffs.

## Test Doubles

- **TST-19.** Real implementations SHOULD be preferred over mocks; test doubles SHOULD be used only when they make tests safer, faster, clearer, or more deterministic.
- **TST-20.** Internal business logic MUST NOT be mocked to make a failing test pass.
- **TST-21.** Time-dependent tests MUST use fake clocks or fixed timestamps unless the test verifies real timing behaviour, and randomness MUST be controlled with fixed seeds or deterministic inputs.

## Test Data

- **TST-22.** Test data MUST be isolated from production data, and tests SHOULD create only the data they need.
- **TST-23.** Tests MUST clean up after themselves when they use shared infrastructure.
- **TST-24.** Sensitive production data MUST NOT be copied into test fixtures unless approved by the data owner or security team via an exception per GEN-05, minimized, and protected to the same standard as production.

## Test Quality

- **TST-25.** Tests MUST be deterministic: the same test produces the same result for the same code. If a test cannot be made deterministic, the reason and risk MUST be documented.
- **TST-26.** Tests in required gates MUST NOT depend on execution order or on state left by other tests, and MUST pass when run in isolation.
- **TST-27.** A flaky, skipped, or disabled test in a required gate MUST reference a tracking issue with an owner and expiry date per GEN-05, and MUST be fixed or deleted by that expiry. A required gate that fails on unchanged code is flaky under this rule.
- **TST-28.** Assertions MUST report expected and actual values (value-comparing matchers, not bare boolean assertions, where the framework supports them); the failure message SHOULD identify the violated expectation without a debugger, the reviewer deciding from the failure output.
- **TST-29.** Test suites SHOULD be organized so fast feedback runs early and deeper checks run in CI or pre-release pipelines.

## CI Quality Gates

- **TST-30.** Formatting, linting, type checks, builds, and the relevant automated tests MUST run in CI and MUST pass before merge to the default branch. This is the canonical CI-gate rule. Gate: CI.
- **TST-31.** Required test gates MUST be reliable enough to block merge (a gate that fails on unchanged code falls under TST-27).
- **TST-32.** Coverage thresholds, when used, MUST be enforced by CI on every pull request; manual overrides require an exception per GEN-05. Tests without assertions MUST NOT be merged, and reviewers MUST reject tests whose only effect is executing lines to satisfy a threshold.
- **TST-33.** Critical code SHOULD have meaningful behavioural coverage even when global coverage numbers are high.
