# Engineering Rules Enforcement

You MUST follow the engineering standards defined in the `rules/` directory of this project.

Before writing code, reviewing, or preparing a change, read `rules/engineering-rules.md` first. It defines the rule enforcement levels and the Definition of Done that apply to implementation work:

- **MUST**: Required. A violation blocks merge or release.
- **SHOULD**: Required by default. Exceptions must be justified.
- **MAY**: Optional guidance to apply when appropriate.

All code changes and implementation work MUST satisfy the Definition of Done in `rules/engineering-rules.md` before being considered complete.

Admin frontend work MUST use Angular.

## Rule Files

| File | Read when |
|---|---|
| `rules/engineering-rules.md` | Always — defines levels, Definition of Done, release readiness |
| `rules/architecture.md` | Creating or modifying modules, services, layers, or folder structure |
| `rules/angular.md` | Writing any Angular component, service, directive, pipe, or template |
| `rules/maintainability.md` | Writing or refactoring any application code |
| `rules/reliability.md` | Handling errors, external services, validation, state, or data integrity |
| `rules/scalability.md` | Writing database queries, caching, background jobs, or performance-sensitive code |
| `rules/security.md` | Handling authentication, authorization, secrets, input, or user data |
| `rules/testing.md` | Writing or reviewing tests of any kind |
| `rules/code-review.md` | Reviewing a pull request or preparing a change for review |

For tasks that cross multiple concerns, read all files whose trigger condition applies. When in doubt, read the relevant file. The rules take precedence over your defaults.

## Mandatory Behavior

- Read the relevant rule files **before** writing code, reviewing, or preparing a change.
- If a rule conflicts with your default behavior, the rule wins.
- Do not surface rule violations after generating code — prevent them before writing a single line.
- If a task cannot be completed without violating a MUST rule, stop and explain why before proceeding.

## If a Rule Is Violated

- Do not silently proceed.
- Identify the specific rule that was or would be violated.
- Correct the implementation.
- Briefly state which rule applies and what was changed.
