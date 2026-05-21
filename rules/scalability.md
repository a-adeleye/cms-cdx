# Scalability Rules

Scalable systems continue to perform acceptably as users, data, traffic, integrations, and business complexity grow.

## Capacity Assumptions

- Features SHOULD state expected usage when scale affects design: users, records, request rate, payload size, data retention, and integration volume.
- Do not design only for demo-sized data unless the feature is explicitly temporary.
- Hardcoded limits MUST be intentional, documented, validated, and surfaced safely to users or operators.
- Growth-sensitive decisions SHOULD identify the next scaling constraint.

## Architecture

- Features MUST be designed so normal growth does not require a major rewrite.
- Read-heavy and write-heavy paths SHOULD be separated when their performance, consistency, or storage needs differ.
- Modules SHOULD be able to scale independently when load profiles differ.
- Shared bottlenecks MUST be visible and intentional.
- Tightly coupled logic that prevents horizontal scaling SHOULD be avoided.

## Database

- Queries MUST fetch only the data needed.
- Large collections MUST NOT be loaded entirely into memory.
- Lists exposed to users or APIs MUST use pagination, limits, or streaming.
- Frequently queried fields SHOULD be indexed based on real query patterns.
- N+1 query patterns MUST be avoided in production paths.
- Sorting, filtering, and joins MUST be reviewed for index support and expected data volume.
- Deeply nested data SHOULD be avoided when nested parts need independent querying or updates.
- Data retention, archival, and deletion strategies SHOULD be defined for high-growth data.

## Performance

- Avoid unnecessary loops, repeated queries, chatty network calls, and expensive synchronous computations.
- Expensive operations SHOULD be measured before and after optimization.
- Payloads MUST be kept as small as practical.
- Heavy features SHOULD be lazy-loaded, streamed, paginated, deferred, or moved to background work.
- Performance-sensitive code SHOULD have clear limits on input size and execution time.
- Hot paths SHOULD avoid avoidable allocation, serialization, and repeated transformation.

## Caching

- Cache only when it improves a measured or strongly expected performance problem.
- Cached data MUST have clear ownership, invalidation, expiry, and consistency expectations.
- Cache keys MUST be safe, bounded, and include all inputs that affect the result.
- Stale data risks MUST be acceptable for the workflow.
- The system MUST continue to behave correctly if the cache is empty or unavailable.

## Background Work

- Heavy, slow, scheduled, or retry-prone work SHOULD run outside request/response flows.
- Queues or background jobs SHOULD be used for long-running tasks.
- Background jobs MUST be idempotent where retries are possible.
- Job processing MUST handle backpressure, retry limits, poison messages, and observability.
- Queue depth and job failure rates SHOULD be monitored for production workloads.

## Concurrency And Load

- Concurrent requests MUST NOT corrupt data or create duplicate side effects.
- Rate limits, throttling, or backpressure SHOULD protect expensive endpoints and integrations.
- Bulk operations MUST have bounded memory and time behaviour.
- Long locks, global locks, and single-thread bottlenecks SHOULD be avoided.
- Workflows that depend on ordering MUST document how ordering is preserved.

## API And Payload Growth

- APIs MUST avoid unbounded response sizes.
- Public contracts SHOULD be versioned or evolved backward-compatibly.
- Clients SHOULD be able to request only the fields or resources they need when payload size is significant.
- Large imports, exports, and reports SHOULD be asynchronous or streamed.
