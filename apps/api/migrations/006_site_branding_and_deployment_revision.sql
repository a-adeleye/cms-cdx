ALTER TABLE sites
  ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS logo_media_id UUID REFERENCES media_assets(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS favicon_media_id UUID REFERENCES media_assets(id) ON DELETE SET NULL;

ALTER TABLE builds
  ADD COLUMN IF NOT EXISTS deploy_revision TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS builds_site_created_at_idx ON builds(site_id, created_at DESC);

UPDATE sites SET deploy_provider = 'none', deploy_config = '{}'
WHERE COALESCE(deploy_provider, 'none') NOT IN ('none', 'firebase', 'cloudflare_pages', 'git_repository');
UPDATE sites SET preview_deploy_provider = 'none', preview_deploy_config = '{}'
WHERE COALESCE(preview_deploy_provider, 'none') NOT IN ('none', 'firebase', 'cloudflare_pages', 'git_repository');
