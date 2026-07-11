# Delivery Rules — CI/CD, Infrastructure, and Releases

Read this file when changing CI/CD pipelines, build or deploy configuration, environments, or infrastructure, or when preparing a release. Migration mechanics live in reliability.md REL-29; feature-flag lifecycle in reliability.md §Feature Flags; the CI quality gates in testing.md §CI Quality Gates.

## Pipelines and Infrastructure

- **DEL-01.** Pipeline and infrastructure definitions MUST be version-controlled and reviewed as code, under the same review rules as application code.
- **DEL-02.** Required checks MUST block merge (TST-30); pipeline changes MUST NOT remove or weaken a required gate without an exception per GEN-05.
- **DEL-03.** Artifacts SHOULD be built once and promoted unchanged across environments.
- **DEL-04.** Production MUST NOT be changed manually outside the pipeline; break-glass changes require the exception process and MUST be reconciled back into code.
- **DEL-05.** Staging or pre-production SHOULD mirror production configuration closely enough that deploy and migration behaviour is representative.
- **DEL-06.** Deployment credentials MUST be least-privilege and scoped per environment (SEC-33).

## Releases

- **DEL-07.** Every release MUST be identifiable (version or commit SHA) and traceable to its changes. Gate: Evidence.
- **DEL-08.** Every release MUST have a documented rollback plan, or an approved written remediation plan where rollback is technically impossible (GEN-22).
- **DEL-09.** Releases of risky changes (engineering-rules.md §Definitions) SHOULD roll out staged or canary, behind a flag where REL-34 applies.
- **DEL-10.** Post-release verification steps MUST be defined for critical flows and executed after deploy; the result is recorded in the release record. Gate: Evidence.
- **DEL-11.** Shared libraries MUST use semantic versioning, and consumer-visible changes SHOULD ship release notes or a changelog entry.

## Environments

- **DEL-12.** Environment-specific configuration MUST be injected through the deployment platform or approved configuration mechanisms, never hardcoded (MNT-35).
- **DEL-13.** Production data MUST NOT be used in non-production environments except via an exception per GEN-05 approved by the data owner or security team, applying the minimization and protection requirements of TST-24.
