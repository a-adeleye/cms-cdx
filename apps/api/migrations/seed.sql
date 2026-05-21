-- Migration Seed: Insert initial seed data for example site and user
-- The default user password is 'admin123' hashed with bcrypt: $2a$10$tZ27d0sQn4Fz.7JpL.LzueEszf1yX2rWv6lJ.x8sD7Nmsv0z5Xk0q

-- 1. Insert Admin User
INSERT INTO users (id, email, password_hash, full_name, role)
VALUES (
  'a29e92ff-08b5-4b0e-9279-d2d0b5e90d0b',
  'admin@example.com',
  '$2a$10$tZ27d0sQn4Fz.7JpL.LzueEszf1yX2rWv6lJ.x8sD7Nmsv0z5Xk0q',
  'Administrator',
  'admin'
) ON CONFLICT (email) DO NOTHING;

-- 2. Insert Example Site
INSERT INTO sites (id, name, slug, domain, blog_path, status, template_key, theme_config, deploy_provider, deploy_config, ai_config, storage_config)
VALUES (
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  'Example Site',
  'example',
  'http://localhost:4321',
  '/articles',
  'active',
  'default-blog',
  '{"primaryColor": "#0058be", "fontFamily": "Hanken Grotesk", "theme": "light"}',
  'none',
  '{}',
  '{"provider": "gemini", "model": "gemini-1.5-flash", "tone": "clear, premium, simple"}',
  '{"provider": "local"}'
) ON CONFLICT (slug) DO NOTHING;

-- 3. Insert Landing Page
INSERT INTO pages (id, site_id, type, title, slug, content_json, seo_title, seo_description, status, published_at)
VALUES (
  'c38e93ff-28b5-4b0e-9279-d2d0b5e90d0d',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  'landing',
  'Home',
  '/',
  '{"headline": "Welcome to the Future of Publishing", "subheadline": "Self-hosted, static-compiled, and AI-powered."}',
  'Example Site - Next-Gen Publishing Platform',
  'A self-hosted static site generated with Astro and powered by a Go API.',
  'published',
  NOW()
);

-- 4. Insert Landing Page Sections
-- Hero Section
INSERT INTO landing_sections (id, site_id, section_key, title, subtitle, content_json, display_order, is_enabled)
VALUES (
  'd48e93ff-38b5-4b0e-9279-d2d0b5e90d0e',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  'hero',
  'Supercharge Your Static Site Deployment',
  'The speed of static pages. The convenience of a modern CMS. Fully self-hosted.',
  '{"cta_text": "Get Started", "cta_url": "/articles", "image_url": "https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800"}',
  1,
  TRUE
);

-- Latest Articles Section
INSERT INTO landing_sections (id, site_id, section_key, title, subtitle, content_json, display_order, is_enabled)
VALUES (
  'e58e93ff-48b5-4b0e-9279-d2d0b5e90d0f',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  'latest_articles',
  'From Our Blog',
  'Stay up to date with the latest architectural practices and insights.',
  '{"limit": 3}',
  2,
  TRUE
);

-- CTA Section
INSERT INTO landing_sections (id, site_id, section_key, title, subtitle, content_json, display_order, is_enabled)
VALUES (
  'f68e93ff-58b5-4b0e-9279-d2d0b5e90d10',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  'cta',
  'Ready to Build Your Own Self-Hosted Network?',
  'Get complete control over your content, assets, and static pipelines today.',
  '{"button_text": "Read the Articles", "button_link": "/articles"}',
  3,
  TRUE
);

-- 5. Insert Author
INSERT INTO authors (id, site_id, name, slug, bio, avatar_url)
VALUES (
  '018e93ff-68b5-4b0e-9279-d2d0b5e90d11',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  'Adeleye Adeyemi',
  'adeleye-adeyemi',
  'Lead Cloud Architect & Developer. Obsessed with high-performance static sites, Golang systems, and clean code.',
  'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150'
) ON CONFLICT (site_id, slug) DO NOTHING;

-- 6. Insert Categories
INSERT INTO categories (id, site_id, name, slug, description)
VALUES (
  '028e93ff-78b5-4b0e-9279-d2d0b5e90d12',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  'Technology',
  'technology',
  'The latest trends in cloud engineering, frameworks, and architecture.'
) ON CONFLICT (site_id, slug) DO NOTHING;

