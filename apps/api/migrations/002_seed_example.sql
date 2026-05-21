INSERT INTO sites (name, slug, domain, blog_path, template_key, status)
VALUES ('Example Site', 'example', 'http://localhost:4321', '/articles', 'default-blog', 'active')
ON CONFLICT (slug) DO NOTHING;

