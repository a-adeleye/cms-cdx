# Security Rules

Read this file when handling authentication, authorization, secrets, input, user data, or reviewing supply-chain risk. Security must be designed into every feature; hidden UI controls, client-side checks, and undocumented assumptions are not security controls.

## Risk Classification

Used by the conditional rules below and by other files.

- **Sensitive data**: credentials, tokens, personal data, payment data, and regulated data (see engineering-rules.md §Definitions).
- **High-risk / privileged surface**: authentication flows, authorization logic, payment handling, administrative actions, data export, and cryptography.
- **SEC-01.** Borderline classifications MUST be decided by the service owner together with a security reviewer and recorded as a label in the PR or an ADR. Gate: Evidence.
- **SEC-02.** Where a rule requires an "approved" algorithm, sanitizer, or system and no approved list exists in the organization, the stated default in that rule MUST be applied and the missing list flagged in the PR (agents: per AGT-06).

## Authentication

- **SEC-03.** Every endpoint, job, and handler MUST require authentication unless it is explicitly declared public, with the declaration and justification visible in code or configuration. Gate: Review.
- **SEC-04.** User identity MUST NOT be trusted from the frontend alone.
- **SEC-05.** Session and token handling MUST implement these controls: session IDs rotated on login and privilege change; sessions invalidated server-side on logout; tokens carry a defined expiry and are revocable server-side; tokens never appear in URLs or logs; browser session cookies set `HttpOnly`, `Secure`, and `SameSite=Lax` or `Strict` (a documented cross-site flow MAY use `SameSite=None`, which still requires `Secure`). Gate: Review.
- **SEC-06.** Passwords, when used, MUST be hashed with an approved password-hashing algorithm (default: Argon2id, scrypt, or bcrypt with current OWASP parameters) and never stored or logged in plain text.
- **SEC-07.** Multi-factor authentication SHOULD be supported for privileged or high-risk access.
- **SEC-08.** Account lifecycle flows MUST implement: logout invalidates the current session and its refresh tokens server-side (a stale token then fails with 401, asserted by a test); account deletion and credential reset invalidate all active sessions and refresh tokens; recovery uses single-use, expiring, unguessable tokens. Gate: Evidence (test attached to the change).
- **SEC-09.** Authentication failures MUST NOT reveal whether an account exists, unless that is an intentional, documented product decision.

## Authorization

- **SEC-10.** Every protected action MUST check authorization on the server, service, or database enforcement layer — for reads, writes, deletes, exports, background jobs, webhooks, and administrative actions alike. This is the canonical list of acceptable enforcement layers.
- **SEC-11.** Hidden UI controls and client-side route guards MUST NOT be relied on for security.
- **SEC-12.** Least-privilege access MUST be used for users, services, jobs, databases, and cloud resources, and object-level and tenant-level access MUST be enforced for user-owned or tenant-owned data.
- **SEC-13.** Permission-change history on high-risk surfaces SHOULD be reviewable as a queryable audit trail, beyond the individual log events SEC-34 requires.

## Secrets

- **SEC-14.** Secrets — API keys, private keys, tokens, passwords, certificates, credentials — MUST NOT be committed to source control, logged, exposed in errors, embedded in client bundles, or sent to analytics tools. Gate: CI (secret scanner) + Review.
- **SEC-15.** Secrets MUST be stored in environment variables, secret managers, or approved secure configuration systems, with development and production secrets separated.
- **SEC-16.** Secret rotation MUST be possible without code changes.
- **SEC-17.** A secret that has been committed, logged, or otherwise exposed MUST be treated as compromised and rotated immediately; removing it from history or logs is not sufficient remediation.

## Input and Output Safety

This section owns the input-validation baseline; reliability.md §Validation owns failure behaviour.

