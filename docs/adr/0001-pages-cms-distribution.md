# ADR 0001: Adopt Pages CMS as the Git-backed editor

- Status: accepted for local development
- Date: 2026-08-10
- Upstream: `hunvreus/pagescms` v2.1.8 at `6f4e860a35d934406580287e7042e5e111e207a1`
- License: MIT (preserved in `apps/pagescms/LICENSE`)
- Rules: ARCH-01, ARCH-04, ARCH-11, ARCH-22, ARCH-27, MNT-19, MNT-24

## Decision

Vendor the pinned Pages CMS snapshot under `apps/pagescms` and extend it through an owned repository-store boundary. Git-tracked JSON and media files are the only content source; Astro reads the same files directly. The existing Go/Postgres CMS path remains available during migration.

The local adapter is enabled only by an explicit development flag, authorizes one repository/ref, confines paths to its canonical root, and uses content digests for optimistic concurrency. Normal Pages CMS GitHub behavior remains available for a later production phase.

AI drafting uses a server-side provider boundary, a deterministic mock by default, and an opt-in OpenAI Responses implementation. Generated values are reviewable form changes and are never saved automatically.

## Consequences

We gain Pages CMS's schema editor and media workflow without rebuilding them. We accept maintenance of a small fork and must consciously rebase upstream releases. Local mode intentionally omits GitHub Actions, collaborators, remote history, scheduling, and deployment.

## Rollback and replacement

Stop using `apps/pagescms` and the `dev:cms`/`dev:local` scripts; the prior `apps/admin` and `apps/api` path remains intact. Content is plain JSON/media in Git and can be consumed by another editor without migration. Upgrades import a new pinned snapshot, retain its MIT license, review the extension diff, and rerun the full edit-to-preview journey.
