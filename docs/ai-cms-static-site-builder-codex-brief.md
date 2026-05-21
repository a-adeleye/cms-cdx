# AI CMS Static Site Builder — Codex Implementation Brief

## Goal

Build a reusable, self-hosted CMS and static site builder that can manage multiple websites from one admin system.

Each website should have:

- Its own landing page template
- Its own articles/blog listing page
- Its own article read/details page
- Its own theme/config
- Its own deploy target
- Dynamic content managed from the CMS
- Static HTML/CSS/JS output for SEO
- AI-assisted content generation

The system should allow one CMS to power many different landing pages and blogs, while deploying each generated website to its own hosting provider.

---

## Core Concept

The CMS is the source of truth.

Astro is only the static page generator.

The builder reads content and configuration from the CMS/database, renders the selected templates, generates static files, and deploys them to the configured hosting target.

Content updates should trigger a rebuild of the affected site so pages like the landing page can automatically update recent articles, featured articles, and other dynamic sections.

---

## Required Stack

Use the following stack:

```txt
Backend API: Go
Database: PostgreSQL
Admin Dashboard: Angular or React
Static Builder: Astro
Storage: S3-compatible storage
Local Storage: MinIO
Production Storage: Cloudflare R2 or any S3-compatible provider
AI Providers: OpenAI, Gemini, Anthropic, configurable
Deployment Providers: Firebase Hosting, Cloudflare Pages, Netlify, SFTP/VPS
Containerization: Docker + Docker Compose
```

The system must not depend on Sanity SaaS or any hosted CMS.

This should be a self-hosted CMS and static site builder.

---

## High-Level Architecture

```txt
Admin Dashboard
  ↓
Go API
  ↓
PostgreSQL
  ↓
Build Service
  ↓
Astro Templates
  ↓
Static HTML/CSS/JS Output
  ↓
Deploy Adapter
  ↓
Target Hosting
```

The app should support many sites from one CMS.

Example:

```txt
CMS
 ├─ anonime.io
 ├─ leatra.com
 ├─ tradersdojo.com
 └─ client-site.com
```

Each site has its own content, pages, articles, theme, and deployment configuration.

---

## Monorepo Structure

Use this structure:

```txt
cms-builder/
  apps/
    api/
      cmd/
      internal/
        auth/
        config/
        database/
        handlers/
        middleware/
        models/
        repositories/
        services/
        storage/
        ai/
        builder/
        deploy/
      migrations/
      Dockerfile

    admin/
      src/
      Dockerfile

    builder/
      astro/
        src/
          pages/
            index.astro
            articles/
              index.astro
              [slug].astro
          layouts/
          components/
          styles/
        public/
        package.json
        astro.config.mjs

  packages/
    templates/
      default-blog/
        landing.astro
        articles.astro
        article.astro
        components/
        styles/

      premium-saas/
        landing.astro
        articles.astro
        article.astro
        components/
        styles/

    deploy-adapters/
      firebase/
      cloudflare/
      netlify/
      sftp/

  sites/
    example/
      site.config.json

  docker-compose.yml
  README.md
```

---

## Main Features

Build the system with these features:

```txt
Authentication
Site management
Page management
Landing page section management
Article/blog management
Authors
Categories
Tags
Media uploads
Theme configuration
Template selection
AI content assistant
Static site generation
Build logs
Deployment adapters
Webhook-triggered rebuilds
Manual rebuild button
Preview build
Published build
```

---

## Database Schema

Use PostgreSQL.

Create the following tables.

---

### users

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  full_name TEXT,
  role TEXT NOT NULL DEFAULT 'admin',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

### sites

```sql
CREATE TABLE sites (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  domain TEXT,
  blog_path TEXT NOT NULL DEFAULT '/articles',
  status TEXT NOT NULL DEFAULT 'active',
  template_key TEXT NOT NULL DEFAULT 'default-blog',
  theme_config JSONB NOT NULL DEFAULT '{}',
  deploy_provider TEXT,
  deploy_config JSONB NOT NULL DEFAULT '{}',
  ai_config JSONB NOT NULL DEFAULT '{}',
  storage_config JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

### pages

```sql
CREATE TABLE pages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  site_id UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  title TEXT NOT NULL,
  slug TEXT,
  content_json JSONB NOT NULL DEFAULT '{}',
  seo_title TEXT,
  seo_description TEXT,
  status TEXT NOT NULL DEFAULT 'draft',
  published_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Allowed page types:

```txt
landing
articles
custom
```

The article details page does not need a page row for every article. It should be generated from the articles table.

