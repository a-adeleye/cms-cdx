# ADR 002: Supromail Astro template builds

Status: accepted (2026-07-28)

## Context

The CMS catalogue previously contained only templates rendered by the API's Go `LocalBuilder`. The repository also contained an Astro builder and a componentized Supromail template package, but neither was used by the API build/deploy workflow. As a result, selecting an Astro template could not create deployable CMS output.

## Decision

Supromail is registered as a trusted, repository-bundled template. The API selects an `AstroBuilder` only when a site's `template_key` is `supromail`; all existing template keys continue through `LocalBuilder` unchanged.

The adapter writes a typed JSON build-data file, starts `npm run build` in the checked-in Astro builder with a two-minute deadline, and passes only the environment values needed by Node. It then relocates Astro's fixed `/articles` output to the site's canonical `blog_path`, allowing Supromail pages and links to honor `/blog`, `/articles`, or another valid path. The generated output is passed to the existing deployment adapters without a new deployment contract.

The public catalogue preview builds a deterministic Supromail sample once per API process, caches the resulting HTML, and applies a restrictive content-security policy. Astro styles are inlined so the preview does not expose a static-asset serving surface.

The accepted template contract is:

- `site`: `name`, `domain`, `blogPath`
- `articles`: article title, slug, excerpt, Markdown body, SEO fields, category, author, date, cover URL, and featured flag
- entrypoints: `landing.astro`, `articles.astro`, and `article.astro`

## Alternatives considered

- Convert Supromail to Go string rendering. Rejected because it would duplicate an already-tested Astro component package and keep each new template as a bespoke backend implementation.
- Make the API execute arbitrary uploaded template code. Rejected for this iteration: untrusted build scripts and dependency installation require isolated build workers, quotas, malware/dependency scanning, and an upload/approval lifecycle.
- Replace every existing template with Astro immediately. Rejected because it enlarges the migration blast radius; the adapter allows incremental adoption and straightforward rollback.

## Consequences

- The API image now includes the trusted Supromail package under the Astro builder directory.
- A Supromail build requires the shipped Node/Astro runtime and can fail or time out independently of the Go renderer; failures are recorded in the existing build history.
- Template uploads are still not supported. This ADR establishes the package and content contract required for that later feature.

## Rollback

To stop using Supromail, assign affected sites an existing Go-rendered template and deploy again. A follow-up migration may remove the catalogue row only after no sites reference it. Reverting this change restores the former Go-only builder without changing existing site content or deployment records.
