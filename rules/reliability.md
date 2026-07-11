# Reliability Rules

Read this file when handling errors, external services, validation, state, data integrity, migrations, or feature flags. Reliable systems behave predictably, fail safely, recover gracefully, and provide enough visibility to diagnose production issues. Sections follow the failure lifecycle: prevent, detect, recover.

**Critical workflow** and **critical dependency** are defined in engineering-rules.md §Definitions; this file does not redefine them.

## Error Handling

- **REL-01.** Expected errors MUST be handled. An error is handled only when it is mapped to a defined outcome: propagated to the caller, retried, replaced by a documented fallback, or surfaced as a user-facing failure. Logging alone is not handling.
- **REL-02.** Errors MUST NOT be silently swallowed. Empty catch blocks, ignored error returns, and catch-log-continue on errors that affect the operation's outcome all count as swallowing. Gate: Review.
- **REL-03.** User-facing error messages MUST be safe and useful, free of internal details (see SEC-24 for what must never leak).
- **REL-04.** API and service boundaries MUST map failures to distinct categories — validation, authorization, not-found, conflict, dependency failure, internal — each with a distinct status or error code (shape per ARCH-18). Internal code MUST NOT collapse these categories before the boundary.
- **REL-05.** Unexpected failures MUST fail safely and preserve data integrity.

## Validation

The input-validation baseline (validate, normalize, constrain) is owned by security.md §Input and Output Safety; this section owns validation failure behaviour.

- **REL-06.** Data from external sources — API requests, form data, environment variables, files, messages, webhooks, third-party responses — MUST be validated against an explicit schema or type contract at the system boundary, before it reaches domain logic. Additional defense-in-depth checks deeper in the system MAY be applied.
- **REL-07.** Client-side validation MUST NOT be treated as a security or integrity boundary.
- **REL-08.** Validation failures MUST return clear, consistent, non-sensitive errors.

## State Management

- **REL-09.** Every piece of shared mutable state MUST have exactly one owning module or service responsible for writes; all other code reads it through that owner's interface. Gate: Review.
- **REL-10.** Duplicated state SHOULD be avoided unless synchronization rules are explicit (see ARCH-24).
- **REL-11.** Loading, success, empty, partial, stale, and error states MUST be represented explicitly where users or workflows depend on them.
- **REL-12.** Where an async hazard exists, the flow MUST apply a named mitigation: duplicate submissions → submit-disable or idempotency key; out-of-order responses → cancel or ignore superseded requests; stale writes → optimistic concurrency or version check; shared mutable state → single-writer ownership (REL-09). The reviewer checks each applicable hazard has a mitigation; novel hazards still require judgment. Gate: Review.
- **REL-13.** Critical state transitions SHOULD be modelled explicitly (state machine or equivalent) rather than inferred from scattered flags.

## External Services

External services are assumed unreliable; the rules below are the required consequences of that assumption. This section is the canonical owner of timeout and retry policy — other files reference it.

- **REL-14.** Every network call MUST have an explicit, finite timeout configured (connect and overall, where the client distinguishes them). Relying on a client-library default does not satisfy this rule unless the default is finite and the choice is documented.
- **REL-15.** Retries MUST be bounded by a maximum attempt count or a total time budget, and MUST only be applied to operations that are idempotent or protected by an idempotency key or deduplication safeguard.
- **REL-16.** Retry logic SHOULD use exponential backoff with jitter where repeated failure can amplify load.
- **REL-17.** Long-running or failure-prone external work SHOULD be moved out of synchronous request flows (see scalability.md §Background Work).
- **REL-18.** Before invoking a state-changing external operation (payment, provisioning, message send), the system MUST persist an operation record — identifier, status, and the data needed to retry, reconcile, or compensate — and MUST update it on completion. Read-only calls are exempt.
- **REL-19.** Idempotency keys or equivalent safeguards SHOULD be used for payments, webhooks, retries, and user-submitted commands that may repeat.
- **REL-20.** Circuit breakers, bulkheads, or graceful degradation SHOULD protect critical dependencies whose failure would cascade.

## Observability

- **REL-21.** Production services MUST emit structured logs and service-level metrics covering, at minimum, traffic, errors, and latency. Metrics SHOULD additionally cover saturation, queue depth, job failures, and dependency health where those exist, and services participating in multi-service request flows SHOULD emit distributed traces. Gate: Review.
- **REL-22.** Log entries emitted while serving a request or processing a job MUST include a correlation or request ID, and that ID MUST propagate across service and background-job boundaries.
- **REL-23.** Log levels MUST follow these semantics: ERROR = requires action; WARN = degraded but handled; INFO = significant state change; DEBUG = disabled in production by default. High-frequency paths SHOULD sample or aggregate log output.
- **REL-24.** Every alert MUST link a runbook or state the expected responder action, MUST fire on a condition tied to user impact or operational risk, and MUST be routed to an owner. Alerts that require no action MUST be converted to dashboards or deleted.
- **REL-25.** Dashboards and runbooks SHOULD exist for critical production flows.

## Data Integrity

- **REL-26.** Multi-step writes that must succeed or fail together MUST be executed in a transaction or an equivalent atomic or compensating mechanism (outbox, saga). A partial failure MUST NOT leave persisted data in a state no single successful operation could produce.
- **REL-27.** Mutable persisted business entities SHOULD carry `createdAt`, `updatedAt`, and actor fields.
- **REL-28.** Destructive operations (hard delete, bulk update, truncation) MUST require a distinct confirmation step or dedicated endpoint/command and MUST be authorized per security.md. They SHOULD be recoverable (soft delete or pre-operation backup) except where retention or privacy obligations require permanent deletion — in which case permanent deletion wins.
- **REL-29.** Schema and data migrations MUST be backward compatible with the currently deployed application versions (expand/migrate/contract): destructive steps (drop, rename, type narrowing) MUST ship only after no deployed code depends on the old shape. Migrations MUST be repeatable, MUST be tested against production-like data, and SHOULD be reversible; an irreversible migration MUST document a remediation plan (GEN-22). This is the canonical migration-safety rule. Gate: Evidence.

## Feature Flags

- **REL-30.** Every feature flag MUST record an owner and a removal condition or date when created; a flag whose removal condition is met is obsolete and MUST be removed (MNT-05).
- **REL-31.** A flag's default state MUST be safe when the flag system is unavailable.
- **REL-32.** Both states of a flag guarding a production path MUST be tested.
- **REL-33.** Flags MUST NOT substitute for authorization.
- **REL-34.** A change SHOULD be behind a feature flag when it cannot be fully verified before production, changes behaviour of a flow designated critical, or is intended for gradual rollout; the author states the classification in the PR.

## Recovery and Operations

- **REL-35.** Critical workflows MUST have a defined recovery path.
- **REL-36.** Background jobs MUST be restart-safe; idempotency requirements are owned by scalability.md SCL-24.
- **REL-37.** Failed jobs MUST be observable and either retryable or manually recoverable.
- **REL-38.** Services MUST handle termination signals by stopping intake of new work and completing or safely aborting in-flight work. Long-running services SHOULD expose health and readiness endpoints (or platform equivalents).
- **REL-39.** Backups, restore procedures, and disaster-recovery expectations SHOULD be defined for production data.
- **REL-40.** New production features MUST record their operational risks before release (GEN-24).

## Incidents

- **REL-41.** User-impacting production incidents MUST receive a blameless postmortem with action items tracked to closure, and incident fixes MUST include regression tests where the defect is testable.
- **REL-42.** Each production system MUST have a defined severity classification and escalation path; recurring incident causes SHOULD feed updates to runbooks and to these rules.
