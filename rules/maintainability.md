# Maintainability Rules

Maintainable code is easy to understand, change, test, review, and delete. Enterprise-grade code must optimize for the next developer as well as the first delivery.

## Code Structure

- Files, modules, classes, functions, and components MUST have one clear responsibility.
- Functions exceeding 30 lines or three levels of nesting are a signal to split into named operations that express intent.
- Complex logic MUST be isolated from framework, UI, transport, and persistence details.
- Duplicate logic SHOULD be removed when the same rule exists in more than one place.
- Reusable abstractions MUST be created only when they remove real duplication or clarify a stable concept.
- Dead code, unused branches, obsolete feature flags, and abandoned TODOs MUST be removed.
- Code paths for critical business rules SHOULD be easy to locate from the feature entry point.

## Naming

- Names MUST describe intent using the language of the business domain.
- Avoid vague names such as `data`, `item`, `temp`, `obj`, `result`, `manager`, or `helper` unless the local context makes the meaning obvious.
- Boolean names SHOULD read as true/false statements, for example `isActive`, `hasPermission`, `canEdit`, or `shouldRetry`.
- Public APIs, events, commands, and database fields MUST use consistent terminology.
- Names MUST NOT hide side effects. A function named like a query SHOULD NOT mutate state.

## Comments And Documentation

- Prefer self-explanatory code over comments that repeat the implementation.
- Comments SHOULD explain why a non-obvious decision exists, not what each line does.
- Outdated comments MUST be removed or corrected immediately.
- Public modules, complex workflows, operational runbooks, and non-obvious configuration MUST be documented.
- Workarounds MUST include the reason, linked issue when available, and removal condition.

## Dependencies

- Do not add dependencies for simple tasks that the language or platform already handles well.
- Every dependency MUST have a clear purpose, owner, and maintenance expectation.
- Dependencies MUST be actively maintained, licensed appropriately, and compatible with security requirements.
- Unused dependencies MUST be removed.
- Dependency upgrades SHOULD be small, frequent, and tested.
- A dependency that changes architecture, runtime behaviour, build output, licensing, or security posture MUST be explicitly reviewed.

## Refactoring

- Refactor when code becomes hard to understand, test, change, or operate.
- Do not mix large refactors with unrelated feature changes.
- Preserve existing behaviour unless the task explicitly requires changing it.
- Behaviour-preserving refactors SHOULD be covered by tests before changes are made.
- Large refactors MUST have a migration strategy that avoids blocking normal delivery.

## Consistency

- Formatting MUST be automated and applied consistently.
- Linting and static analysis rules SHOULD be enforced in CI.
- Existing project conventions SHOULD be followed unless there is a documented reason to change them.
- Similar problems SHOULD be solved in similar ways across the codebase.
- New patterns MUST be easy for the team to understand and repeat.

## Configuration

- Configuration MUST be explicit, typed or validated where practical, and documented when non-obvious.
- Environment-specific values MUST NOT be hardcoded into application logic.
- Missing or invalid required configuration MUST fail fast with a safe, actionable error.
- Defaults MUST be safe for production or clearly marked as development-only.
