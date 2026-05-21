INSERT INTO authors (site_id, name, slug, bio)
SELECT id, 'Ava Carter', 'ava-carter', 'Editor and strategist'
FROM sites
WHERE slug = 'example'
ON CONFLICT (site_id, slug) DO NOTHING;

INSERT INTO categories (site_id, name, slug, description)
SELECT id, 'Strategy', 'strategy', 'Planning and positioning'
FROM sites
WHERE slug = 'example'
ON CONFLICT (site_id, slug) DO NOTHING;

INSERT INTO categories (site_id, name, slug, description)
SELECT id, 'Product', 'product', 'Product and feature thinking'
FROM sites
WHERE slug = 'example'
ON CONFLICT (site_id, slug) DO NOTHING;

INSERT INTO tags (site_id, name, slug)
SELECT id, 'Privacy', 'privacy'
FROM sites
WHERE slug = 'example'
ON CONFLICT (site_id, slug) DO NOTHING;

INSERT INTO tags (site_id, name, slug)
SELECT id, 'SEO', 'seo'
FROM sites
WHERE slug = 'example'
ON CONFLICT (site_id, slug) DO NOTHING;

INSERT INTO tags (site_id, name, slug)
SELECT id, 'Static Sites', 'static-sites'
FROM sites
WHERE slug = 'example'
ON CONFLICT (site_id, slug) DO NOTHING;

INSERT INTO landing_sections (site_id, section_key, title, subtitle, display_order, is_enabled, content_json)
SELECT id, 'hero', 'Build once, deploy many', 'Static sites from one CMS', 0, TRUE, '{"headline":"CMS-backed static sites"}'::jsonb
FROM sites
WHERE slug = 'example'
ON CONFLICT (site_id, section_key) DO NOTHING;

INSERT INTO landing_sections (site_id, section_key, title, subtitle, display_order, is_enabled, content_json)
SELECT id, 'latest_articles', 'Latest Articles', 'Pulled from published content', 1, TRUE, '{}'::jsonb
FROM sites
WHERE slug = 'example'
ON CONFLICT (site_id, section_key) DO NOTHING;

INSERT INTO landing_sections (site_id, section_key, title, subtitle, display_order, is_enabled, content_json)
SELECT id, 'cta', 'Get Started', 'Example call to action', 2, TRUE, '{"buttonLabel":"Start building"}'::jsonb
FROM sites
WHERE slug = 'example'
ON CONFLICT (site_id, section_key) DO NOTHING;

INSERT INTO articles (
  site_id,
  author_id,
  category_id,
  title,
  slug,
  excerpt,
  content_markdown,
  status,
  is_featured,
  published_at,
  seo_title,
  seo_description,
  generated_by_ai,
  human_reviewed
)
SELECT
  s.id,
  a.id,
  c.id,
  'What is a private CMS?',
  'what-is-private-cms',
  'A practical introduction to self-hosted content systems.',
  '# What is a private CMS?\n\nA private CMS keeps content, workflow, and deployment under your control.',
  'published',
  TRUE,
  NOW(),
  'What is a private CMS?',
  'Understand private CMS architecture and tradeoffs.',
  FALSE,
  TRUE
FROM sites s
JOIN authors a ON a.site_id = s.id AND a.slug = 'ava-carter'
JOIN categories c ON c.site_id = s.id AND c.slug = 'strategy'
WHERE s.slug = 'example'
ON CONFLICT (site_id, slug) DO NOTHING;

INSERT INTO articles (
  site_id,
  author_id,
  category_id,
  title,
  slug,
  excerpt,
  content_markdown,
  status,
  is_featured,
  published_at,
  seo_title,
  seo_description,
  generated_by_ai,
  human_reviewed
)
SELECT
  s.id,
  a.id,
  c.id,
  'How static rendering helps SEO',
  'how-static-rendering-helps-seo',
  'Static HTML gives search engines predictable, fast pages.',
  '# How static rendering helps SEO\n\nStatic output reduces runtime dependency and improves crawlability.',
  'published',
  FALSE,
  NOW(),
  'How static rendering helps SEO',
  'Why static HTML still matters for discovery.',
  TRUE,
  FALSE
FROM sites s
JOIN authors a ON a.site_id = s.id AND a.slug = 'ava-carter'
JOIN categories c ON c.site_id = s.id AND c.slug = 'product'
WHERE s.slug = 'example'
ON CONFLICT (site_id, slug) DO NOTHING;

INSERT INTO articles (
  site_id,
  author_id,
  category_id,
  title,
  slug,
  excerpt,
  content_markdown,
  status,
  is_featured,
  published_at,
  seo_title,
  seo_description,
  generated_by_ai,
  human_reviewed
)
SELECT
  s.id,
  a.id,
  c.id,
  'Choosing deployment targets per site',
  'choosing-deployment-targets-per-site',
  'Each site can deploy to a different provider from one CMS.',
  '# Choosing deployment targets per site\n\nDeployment adapters keep provider logic isolated.',
  'published',
  FALSE,
  NOW(),
  'Choosing deployment targets per site',
  'Use deployment adapters to keep providers isolated.',
  FALSE,
  TRUE
FROM sites s
JOIN authors a ON a.site_id = s.id AND a.slug = 'ava-carter'
JOIN categories c ON c.site_id = s.id AND c.slug = 'strategy'
WHERE s.slug = 'example'
ON CONFLICT (site_id, slug) DO NOTHING;
