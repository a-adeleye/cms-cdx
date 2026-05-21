# Security Rules

Security must be designed into every feature. Hidden UI controls, client-side checks, and undocumented assumptions are not security controls.

## Authentication

- Authentication MUST be required for protected actions.
- User identity MUST NOT be trusted from the frontend alone.
- Session and token handling MUST be secure against theft, replay, fixation, and unintended persistence.
- Passwords, when used, MUST be hashed with an approved password hashing algorithm and never stored or logged in plain text.
- Multi-factor authentication SHOULD be supported for privileged or high-risk access.
- Logout, session revocation, account deletion, and account recovery flows MUST be safe and complete.
- Authentication failures MUST NOT leak whether a specific account exists unless that is an intentional product decision.

## Authorization

- Every protected action MUST check authorization on the server, service, or database enforcement layer.
- Hidden UI controls MUST NOT be relied on for security.
- Authorization MUST be checked for reads, writes, deletes, exports, background jobs, webhooks, and administrative actions.
- Least privilege access MUST be used for users, services, jobs, databases, and cloud resources.
- Object-level and tenant-level access MUST be enforced for user-owned or tenant-owned data.
- Permission changes SHOULD be auditable for sensitive systems.

## Secrets

- Secrets, API keys, private keys, tokens, passwords, certificates, and credentials MUST NOT be committed to source control.
- Secrets MUST be stored in environment variables, secret managers, or approved secure configuration systems.
- Secrets MUST NOT be logged, exposed in errors, embedded in client bundles, or sent to analytics tools.
- Secret rotation MUST be possible without code changes.
- Development and production secrets MUST be separated.

## Input And Output Safety

- All external input MUST be validated, normalized, and constrained before use.
- Injection attacks MUST be prevented for SQL, NoSQL, shell commands, templates, LDAP, XML, and other interpreters.
- Output MUST be escaped or encoded for the target context.
- Unsafe HTML MUST NOT be rendered unless sanitized with an approved sanitizer.
- File uploads MUST validate type, size, content, storage location, and access permissions.
- Deserialization of untrusted data MUST be avoided or strictly controlled.

## Data Protection And Privacy

- Store only the data needed for a defined business purpose.
- Sensitive data MUST be protected in transit and at rest.
- Sensitive data MUST NOT be exposed in logs, analytics, errors, URLs, screenshots, or client-side storage unless explicitly approved.
- Internal IDs, stack traces, infrastructure details, and implementation details SHOULD NOT be exposed to untrusted users.
- User-owned and tenant-owned data MUST have proper access rules.
- Data retention, deletion, export, and privacy obligations SHOULD be defined for personal or regulated data.

## Web And API Security

- All production traffic MUST be served over HTTPS. HTTP MUST be redirected to HTTPS or blocked.
- CSRF protection MUST be used where browser-based authenticated state-changing requests are possible.
- CORS MUST be restrictive and intentional.
- Cookies, when used for authentication, MUST use secure attributes appropriate to the application.
- APIs MUST validate content type, payload size, schema, authentication, authorization, and rate limits where needed.
- Public endpoints SHOULD include abuse prevention such as throttling, bot protection, replay protection, or abuse monitoring.
- Error responses MUST avoid leaking stack traces, secrets, queries, or infrastructure details.

## Audit And Monitoring

- Security-sensitive events SHOULD be logged with actor, action, target, result, time, and correlation ID.
- Audit logs MUST avoid storing secrets or excessive personal data.
- Suspicious activity SHOULD be detectable for authentication, authorization, payment, data export, and administrative actions.
- Critical security alerts MUST be actionable and routed to an owner.

## Dependencies And Supply Chain

- Dependencies MUST be reviewed before adding them.
- Dependencies MUST be kept updated and scanned for known vulnerabilities where tooling exists.
- Abandoned, unmaintained, or suspicious packages SHOULD be avoided.
- Package installation scripts, transitive dependency risk, license risk, and maintainer reputation SHOULD be considered for critical systems.
- Build and deployment pipelines MUST protect credentials and production artifacts.

## Security Testing And Response

- Security-sensitive changes MUST include relevant tests or review evidence.
- Known vulnerabilities MUST be triaged by severity and tracked to resolution or accepted risk.
- High-risk features SHOULD receive focused security review before release.
- Incident response expectations SHOULD be defined for production systems.
