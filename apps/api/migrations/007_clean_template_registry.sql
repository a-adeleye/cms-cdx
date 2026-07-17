UPDATE sites SET template_key = 'default-blog'
WHERE template_key NOT IN ('default-blog', 'premium-saas', 'anonime');

DELETE FROM templates
WHERE slug NOT IN ('default-blog', 'premium-saas', 'anonime');

