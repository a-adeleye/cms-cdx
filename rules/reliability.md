# Reliability Rules

Reliable systems behave predictably, fail safely, recover gracefully, and provide enough visibility to diagnose production issues.

## Error Handling

- Expected errors MUST be handled explicitly.
- Errors MUST NOT be silently swallowed.
- User-facing error messages MUST be safe, useful, and free of internal details.
- Logs MUST include enough context to debug the issue without exposing secrets or sensitive data.
- Unexpected failures MUST fail safely and preserve data integrity.
- Error handling MUST distinguish validation errors, authorization failures, not-found cases, dependency failures, conflict states, and internal failures where relevant.
- Critical failures SHOULD include correlation IDs or request IDs for traceability.

## Validation

- All external input MUST be validated before use.
- API requests, form data, environment variables, files, messages, webhooks, and third-party responses MUST be validated.
- Client-side validation MUST NOT be treated as a security or integrity boundary.
- Validation failures MUST return clear, consistent, non-sensitive errors.
- Invalid external data SHOULD be rejected before it reaches domain logic.

## State Management

- State MUST be predictable and have a clear owner.
- Duplicated state SHOULD be avoided unless synchronization rules are explicit.
- Loading, success, empty, partial, stale, and error states MUST be represented clearly where users or workflows depend on them.
- Async flows MUST guard against race conditions, stale writes, duplicate submissions, and out-of-order responses.
- Critical state transitions SHOULD be modelled explicitly rather than inferred from scattered flags.

## External Services

- External services MUST be assumed unreliable.
- Network calls MUST have timeouts.
- Retries MUST be bounded and used only when safe.
- Retry logic SHOULD use backoff and jitter to avoid retry storms.
- Non-idempotent operations MUST NOT be retried blindly.
- Long-running or failure-prone external work SHOULD be moved out of synchronous request flows.
- The system MUST store enough information to recover, reconcile, or safely compensate failed external operations.
- Circuit breakers, bulkheads, or graceful degradation SHOULD be used for critical dependencies when failure would cascade.

## Observability

- Production systems MUST emit useful logs, metrics, and traces appropriate to their risk and complexity.
- Logs MUST be structured where practical and include correlation identifiers.
- Metrics SHOULD cover traffic, latency, errors, saturation, queue depth, job failures, and dependency health.
- Alerts MUST be actionable and tied to user impact or operational risk.
- Dashboards and runbooks SHOULD exist for critical production flows.

## Data Integrity

- Partial updates that can leave data inconsistent MUST be avoided.
- Transactions or atomic operations MUST be used where consistency is required.
- Idempotency keys or equivalent safeguards SHOULD be used for payments, webhooks, retries, and user-submitted commands that may repeat.
- Audit fields SHOULD be present where useful, such as `createdAt`, `updatedAt`, `createdBy`, and `updatedBy`.
- Destructive changes MUST require explicit intent and safe authorization.
- Data migrations MUST be repeatable, reversible where practical, and safe for existing production data.

## Recovery And Operations

- Critical workflows MUST have a defined recovery path.
- Background jobs MUST be restart-safe and idempotent where practical.
- Failed jobs MUST be observable and retryable or manually recoverable.
- Backups, restore procedures, and disaster recovery expectations SHOULD be defined for production data.
- New production features SHOULD identify their operational risks before release.
