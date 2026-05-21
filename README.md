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

- Admin: `http://localhost:3000`
- API: `http://localhost:8081`
- Postgres: `localhost:5433`
- MinIO: `http://localhost:9002` and `http://localhost:9003`

## Notes

- The CMS is the source of truth.
- Astro is only the static generator.
- Public sites are static and SEO-friendly.
- AI-assisted content is saved as draft or review only.
