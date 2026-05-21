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
docker compose up --build
```

## Notes

- The CMS is the source of truth.
- Astro is only the static generator.
- Public sites are static and SEO-friendly.
- AI-assisted content is saved as draft or review only.
