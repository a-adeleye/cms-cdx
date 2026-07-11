# Architecture Rules

Read this file when creating or modifying modules, services, layers, or folder structure, or when designing or changing APIs, events, or webhooks. Architecture must make the system easier to change, test, scale, secure, and operate; that purpose is the tie-breaker when specific rules leave room for judgment.

## Modules

A **module** is a declared unit in the repository's module artifact — workspace/package configuration, Nx project graph, or a `MODULES.md` registry — listing its name, owner, purpose, public entry point, and owned data stores. All boundary rules below apply to declared modules.

- **ARCH-01.** Each repository MUST declare its modules in such an artifact; a module's public surface is what its declared entry point exports (e.g. an `index.ts` barrel or a package `exports` field). Everything else is internal.
- **ARCH-02.** Each module MUST state its purpose in one sentence in the module artifact, and its tests MUST run without booting sibling modules. Gate: Review.
- **ARCH-03.** Modules MUST NOT depend on another module's internals — its non-exported files, database tables, private classes, or private state. Cross-module deep imports MUST fail CI via boundary tooling (dependency-cruiser, ESLint boundary rules, Nx `enforce-module-boundaries`, ArchUnit, or equivalent). Gate: CI.
- **ARCH-04.** Communication between modules MUST happen through explicit interfaces, events, APIs, or contracts.
- **ARCH-05.** Modules MUST NOT form dependency cycles; this MUST be enforced by the same boundary tooling as ARCH-03. File-level import cycles within a module SHOULD NOT be introduced. Gate: CI.
- **ARCH-06.** A pull request that introduces a new module-to-module dependency edge MUST declare the edge and its justification in the PR description and update the boundary-tooling configuration; the owning team of the depended-on module approves it in review. Agents MUST NOT add a new edge silently. Gate: Evidence.
- **ARCH-07.** Shared code MUST NOT import from feature modules and MUST NOT branch on feature or product identifiers. Trivial code SHOULD be duplicated rather than widening a shared module.
- **ARCH-08.** Cross-cutting concerns (logging, configuration, authorization, validation, telemetry) SHOULD be standardized once instead of reimplemented per module.

## Separation of Concerns

- **ARCH-09.** Code SHOULD be separated into these concerns: Domain (business rules and domain model), Application (use cases, orchestration, workflows), Infrastructure (databases, APIs, queues, external services, file systems), and UI/delivery (presentation, controllers, routes, transport adapters).
- **ARCH-10.** Business rules MUST be unit-testable without a running framework, database, or network. DI registration metadata and serialization annotations do not count as framework dependence; framework APIs inside business-rule control flow do. This applies whether or not distinct layers are implemented.
- **ARCH-11.** Domain and Application code MUST access I/O-performing dependencies (databases, network services, file systems, message brokers, clock, randomness) through interfaces they own. Pure in-process libraries MAY be used directly.
- **ARCH-12.** No code path MAY circumvent domain validation, authorization, or business rules — UI and transport layers included.

## Contracts and Interfaces

- **ARCH-13.** Public module boundaries MUST be represented by explicit contracts: interfaces, DTOs, schemas, API definitions, events, or command/query models.
- **ARCH-14.** Contracts MUST use stable domain language. Persistence entities, ORM models, and framework types MUST NOT appear in public contracts; map to DTOs or schemas at the boundary. Gate: Review.
- **ARCH-15.** Breaking contract changes MUST be documented, MUST be communicated to consumers (a linked announcement or a versioning/deprecation plan in the PR or release record), and MUST be versioned or migrated safely. Gate: Evidence.
- **ARCH-16.** Code in this organization MUST NOT rely on undocumented behaviour of APIs it consumes. Providers MAY change undocumented behaviour without the breaking-change process and MUST say so in their API documentation.

## API Design

Applies to HTTP APIs, events, and webhooks exposed beyond a single module.

- **ARCH-17.** HTTP methods and status codes MUST follow standard semantics (safe GET, idempotent PUT/DELETE, 4xx for caller errors, 5xx for server errors).
- **ARCH-18.** Error responses MUST use one documented shape per API, carrying a machine-readable error code and a safe human-readable message. The error categories in reliability.md REL-04 map to distinct status/error codes.
- **ARCH-19.** Collection endpoints MUST paginate (see SCL-11); unsafe operations that may be retried MUST support idempotency keys or equivalent deduplication (this tightens REL-19 to MUST for APIs exposed beyond a single module).
- **ARCH-20.** Each public API MUST state its versioning and evolution strategy in its documentation.
- **ARCH-21.** Timestamps in contracts SHOULD be UTC ISO 8601; identifiers and field names SHOULD use one documented naming convention per API.

## Data Ownership

- **ARCH-22.** Every persistent data model MUST have one owning module or service, recorded in the module artifact.
- **ARCH-23.** Non-owning modules MUST NOT write to another module's data store, and MUST NOT query its tables or collections directly; they SHOULD use published APIs, events, or read models. Direct read access requires an exception per GEN-05.
- **ARCH-24.** Data MAY be duplicated across modules only when ownership, freshness, synchronization, and recovery rules are documented where the duplicate is introduced. Gate: Evidence.

## File and Folder Structure

- **ARCH-25.** Code SHOULD be organized by feature, module, or bounded context — not primarily by file type when that scatters one feature across unrelated folders.
- **ARCH-26.** Shared code SHOULD live in a designated shared location (e.g. `/shared` or `/lib`), subject to ARCH-07. Each feature folder corresponds to one declared module or a slice within one.

Preferred:

```text
/features
  /auth
    auth.service.ts
    auth.model.ts
    auth.controller.ts
  /orders
    orders.service.ts
    orders.model.ts
```

Avoid as the primary structure:

```text
/services
/controllers
/models
```

## Architecture Decisions

- **ARCH-27.** An ADR (per GEN-11, default `docs/adr/`) MUST be recorded, with context, decision, alternatives, and consequences, when a change: introduces a new framework, database, messaging system, or runtime platform; adds a new cross-module dependency direction; makes a layering or boundary exception; breaks a public contract; or starts a large rewrite. A reviewer MAY require an ADR for comparable decisions not listed. Gate: Evidence.
- **ARCH-28.** Architecture SHOULD evolve incrementally. Large rewrites MUST have a migration plan, risk assessment, and rollback strategy before the first PR merges.
- **ARCH-29.** Exceptions to layering or dependency rules MUST follow the exception format in engineering-rules.md GEN-05 (reason, risk, owner, expiry date, mitigation).
