ALTER TABLE sites
  ADD COLUMN IF NOT EXISTS preview_deploy_provider TEXT,
  ADD COLUMN IF NOT EXISTS preview_deploy_config JSONB NOT NULL DEFAULT '{}';
