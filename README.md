# CMS Builder

Reusable self-hosted CMS and static site builder.

## Structure

- `apps/api`: Go REST API and build/deploy orchestration
- `apps/admin`: Angular admin dashboard
- `apps/builder`: Astro static site builder
- `packages/templates`: reusable site templates
- `packages/deploy-adapters`: deployment adapter definitions
- `sites/example`: seed config for the example site

## Local Development

```bash
npm install
npm run dev:admin
docker compose up --build
```

Local host ports:

- Admin: `http://localhost:9960`
- API: `http://localhost:8081`
- Postgres: `localhost:5433`
- MinIO: `http://localhost:9002` and `http://localhost:9003`

## Notes

- The CMS is the source of truth.
- Astro is only the static generator.
- Public sites are static and SEO-friendly.
- AI-assisted content is saved as draft or review only.
- Each site records whether it is an application blog or standalone blog, so future AI writing can apply the appropriate conversion-aware or editorial guidance.

## Cloudflare Pages Deployments

Published builds can deploy to a Cloudflare Pages Direct Upload project when the site's deployment provider is `cloudflare_pages`. See `packages/deploy-adapters/cloudflare/README.md` for the site configuration, runtime credentials, post-deploy checks, and rollback procedure.

## Repository Deployments

Use `git_repository` when an existing landing-page repository should keep control of the root site while the CMS publishes a nested blog. The CMS replaces only the configured directory, commits the generated blog, and pushes the configured preview or production branch. See `packages/deploy-adapters/repository/README.md` and `docs/adr/001-deployment-modes-and-site-branding.md`.
