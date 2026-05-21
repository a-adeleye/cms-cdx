# Architecture Rules

Architecture must make the system easier to change, test, scale, secure, and operate. Simplicity is preferred, but accidental coupling is not acceptable.

## Modular Architecture

- The system MUST be built from independent, well-defined modules.
- Each module MUST have a clear responsibility, owner, public interface, and boundary.
- Each module MUST be independently understandable and testable.
- Modules MUST NOT depend on internal implementation details of other modules.
- Communication between modules MUST happen through explicit interfaces, events, APIs, or contracts.
- Shared code MUST be minimal, generic, stable, and free of hidden product-specific behaviour.
- Cross-cutting concerns such as logging, configuration, authorization, validation, and telemetry SHOULD be standardized instead of reimplemented differently in every module.

## Separation of Concerns

- Code SHOULD be separated into these concerns where the application size justifies it:
  - Domain: business rules and domain model.
  - Application: use cases, orchestration, commands, workflows.
  - Infrastructure: databases, APIs, queues, external services, file systems.
  - UI or delivery: presentation, controllers, routes, screens, CLI, public transport adapters.
- Business logic MUST NOT depend on UI, framework, database, or transport-specific code.
- Infrastructure SHOULD be replaceable without changing domain behaviour.
- UI and transport layers MUST NOT bypass application or domain rules.
- Framework-specific code SHOULD be kept at the edges of the system.

## Dependency Rules

- Dependencies MUST point toward stable business rules and away from volatile implementation details.
- Circular dependencies are not allowed.
- Lower-level details MUST be injected behind interfaces when used by higher-level policy.
- A module MUST NOT reach directly into another module's database tables, private files, internal classes, or private state.
- New dependencies between modules SHOULD be reviewed as architecture changes, not treated as incidental imports.

## Contracts And Interfaces

- Public module boundaries MUST be represented by explicit contracts: interfaces, DTOs, schemas, API definitions, events, or command/query models.
- Contracts MUST use stable domain language and avoid leaking persistence or framework internals.
- Breaking contract changes MUST be versioned, migrated safely, or coordinated with all consumers.
- External consumers MUST NOT rely on undocumented behaviour.
- Internal data structures MUST NOT be exposed directly when that would create coupling.

## Data Ownership

- Every persistent data model MUST have a clear owning module or service.
- Only the owning module SHOULD write to its data.
- Other modules SHOULD access owned data through approved APIs, queries, events, or read models.
- Shared database access MUST NOT become an uncontrolled integration contract.
- Data duplication is allowed only when ownership, freshness, synchronization, and recovery rules are clear.

## File And Folder Structure

- Organize by feature, module, or bounded context by default.
- Organizing primarily by file type SHOULD be avoided when it scatters one feature across unrelated folders.

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

- Significant choices MUST include rationale, tradeoffs, and expected consequences.
- New frameworks, databases, messaging systems, or runtime platforms MUST have a documented reason.
- Architecture SHOULD evolve incrementally. Large rewrites require a migration plan, risk assessment, and rollback strategy.
- Exceptions to layering or dependency rules MUST be documented with a clear expiry or remediation path.