---

### landing_sections

```sql
CREATE TABLE landing_sections (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  site_id UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  section_key TEXT NOT NULL,
  title TEXT,
  subtitle TEXT,
  content_json JSONB NOT NULL DEFAULT '{}',
  display_order INTEGER NOT NULL DEFAULT 0,
  is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Example section keys:

```txt
hero
features
how_it_works
latest_articles
featured_articles
cta
faq
footer
```

The latest_articles and featured_articles sections should be dynamic and should pull from the articles table during build.

---

### articles

```sql
CREATE TABLE articles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  site_id UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  author_id UUID,
  category_id UUID,
  title TEXT NOT NULL,
  slug TEXT NOT NULL,
  excerpt TEXT,
  content_markdown TEXT NOT NULL,
  cover_image_url TEXT,
  status TEXT NOT NULL DEFAULT 'draft',
  is_featured BOOLEAN NOT NULL DEFAULT FALSE,
  published_at TIMESTAMPTZ,
  seo_title TEXT,
  seo_description TEXT,
  canonical_url TEXT,
  generated_by_ai BOOLEAN NOT NULL DEFAULT FALSE,
  human_reviewed BOOLEAN NOT NULL DEFAULT FALSE,
  ai_prompt TEXT,
  ai_model TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(site_id, slug)
);
```

Allowed statuses:

```txt
draft
review
published
archived
```

Only published articles should appear in generated pages.

---

### authors

```sql
CREATE TABLE authors (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  site_id UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  slug TEXT NOT NULL,
  bio TEXT,
  avatar_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(site_id, slug)
);
```

---

### categories

```sql
CREATE TABLE categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  site_id UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  slug TEXT NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(site_id, slug)
);
```

---

### tags

```sql
CREATE TABLE tags (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  site_id UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  slug TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(site_id, slug)
);
```

---

### article_tags

```sql
CREATE TABLE article_tags (
  article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
  tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (article_id, tag_id)
);
```

---

### media_assets

```sql
CREATE TABLE media_assets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  site_id UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  file_name TEXT NOT NULL,
  file_url TEXT NOT NULL,
  mime_type TEXT,
  size_bytes BIGINT,
  storage_provider TEXT NOT NULL DEFAULT 's3',
  storage_key TEXT,
  alt_text TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

### builds

