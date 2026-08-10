import { mkdtemp, mkdir, rm, writeFile } from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';
import { afterEach, describe, expect, it } from 'vitest';
import { loadBuildData } from './content';

const roots: string[] = [];

const writeJson = async (root: string, relativePath: string, value: unknown) => {
  const filePath = path.join(root, relativePath);
  await mkdir(path.dirname(filePath), { recursive: true });
  await writeFile(filePath, JSON.stringify(value));
};

const makeContentRoot = async () => {
  const root = await mkdtemp(path.join(os.tmpdir(), 'cms-builder-content-'));
  roots.push(root);
  await writeJson(root, 'site.json', {
    name: 'Test site', domain: 'https://example.test', blogPath: '/articles', template: 'anonime', theme: 'emerald',
    seo: { title: 'Test site', description: 'A test site.' },
  });
  await writeJson(root, 'home.json', { seo: { title: 'Home', description: 'Home.' }, sections: [{ type: 'latest_articles', enabled: true, heading: 'Latest' }] });
  await writeJson(root, 'authors/editor.json', { slug: 'editor', name: 'Editor' });
  await writeJson(root, 'categories/guides.json', { slug: 'guides', name: 'Guides' });
  await writeJson(root, 'tags/cms.json', { slug: 'cms', name: 'CMS' });
  return root;
};

afterEach(async () => {
  delete process.env.CMS_CONTENT_ROOT;
  await Promise.all(roots.splice(0).map((root) => rm(root, { recursive: true, force: true })));
});

describe('loadBuildData', () => {
  it('resolves references, filters drafts, and sorts published articles newest first', async () => {
    const root = await makeContentRoot();
    const article = (slug: string, date: string, status: string) => ({
      title: slug, slug, excerpt: 'Excerpt', body: 'Body', author: 'editor', category: 'guides', tags: ['cms'],
      publishedAt: date, status, featured: false, seo: { title: slug, description: 'Description' },
    });
    await writeJson(root, 'articles/older.json', article('older', '2026-01-01', 'published'));
    await writeJson(root, 'articles/newer.json', article('newer', '2026-02-01', 'published'));
    await writeJson(root, 'articles/draft.json', article('draft', '2026-03-01', 'draft'));
    process.env.CMS_CONTENT_ROOT = root;

    const result = loadBuildData();

    expect(result.articles.map((item) => item.slug)).toEqual(['newer', 'older']);
    expect(result.articles[0].authorName).toBe('Editor');
    expect(result.articles[0].categoryName).toBe('Guides');
  });

  it('rejects unsafe content URLs', async () => {
    const root = await makeContentRoot();
    await writeJson(root, 'articles/unsafe.json', {
      title: 'Unsafe', slug: 'unsafe', excerpt: 'Excerpt', body: 'Body', author: 'editor', category: 'guides', tags: [],
      publishedAt: '2026-01-01', status: 'published', featured: false, coverImage: 'javascript:alert(1)',
      seo: { title: 'Unsafe', description: 'Description' },
    });
    process.env.CMS_CONTENT_ROOT = root;

    expect(() => loadBuildData()).toThrow(/unsafe|root-relative/i);
  });
});