- **SEC-18.** All external input MUST be validated, normalized, and constrained at the trust boundary before use (mechanics per REL-06).
- **SEC-19.** Injection MUST be prevented for SQL, NoSQL, shell commands, templates, LDAP, XML, and other interpreters — parameterized queries and safe APIs, never string assembly from untrusted input. Output MUST be escaped or encoded for its target context, and unsafe HTML MUST NOT be rendered unless sanitized with an approved sanitizer (default: the framework's built-in sanitizer or DOMPurify).
- **SEC-20.** Server-side requests to user-influenced URLs MUST validate the destination against an allowlist and block internal address ranges (SSRF). File paths derived from external input MUST be canonicalized and confined to an allowed root. Redirect targets from user input MUST be allowlisted.
- **SEC-21.** Native object deserialization (e.g. pickle, Java serialization, unsafe YAML load) MUST NOT be used on untrusted data; untrusted data MUST be parsed with data-only formats (JSON, protobuf) and validated against a schema. Exceptions per GEN-05.
- **SEC-22.** File uploads MUST validate type, size, content, storage location, and access permissions, and APIs MUST validate content type, payload size, and schema for requests that carry a body.

## Data Protection and Privacy

- **SEC-23.** Newly collected personal or sensitive data MUST have a defined business purpose recorded where the data model lives; data minimization SHOULD be applied to all other data.
- **SEC-24.** Sensitive data MUST be protected in transit and at rest, and MUST NOT appear in logs, analytics, error messages, URLs, screenshots or recordings attached to reviews, or client-side storage without an exception per GEN-05. Responses to untrusted users MUST NOT include stack traces, raw queries, secrets, or infrastructure details; internal IDs and implementation details SHOULD NOT be exposed.
- **SEC-25.** For regulated data, retention, deletion, export, and privacy obligations MUST be defined before storing it, and deletion/export requests MUST be technically satisfiable (no orphaned copies in logs, caches, or backups beyond stated policy). For other personal data these SHOULD be defined. Data residency constraints MUST be checked before adding new storage locations or third-party processors.

## Web and API Security

- **SEC-26.** All production traffic MUST be served over HTTPS; HTTP MUST be redirected to HTTPS or blocked.
- **SEC-27.** CSRF protection MUST be used wherever browser-based authenticated state-changing requests are possible, enforced by the backend regardless of any client-side helper.
- **SEC-28.** CORS allowlists MUST enumerate specific trusted origins; wildcard origins or reflecting the request `Origin` header MUST NOT be used on endpoints that require credentials or return non-public data.
- **SEC-29.** Web responses MUST set `Strict-Transport-Security` and `X-Content-Type-Options: nosniff`. Pages rendering authenticated content MUST prevent framing (`frame-ancestors` or `X-Frame-Options`) unless embedding is an intentional, documented feature. A Content-Security-Policy SHOULD be defined for web applications.
- **SEC-30.** Publicly reachable or unauthenticated endpoints SHOULD enforce rate limits and abuse prevention (throttling, bot protection, replay protection, or abuse monitoring).

## Dependencies and Supply Chain

Dependency lifecycle rules (adding, pinning, upgrading) live in maintainability.md §Dependencies; this section owns the security controls.

- **SEC-31.** Dependencies MUST be scanned for known vulnerabilities by a CI scanner in ecosystems with mainstream tooling; findings MUST be triaged by severity and tracked to resolution or an accepted-risk exception (GEN-05). Gate: CI.
- **SEC-32.** Abandoned, unmaintained, or suspicious packages SHOULD be avoided, and install scripts, transitive-dependency risk, and maintainer reputation SHOULD be assessed and documented for critical systems.
- **SEC-33.** Build and deployment pipelines MUST protect credentials and production artifacts (see delivery.md).

## Audit, Monitoring, and Response

- **SEC-34.** Production systems MUST log these security events with actor, action, target, result, time, and correlation ID: authentication success and failure, authorization denials, permission changes, and administrative actions. Broader security-event logging SHOULD be added by risk. Audit logs MUST NOT contain secrets or excessive personal data.
- **SEC-35.** Systems MUST define alerts for critical security events; each alert follows REL-24 (runbook or stated action, routed owner).
- **SEC-36.** Changes to high-risk surfaces MUST include relevant tests or review evidence, and SHOULD receive focused security review before release.
- **SEC-37.** Incident response expectations MUST be defined for production systems (severity, escalation path — see REL-42); a leaked secret triggers SEC-17.
