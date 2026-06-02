INSERT INTO users (email, password_hash, full_name, role)
VALUES ('admin@example.com', crypt('admin123', gen_salt('bf')), 'Admin User', 'admin')
ON CONFLICT (email) DO NOTHING;

INSERT INTO sites (name, slug, domain, blog_path, template_key, status, deploy_provider, preview_deploy_provider)
VALUES ('Example Site', 'example', 'http://localhost:4321', '/articles', 'default-blog', 'active', 'none', 'none')
ON CONFLICT (slug) DO NOTHING;