INSERT INTO categories (id, site_id, name, slug, description)
VALUES (
  '038e93ff-88b5-4b0e-9279-d2d0b5e90d13',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  'Design Systems',
  'design-systems',
  'Best practices for modern, responsive, and beautiful layouts.'
) ON CONFLICT (site_id, slug) DO NOTHING;

-- 7. Insert Tags
INSERT INTO tags (id, site_id, name, slug)
VALUES (
  '048e93ff-98b5-4b0e-9279-d2d0b5e90d14',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  'Golang',
  'golang'
) ON CONFLICT (site_id, slug) DO NOTHING;

INSERT INTO tags (id, site_id, name, slug)
VALUES (
  '058e93ff-a8b5-4b0e-9279-d2d0b5e90d15',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  'Astro',
  'astro'
) ON CONFLICT (site_id, slug) DO NOTHING;

INSERT INTO tags (id, site_id, name, slug)
VALUES (
  '068e93ff-b8b5-4b0e-9279-d2d0b5e90d16',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  'Vanilla CSS',
  'vanilla-css'
) ON CONFLICT (site_id, slug) DO NOTHING;

-- 8. Insert Articles
-- Article 1: Golang Power
INSERT INTO articles (id, site_id, author_id, category_id, title, slug, excerpt, content_markdown, cover_image_url, status, is_featured, published_at, seo_title, seo_description, canonical_url, generated_by_ai, human_reviewed)
VALUES (
  '078e93ff-c8b5-4b0e-9279-d2d0b5e90d17',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  '018e93ff-68b5-4b0e-9279-d2d0b5e90d11',
  '028e93ff-78b5-4b0e-9279-d2d0b5e90d12',
  'Building Lightning-Fast Backends with Go',
  'building-lightning-fast-backends-with-go',
  'Learn how Go standard library and custom routers allow you to achieve sub-millisecond API response times.',
  '# Building Lightning-Fast Backends with Go

Golang has emerged as the language of choice for building scalable, high-performance web APIs. Its simplicity, powerful concurrency model, and lightning-fast compilation make it extremely suitable for modern backend development.

## Why Go is Ideal for CMS Builders

When constructing a static site builder, your API needs to efficiently fetch data from PostgreSQL, parse schemas, execute build steps, and communicate with S3. Go enables you to run all these processes concurrently with minimal memory footprints.

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from the Go API!")
    })
    http.ListenAndServe(":8080", nil)
}
```

## Concurrency via Goroutines

Unlike other languages that rely heavily on complex async event loops or multi-threading pools, Go leverages simple goroutines. A goroutine is a lightweight thread managed by the Go runtime, requiring only a few kilobytes of initial stack space.

We use this concurrency to spawn background static compiler tasks in our CMS builder without blocking incoming API calls.',
  'https://images.unsplash.com/photo-1618401471353-b98aedd07871?w=800',
  'published',
  TRUE,
  NOW(),
  'Go Performance Backend Guide - CMS Builder',
  'How to build highly efficient Go microservices for content orchestration.',
  'http://localhost:4321/articles/building-lightning-fast-backends-with-go',
  FALSE,
  TRUE
) ON CONFLICT (site_id, slug) DO NOTHING;

-- Article 2: Astro Static Site Builder
INSERT INTO articles (id, site_id, author_id, category_id, title, slug, excerpt, content_markdown, cover_image_url, status, is_featured, published_at, seo_title, seo_description, canonical_url, generated_by_ai, human_reviewed)
VALUES (
  '088e93ff-d8b5-4b0e-9279-d2d0b5e90d18',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  '018e93ff-68b5-4b0e-9279-d2d0b5e90d11',
  '028e93ff-78b5-4b0e-9279-d2d0b5e90d12',
  'Astro: The Next-Gen Static Generator',
  'astro-the-next-gen-static-generator',
  'An in-depth look at Astro''s island architecture and static generation capabilities for fast SEO ranking.',
  '# Astro: The Next-Gen Static Generator

Astro is a modern web framework designed for content-focused websites. By delivering zero client-side JavaScript by default, it makes pages load instantly and ranks highly on search engines.

## Island Architecture

Astro''s primary innovation is "Island Architecture". Instead of bundling your entire application into a massive single-page application (SPA), Astro compiles static HTML and allows you to inject interactive "islands" of React, Vue, or Angular only where dynamic client logic is needed.

## Re-building Static Files

In our CMS flow, Astro functions as a pure static compiler. It pulls structural JSON from our Go API, feeds it into templates, and writes static static assets into a build folder ready for direct S3 or SFTP deployment.

This represents the ultimate separation of concerns: the CMS handles data state, while Astro processes presentation without runtime overhead.',
  'https://images.unsplash.com/photo-1507238691740-187a5b1d37b8?w=800',
  'published',
  FALSE,
  NOW(),
  'Astro Static Site Framework Deep-Dive',
  'Discover why zero-JS templates and island rendering are excellent for SEO content.',
  'http://localhost:4321/articles/astro-the-next-gen-static-generator',
  FALSE,
  TRUE
) ON CONFLICT (site_id, slug) DO NOTHING;

-- Article 3: Premium Styling with Vanilla CSS
INSERT INTO articles (id, site_id, author_id, category_id, title, slug, excerpt, content_markdown, cover_image_url, status, is_featured, published_at, seo_title, seo_description, canonical_url, generated_by_ai, human_reviewed)
VALUES (
  '098e93ff-e8b5-4b0e-9279-d2d0b5e90d19',
  'b28e93ff-18b5-4b0e-9279-d2d0b5e90d0c',
  '018e93ff-68b5-4b0e-9279-d2d0b5e90d11',
  '038e93ff-88b5-4b0e-9279-d2d0b5e90d13',
  'Crafting Premium Designs with Vanilla CSS',
  'crafting-premium-designs-with-vanilla-css',
  'Unlock modern CSS features like Custom Properties, Grid, Flexbox, and complex animations without bloated libraries.',
  '# Crafting Premium Designs with Vanilla CSS

Many developers immediately reach for bulky styling frameworks. However, vanilla CSS has evolved into a robust, elegant, and standard styling mechanism that offers ultimate performance and total design control.

## Leveraging CSS Variables

By establishing design tokens as CSS custom properties, you can easily implement light/dark mode and switch brand themes on the fly.

```css
:root {
  --primary-color: #0058be;
  --surface-bg: #f7f9fb;
  --font-family: ''Hanken Grotesk'', sans-serif;
  --rounded-md: 0.5rem;
}
```

## Layout Control with CSS Grid

Grid and Flexbox allow you to build fluid layouts that automatically snap to tablet and mobile viewports with extremely minimal code. 

By avoiding heavy utility-class layers, your generated static pages remain highly readable, SEO-compliant, and exceptionally lightweight.',
  'https://images.unsplash.com/photo-1550751827-4bd374c3f58b?w=800',
  'published',
  FALSE,
  NOW(),
  'Vanilla CSS Design Token Customization',
  'How to build flexible components with native CSS custom properties.',
  'http://localhost:4321/articles/crafting-premium-designs-with-vanilla-css',
  FALSE,
  TRUE
) ON CONFLICT (site_id, slug) DO NOTHING;

-- 9. Insert Article Tags
INSERT INTO article_tags (article_id, tag_id) VALUES ('078e93ff-c8b5-4b0e-9279-d2d0b5e90d17', '048e93ff-98b5-4b0e-9279-d2d0b5e90d14');
INSERT INTO article_tags (article_id, tag_id) VALUES ('078e93ff-c8b5-4b0e-9279-d2d0b5e90d17', '058e93ff-a8b5-4b0e-9279-d2d0b5e90d15');
INSERT INTO article_tags (article_id, tag_id) VALUES ('088e93ff-d8b5-4b0e-9279-d2d0b5e90d18', '058e93ff-a8b5-4b0e-9279-d2d0b5e90d15');
INSERT INTO article_tags (article_id, tag_id) VALUES ('098e93ff-e8b5-4b0e-9279-d2d0b5e90d19', '068e93ff-b8b5-4b0e-9279-d2d0b5e90d16');
