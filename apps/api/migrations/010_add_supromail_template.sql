INSERT INTO templates (name, slug)
VALUES ('Supromail', 'supromail')
ON CONFLICT (slug) DO UPDATE
SET name = EXCLUDED.name,
    updated_at = NOW();
