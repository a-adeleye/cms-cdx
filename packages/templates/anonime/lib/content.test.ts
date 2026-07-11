import assert from 'node:assert/strict';
import test from 'node:test';

import { parseArticleMarkdown } from './markdown.ts';
import {
  adaptArticle,
  normalizeBlogPath,
  normalizePositiveInteger,
  resolveArticles,
  resolveFeaturedArticle,
  resolveSite,
  safeMediaUrl,
} from './model.ts';

test('parseArticleMarkdown keeps raw HTML as inert text', () => {
  const parsed = parseArticleMarkdown('Hello <script>alert(1)</script> [open](javascript:alert(1))');
  const paragraph = parsed.blocks[0];

  assert.equal(paragraph.kind, 'paragraph');
  const children = paragraph.kind === 'paragraph' ? paragraph.children : [];
  assert.equal(children.every((node) => node.kind === 'text'), true);
  assert.equal(children.map((node) => node.kind === 'text' ? node.value : '').join(''), 'Hello <script>alert(1)</script> open)');
});

test('parseArticleMarkdown preserves supported semantic formatting', () => {
  const parsed = parseArticleMarkdown('## A useful heading\n\n**Bold** and [safe](https://example.com).');

  assert.equal(parsed.blocks[0].kind, 'heading');
  assert.equal(parsed.blocks[1].kind, 'paragraph');
  assert.deepEqual(parsed.tableOfContents, [
    { depth: 2, label: 'A useful heading', id: 'a-useful-heading' },
  ]);
});

test('parseArticleMarkdown makes duplicate heading anchors deterministic', () => {
  const parsed = parseArticleMarkdown('## Repeated\n\n## Repeated');

  assert.deepEqual(parsed.tableOfContents.map((entry) => entry.id), ['repeated', 'repeated-2']);
});

test('parseArticleMarkdown omits a Markdown H1 that duplicates the page title', () => {
  const parsed = parseArticleMarkdown('# Article title\n\nOpening paragraph.', 'Article title');

  assert.equal(parsed.blocks.length, 1);
  assert.equal(parsed.blocks[0].kind, 'paragraph');
});

test('parseArticleMarkdown rejects NUL characters without echoing article content', () => {
  assert.throws(() => parseArticleMarkdown('private\0content'), /unsupported control character/);
});

test('normalizeBlogPath confines navigation to a root-relative path', () => {
  assert.equal(normalizeBlogPath('/privacy/articles/'), '/privacy/articles');
  assert.equal(normalizeBlogPath('javascript:alert(1)'), '/articles');
  assert.equal(normalizeBlogPath('//example.com/articles'), '/articles');
  assert.equal(normalizeBlogPath('/'), '/articles');
});

test('safeMediaUrl allows web and root-relative media only', () => {
  assert.equal(safeMediaUrl('https://cdn.example.com/image.png'), 'https://cdn.example.com/image.png');
  assert.equal(safeMediaUrl('/media/image.png'), '/media/image.png');
  assert.equal(safeMediaUrl('data:text/html,unsafe'), undefined);
  assert.equal(safeMediaUrl('https://user:secret@example.com/image.png'), undefined);
});

test('resolveArticles distinguishes preview fixtures from a real empty collection', () => {
  assert.ok(resolveArticles().length > 0);
  assert.deepEqual(resolveArticles([]), []);
});

test('resolveArticles rejects malformed CMS payloads at the template boundary', () => {
  assert.throws(() => resolveArticles('not-an-array'), /provided as an array/);
  assert.throws(() => resolveArticles([{ title: 42, slug: 'invalid-title' }]), /title must be a string/);
  assert.throws(() => resolveSite({ name: 42 }), /name must be a string/);
});

test('adaptArticle rejects unsafe slugs instead of publishing a fixture route', () => {
  assert.throws(
    () => adaptArticle({ title: 'A new article', slug: '../unsafe' }, 0),
    /slug must contain lowercase letters/,
  );
  assert.throws(
    () => adaptArticle({ title: 'Uppercase route', slug: 'Uppercase-route' }, 0),
    /slug must contain lowercase letters/,
  );
});

test('resolveFeaturedArticle requires content for a CMS read page', () => {
  assert.throws(
    () => resolveFeaturedArticle({ title: 'Empty article', slug: 'empty-article', contentMarkdown: '' }),
    /contentMarkdown is required/,
  );
});

test('adaptArticle derives neutral metadata and reading time for valid CMS content', () => {
  const article = adaptArticle(
    {
      title: 'A new article',
      slug: 'a-new-article',
      contentMarkdown: Array.from({ length: 221 }, () => 'word').join(' '),
      publishedAt: '2026-02-10',
    },
    0,
  );

  assert.equal(article.slug, 'a-new-article');
  assert.equal(article.readingTime, 2);
  assert.equal(article.publishedLabel, 'Feb 10, 2026');
  assert.equal(article.author, 'Editorial team');
  assert.equal(article.category, 'Privacy');
});

test('normalizePositiveInteger rejects unsafe pagination values', () => {
  assert.equal(normalizePositiveInteger(undefined, 1), 1);
  assert.equal(normalizePositiveInteger(3, 1), 3);
  assert.throws(() => normalizePositiveInteger('3', 1), /positive integers/);
  assert.throws(() => normalizePositiveInteger(0, 1), /positive integers/);
});