```sql
CREATE TABLE builds (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  site_id UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  status TEXT NOT NULL DEFAULT 'queued',
  build_type TEXT NOT NULL DEFAULT 'published',
  logs TEXT,
  output_path TEXT,
  deploy_provider TEXT,
  deploy_status TEXT,
  deploy_url TEXT,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Allowed build statuses:

```txt
queued
running
success
failed
```

Allowed build types:

```txt
preview
published
```

---

## API Requirements

Create REST APIs.

Base path:

```txt
/api/v1
```

---

### Auth

```txt
POST /auth/login
POST /auth/logout
GET  /auth/me
```

Use JWT auth.

---

### Sites

```txt
GET    /sites
POST   /sites
GET    /sites/:id
PUT    /sites/:id
DELETE /sites/:id
```

---

### Pages

```txt
GET    /sites/:siteId/pages
POST   /sites/:siteId/pages
GET    /pages/:id
PUT    /pages/:id
DELETE /pages/:id
```

---

### Landing Sections

```txt
GET    /sites/:siteId/sections
POST   /sites/:siteId/sections
PUT    /sections/:id
DELETE /sections/:id
```

---

### Articles

```txt
GET    /sites/:siteId/articles
POST   /sites/:siteId/articles
GET    /articles/:id
PUT    /articles/:id
DELETE /articles/:id
POST   /articles/:id/publish
POST   /articles/:id/unpublish
```

---

### Authors

```txt
GET    /sites/:siteId/authors
POST   /sites/:siteId/authors
PUT    /authors/:id
DELETE /authors/:id
```

---

### Categories

```txt
GET    /sites/:siteId/categories
POST   /sites/:siteId/categories
PUT    /categories/:id
DELETE /categories/:id
```

---

### Tags

```txt
GET    /sites/:siteId/tags
POST   /sites/:siteId/tags
PUT    /tags/:id
DELETE /tags/:id
```

---

### Media

```txt
POST   /sites/:siteId/media
GET    /sites/:siteId/media
DELETE /media/:id
```

Media upload should use S3-compatible storage.

Use MinIO locally and allow switching to Cloudflare R2 or another S3 provider in production.

---

### AI Content Assistant

```txt
POST /sites/:siteId/ai/generate-ideas
POST /sites/:siteId/ai/generate-outline
POST /sites/:siteId/ai/generate-draft
POST /sites/:siteId/ai/generate-seo
POST /sites/:siteId/ai/generate-tags
POST /sites/:siteId/ai/generate-image-prompt
```

AI should never publish content automatically.

Generated content should be saved as draft or review only.

---

### Build and Deploy

```txt
POST /sites/:siteId/build
GET  /sites/:siteId/builds
GET  /builds/:id
POST /builds/:id/deploy
```

Build request body:

```json
{
  "buildType": "published",
  "deploy": true
}
```

---

## Admin Dashboard Requirements

Build an admin dashboard with the following sections:

```txt
Login
Dashboard
Sites
Site Settings
Landing Page Editor
Articles
Article Editor
Authors
Categories
Tags
Media Library
AI Assistant
Builds
Deployment Settings
```

---

## Admin UX Requirements

The article editor should support:

```txt
Title
Slug
Excerpt
Cover image
Markdown body
SEO title
SEO description
Canonical URL
Author
Category
Tags
Featured toggle
Status
AI generation panel
Preview button
Save draft
Mark reviewed
Publish
```

The landing page editor should support:

```txt
Manage sections
Reorder sections
Enable/disable sections
Edit section content JSON
Preview section output
```

The site settings page should support:

```txt
Domain
Blog path
Template
Theme config
Storage config
Deployment config
AI provider config
```

---

## Astro Builder Requirements

The builder should generate static HTML for each site.

Required pages:

```txt
/
{blog_path}/
{blog_path}/{article_slug}/
```

Example:

```txt
/
/articles/
/articles/what-is-private-email/
```

or:

```txt
/
/blog/
/blog/what-is-private-email/
```

depending on the site config.

The builder must:

```txt
Load site config
Load landing page content
Load enabled landing sections
Load published articles
Load featured articles
Load categories and tags
Render landing page
Render articles list page
Render article read pages
Generate sitemap.xml
Generate robots.txt
Generate RSS feed
Generate OpenGraph metadata
Generate canonical URLs
Generate static assets
Output to dist/sites/{site_slug}
```

---

## Landing Page Behavior

The landing page should be template-driven.

It must be able to render static sections and dynamic sections.

Example static sections:

```txt
hero
features
how_it_works
faq
cta
footer
```

Example dynamic sections:

```txt
latest_articles
featured_articles
category_articles
```

When a new article is published, the landing page should be rebuilt so recent article sections update automatically.

---

## Article List Page Behavior

The articles page should show:

```txt
Page title
Page description
Featured article
Recent articles
Category filter links
Article cards
Pagination or load more
```

For static output, prefer paginated static pages:

```txt
/articles/
/articles/page/2/
/articles/page/3/
```

---

## Article Read Page Behavior

Each article read page should include:

```txt
Title
Excerpt
Cover image
Author
Published date
Category
Tags
Markdown content rendered as HTML
Table of contents if possible
Related articles
SEO metadata
OpenGraph metadata
Canonical URL
```

---

## Markdown Rendering

Articles are stored in the database as Markdown.

During Astro build, Markdown should be rendered to HTML.

Support:

```txt
Headings
Paragraphs
Links
Images
Code blocks
Blockquotes
Tables
Lists
Embeds if safe
```

Sanitize rendered HTML before output.

---

## SEO Requirements

Every generated page should include:

```txt
<title>
<meta name="description">
Canonical URL
OpenGraph title
OpenGraph description
OpenGraph image
Twitter card metadata
Structured data where useful
Sitemap
Robots.txt
RSS feed
```

Generate structured data for articles using JSON-LD:

```txt
BlogPosting
```

---

## AI Assistant Requirements

AI should be provider-based.

Create interface:

```go
type AIProvider interface {
    GenerateIdeas(input GenerateIdeasInput) ([]Idea, error)
    GenerateOutline(input GenerateOutlineInput) (Outline, error)
    GenerateDraft(input GenerateDraftInput) (Draft, error)
    GenerateSEO(input GenerateSEOInput) (SEOResult, error)
    GenerateTags(input GenerateTagsInput) ([]string, error)
    GenerateImagePrompt(input GenerateImagePromptInput) (string, error)
}
```

Supported providers:

```txt
openai
gemini
anthropic
none
```

The site should store AI configuration in `sites.ai_config`.

Example:

```json
{
  "provider": "openai",
  "model": "gpt-4.1-mini",
  "tone": "clear, premium, simple",
  "brand_context": "Privacy-focused app for personas, private drops, encrypted temporary conversations."
}
```

AI output must be saved as draft content.

Never publish without human review.

---

## Deploy Adapter Requirements

Create a deploy interface:

```go
type DeployAdapter interface {
    Deploy(ctx context.Context, site Site, build Build, outputPath string) (*DeployResult, error)
}
```

Supported deploy providers:

```txt
firebase
cloudflare_pages
netlify
sftp
none
```

Each site should define its deploy config.

Example:

```json
{
  "provider": "firebase",
  "projectId": "anonime-landing",
  "siteId": "anonime",
  "tokenSecretRef": "FIREBASE_TOKEN"
}
```

Do not store raw deployment secrets in plain text in the database.

Use environment variables or secret references.

---

## Storage Requirements

Use S3-compatible storage.

Local:

```txt
MinIO
```

Production:

```txt
Cloudflare R2
AWS S3
Any S3-compatible storage
```

Create storage interface:

```go
type StorageProvider interface {
    Upload(ctx context.Context, file UploadFile) (*StoredFile, error)
    Delete(ctx context.Context, key string) error
    GetPublicURL(key string) string
}
```

---

## Docker Compose

Create docker-compose for local development:

```txt
postgres
minio
api
admin
builder
```

PostgreSQL should run with persistent volume.

MinIO should expose console and API ports.

---

## Environment Variables

Use environment variables:

```txt
DATABASE_URL=
JWT_SECRET=
API_PORT=

