# Testing Rules

Tests must prove important behaviour, prevent regressions, and support safe change. Tests are part of the product, not disposable scaffolding.

## General

- Critical business logic MUST be tested.
- Bug fixes MUST include a regression test unless the defect is not practically testable; the exception must be documented.
- Tests SHOULD be clear, reliable, fast enough to run regularly, and easy to diagnose when they fail.
- Tests MUST focus on observable behaviour instead of private implementation details.
- Risky changes MUST have proportionate test coverage across unit, integration, contract, UI, or end-to-end levels.

## Structure

- Tests SHOULD follow Arrange-Act-Assert or an equivalent clear structure.
- One test SHOULD verify one concept or behaviour.
- Tests MUST NOT assert multiple unrelated behaviours in a single case.
- Test setup SHOULD be explicit enough for the reader to understand the scenario.
- Shared test helpers MUST improve clarity and must not hide important behaviour.

## Naming

- Test names MUST describe the scenario and expected outcome.

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

- Unit tests SHOULD cover pure logic, utilities, validators, calculations, state transitions, and domain rules.
- Unit tests MUST cover normal cases, edge cases, and failure cases for important logic.
- Unit tests MUST NOT hit the network, real external services, production databases, or the file system unless the unit is specifically a file-system utility using isolated test storage.
- Unit tests SHOULD be fast and deterministic.

## Integration Tests

- Integration tests MUST cover important flows across modules, services, persistence, queues, or framework boundaries.
- Real implementations SHOULD be used unless the dependency is external, slow, unstable, costly, or unsafe.
- Mock only at system boundaries such as external APIs, third-party services, email, payments, identity providers, and infrastructure that cannot run safely in tests.
- Database integration tests MUST use isolated test data and must not depend on test execution order.

## Contract And API Tests

- Public APIs, events, webhooks, and shared schemas SHOULD have contract tests or schema validation.
- Breaking contract changes MUST be detected before release.
- Consumer-impacting changes SHOULD include compatibility tests or migration evidence.

## UI And End-To-End Tests

- Important user journeys SHOULD be tested end to end.
- UI tests SHOULD cover loading, empty, success, validation, authorization, and error states.
- End-to-end tests SHOULD focus on critical journeys, not every visual detail.
- Visual or snapshot tests MUST be stable and intentional, not broad snapshots that fail on noise.

## Test Doubles

- Prefer real implementations over mocks when practical.
- Do not mock internal business logic to make tests easier to pass.
- Use fakes, stubs, spies, or mocks only when they make the test safer, faster, clearer, or more deterministic.
- Use fake clocks or fixed timestamps instead of timing-based assertions.
- Randomness MUST be controlled with fixed seeds or deterministic inputs.

## Test Data

- Test data MUST be isolated from production data.
- Tests SHOULD create only the data they need.
- Tests MUST clean up after themselves when they use shared infrastructure.
- Sensitive production data MUST NOT be copied into test fixtures unless explicitly approved, minimized, and protected.

## Test Quality

- Tests MUST be deterministic. The same test must produce the same result for the same code.
- Flaky tests MUST be fixed or removed from required gates until fixed; they must not be ignored long-term.
- If a test cannot be made deterministic, the reason and risk MUST be documented.
- Failing tests MUST provide actionable failure messages.
- Test suites SHOULD be organized so fast feedback runs early and deeper checks run in CI or pre-release pipelines.

## CI Quality Gates

- Formatting, linting, type checks, builds, and relevant automated tests SHOULD run in CI.
- Required test gates MUST be reliable enough to block merge.
- Coverage thresholds, when used, MUST be enforced consistently and should not reward low-value tests.
- Critical paths SHOULD have meaningful behavioural coverage even if global coverage is high.
