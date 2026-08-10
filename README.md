# cms-cdx

Git-backed local CMS and Astro site builder, based on an MIT-licensed Pages CMS distribution.

## Local end-to-end workflow

```powershell
npm install
npm run dev:local
```

Open Pages CMS at `http://127.0.0.1:3000` and the live Astro preview at `http://127.0.0.1:4321`. Choose `local / cms-cdx`, open `working-tree`, and edit an article. Saving updates `apps/builder/content` directly; Astro refreshes from those Git-tracked files.

The default AI provider is a deterministic mock, so generate → review → save → preview works without credentials. To use OpenAI, set `AI_PROVIDER=openai`, `OPENAI_API_KEY`, and optionally `OPENAI_MODEL` before starting. The key remains server-side.

Local repository mode is deliberately development-only, bound to loopback, confined to the configured repository root, and unavailable from production builds.

## Structure

- `apps/pagescms`: pinned Pages CMS v2.1.8 distribution with local repository and AI extensions
- `.pages.yml`: editor schemas for site, landing page, articles, taxonomy, authors, and media
- `apps/builder/content`: canonical JSON content in Git
- `apps/builder`: Astro static site and live preview
- `apps/api` and `apps/admin`: previous implementation retained as a rollback path
- `packages/templates`: reusable site templates

## Commands

- `npm run dev:local`: start CMS and Astro preview together
- `npm run dev:cms`: start only Pages CMS on port 3000
- `npm run dev:builder`: start only Astro on port 4321
- `npm run build:cms`: build the Pages CMS distribution
- `npm run build:builder`: build the Astro site
- `npm run test:cms`: run local adapter and AI tests

## Optional OpenAI configuration

```powershell
$env:AI_PROVIDER = "openai"
$env:OPENAI_API_KEY = "..."
$env:OPENAI_MODEL = "gpt-5.6-luna"
npm run dev:local
```

The default mock provider is deterministic and intended for local workflow testing. AI output only populates editable form fields; it is never saved automatically.

## Previous stack

The Go API, Angular admin, Postgres, and MinIO compose services remain available for comparison and rollback. They are not the content source for the Pages CMS path.