S3_ENDPOINT=
S3_REGION=
S3_BUCKET=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_PUBLIC_URL=

OPENAI_API_KEY=
GEMINI_API_KEY=
ANTHROPIC_API_KEY=

FIREBASE_TOKEN=
CLOUDFLARE_API_TOKEN=
NETLIFY_AUTH_TOKEN=
```

---

## Build Trigger Rules

Trigger rebuild when:

```txt
Article is published
Published article is updated
Published article is unpublished
Landing page content changes
Landing sections change
Site theme changes
Site template changes
Deploy config changes
```

Draft changes should not trigger a published rebuild.

---

## Preview Requirements

The system should support preview builds.

Preview builds should:

```txt
Include draft content
Generate to a temporary output folder
Not deploy to production unless explicitly requested
Return a preview URL if preview hosting is configured
```

---

## Security Requirements

Implement:

```txt
JWT authentication
Password hashing with bcrypt or argon2
Role-based permissions
Input validation
Markdown sanitization
File upload validation
Private deployment secrets
CORS config
Rate limiting for AI endpoints
Audit logs for publishing and deployment actions
```

---

## Audit Logs

Create audit log table:

```sql
CREATE TABLE audit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  site_id UUID REFERENCES sites(id),
  action TEXT NOT NULL,
  entity_type TEXT,
  entity_id UUID,
  metadata JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Log:

```txt
login
create_site
update_site
create_article
update_article
publish_article
unpublish_article
generate_ai_content
trigger_build
deploy_site
```

---

## Required Seed Data

Create seed data for one example site:

```txt
Site name: Example Site
Slug: example
Domain: http://localhost:4321
Blog path: /articles
Template: default-blog
```

Create:

```txt
Landing page
Hero section
Latest articles section
CTA section
One author
Two categories
Three tags
Three sample articles
```

---

## Acceptance Criteria

The project is complete when:

```txt
Admin can log in
Admin can create/edit a site
Admin can edit landing page sections
Admin can create/edit/publish articles
Admin can upload media
Admin can generate AI draft content
Admin can trigger a build
Builder generates static HTML files
Landing page includes latest articles
Articles page lists published articles
Article read page renders Markdown as HTML
Sitemap, robots.txt, RSS are generated
Deployment adapter can deploy or simulate deployment
System runs locally with Docker Compose
No Sanity SaaS dependency exists
```

---

## Important Implementation Rules

```txt
Do not use Sanity SaaS.
Do not use WordPress.
Do not make Astro the CMS.
Do not render the public blog with client-side JavaScript only.
Do not auto-publish AI content.
Do not store raw provider secrets in the database.
Do not hardcode one website into the system.
Do not build only for one hosting provider.
```

---

## Final Product Definition

The final product is a reusable self-hosted CMS and static site builder.

It should let the user create multiple websites, manage content for each one, generate SEO-friendly static HTML pages with Astro, and deploy those pages to different hosting providers using config.

The public sites should be fast, static, SEO-friendly, and independent of the CMS runtime after deployment.

The CMS should be reusable across many brands, landing pages, and blogs.

