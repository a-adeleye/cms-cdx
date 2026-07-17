ALTER TABLE sites
  ADD COLUMN IF NOT EXISTS content_context TEXT NOT NULL DEFAULT 'standalone_blog';

ALTER TABLE sites
  DROP CONSTRAINT IF EXISTS sites_content_context_check;

ALTER TABLE sites
  ADD CONSTRAINT sites_content_context_check
  CHECK (content_context IN ('application_blog', 'standalone_blog'));
