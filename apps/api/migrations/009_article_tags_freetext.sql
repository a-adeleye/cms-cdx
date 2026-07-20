ALTER TABLE articles ADD COLUMN tags TEXT NOT NULL DEFAULT '';

UPDATE articles a
SET tags = COALESCE((
  SELECT string_agg(t.name, ', ' ORDER BY t.name)
  FROM article_tags at
  JOIN tags t ON t.id = at.tag_id
  WHERE at.article_id = a.id
), '');

DROP TABLE article_tags;

DROP TABLE tags;
