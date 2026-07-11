# Scalability Rules

Read this file when writing database queries, caching, background jobs, or performance-sensitive code. Scalable systems continue to perform acceptably as users, data, traffic, integrations, and business complexity grow.

## Definitions

- **Unbounded collection**: data whose size grows with users, tenants, or time and has no enforced cap.
- **High-growth data**: an unbounded collection with sustained growth.
- **Expensive endpoint/operation**: one that exceeds the service's latency budget or calls a paid or rate-limited integration.
- **Hot path**: code executed per request or per item in high-volume processing.

Where a rule below still requires judgment, the reviewer decides using the capacity assumptions stated for the feature (SCL-01).

## Capacity Assumptions

- **SCL-01.** Features SHOULD state expected usage when scale affects design: users, records, request rate, payload size, data retention, and integration volume — in the PR or design record.
- **SCL-02.** Features MUST NOT be designed or tested only against demo-sized data unless the feature is explicitly temporary.
- **SCL-03.** Hardcoded limits MUST be named constants with a documented rationale and MUST be enforced by validation; exceeding a limit MUST produce a clear error to the caller or operator — never silent truncation or data loss.
- **SCL-04.** Growth-sensitive designs SHOULD document the first constraint expected to break and the intended mitigation.

## Architecture

- **SCL-05.** Designs MUST handle the usage stated under SCL-01 without requiring changes to module boundaries or storage technology; where no capacity assumptions were stated, the reviewer requires them first.
- **SCL-06.** Read-heavy and write-heavy paths SHOULD be separated when measured or stated load profiles differ by roughly an order of magnitude, or when their consistency requirements conflict. Do not introduce CQRS or read replicas below that bar.
- **SCL-07.** Services intended to run as multiple instances MUST NOT rely on instance-local state as a source of truth (in-memory sessions, local files as records, in-process schedulers without leader election); such state MUST be externalized to shared storage.
- **SCL-08.** A change that routes new or significantly increased load through a shared serialized resource (single writer, global lock, singleton service, single queue consumer) MUST name that resource and its expected headroom in the PR description or design record. Gate: Evidence.

## Database

- **SCL-09.** Hand-written queries MUST NOT use `SELECT *`, and queries MUST NOT fetch rows only to discard them with filtering the database could have applied (app-side filtering is acceptable only where the predicate cannot be expressed in the query). Column-level projection SHOULD be used when the entity is large or the path is hot; default ORM entity hydration is otherwise acceptable.
- **SCL-10.** Unbounded collections MUST NOT be loaded entirely into memory; bounded reference data (fixed enumerations, small config tables) is exempt.
- **SCL-11.** Lists exposed to users or APIs MUST use pagination, limits, or streaming. This is the canonical pagination rule (ARCH-19 references it).
- **SCL-12.** N+1 query patterns MUST be avoided in production paths, and sorting, filtering, and joins MUST be reviewed for index support at the expected data volume. Fields SHOULD be indexed based on the queries the code actually issues, not speculatively.
- **SCL-13.** Data retention, archival, and deletion strategies SHOULD be defined for high-growth data, and deeply nested documents SHOULD be avoided when nested parts need independent querying or updates.

## Performance

- **SCL-14.** Code SHOULD avoid unnecessary loops, repeated queries, chatty network calls, and expensive synchronous computations; hot paths SHOULD avoid avoidable allocation, serialization, and repeated transformation.
- **SCL-15.** Expensive operations SHOULD be measured before and after optimization.
- **SCL-16.** Responses MUST NOT include fields no consumer uses, and responses exceeding the organization's size threshold (project-configurable; default 1 MB) MUST be paginated, streamed, or compressed, or the exception justified in the PR.
- **SCL-17.** Heavy features SHOULD be lazy-loaded, streamed, paginated, deferred, or moved to background work, and performance-sensitive code SHOULD have explicit limits on input size and execution time.

## Caching

- **SCL-18.** New caches SHOULD be justified by a measured or documented expected performance problem.
- **SCL-19.** Each cache MUST declare its maximum acceptable staleness (TTL or explicit invalidation trigger) where it is introduced; the feature owner accepts that staleness in review. Gate: Evidence.
- **SCL-20.** Cache keys MUST include every input that affects the result, including tenant/user scope for user-specific data. Raw unbounded user input MUST be hashed or normalized before use in a key so key cardinality stays bounded.
- **SCL-21.** Caches MUST have a bounded size or an eviction policy; caches fronting expensive computations SHOULD use stampede protection (single-flight locking or jittered TTLs).
- **SCL-22.** The system MUST continue to behave correctly if the cache is empty or unavailable.

## Background Work

- **SCL-23.** Work that is long-running, scheduled, or retry-prone SHOULD run via queues or background jobs, outside request/response flows.
- **SCL-24.** Background jobs that can be retried or restarted MUST be idempotent or protected by deduplication; exceptions are documented per GEN-05. This is the canonical job-idempotency rule (REL-36 references it).
- **SCL-25.** Job consumers MUST define a retry limit and a dead-letter or poison-message destination; producers MUST have a defined behaviour when the queue is full or consumers lag.
- **SCL-26.** Queue depth and job failure rates SHOULD be monitored for production workloads (telemetry baseline per REL-21).

## Concurrency and Load

- **SCL-27.** Concurrent requests MUST NOT corrupt data or create duplicate side effects. Where duplicate execution is plausible (retries, double submits, concurrent writers to the same records), the protecting mechanism — unique constraint, transaction isolation or locking, optimistic concurrency, or idempotency key — MUST be identifiable in the code or stated in the PR. Gate: Review.
- **SCL-28.** Rate limits, throttling, or backpressure SHOULD protect expensive endpoints and integrations.
- **SCL-29.** Bulk operations MUST have bounded memory and time behaviour.
- **SCL-30.** Long-held locks, global locks, and single-thread bottlenecks SHOULD be avoided (declare unavoidable ones per SCL-08).
- **SCL-31.** Workflows that depend on ordering MUST document how ordering is preserved.

## API and Payload Growth

- **SCL-32.** Public contracts SHOULD be versioned or evolved backward-compatibly (breaking changes per ARCH-15).
- **SCL-33.** Clients SHOULD be able to request only the fields or resources they need when payloads exceed the SCL-16 threshold.
- **SCL-34.** Large imports, exports, and reports SHOULD be asynchronous or streamed. Inbound payload-size limits are owned by SEC-22.
