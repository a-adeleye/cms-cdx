CREATE TABLE IF NOT EXISTS templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO templates (name, slug)
VALUES
  ('Default Blog', 'default-blog'),
  ('Premium SaaS', 'premium-saas'),
  ('Anonime', 'anonime')
ON CONFLICT (slug) DO UPDATE
SET name = EXCLUDED.name,
    updated_at = NOW();
