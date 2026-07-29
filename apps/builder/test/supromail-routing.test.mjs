import assert from 'node:assert/strict';
import { mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join } from 'node:path';
import { spawnSync } from 'node:child_process';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const builderRoot = dirname(dirname(fileURLToPath(import.meta.url)));
const repositoryRoot = dirname(dirname(builderRoot));
const astroCli = join(repositoryRoot, 'node_modules', 'astro', 'astro.js');
const templateRoot = join(repositoryRoot, 'packages', 'templates', 'supromail');

test('Supromail renders CMS content through the production Astro routes', () => {
  const temporaryRoot = mkdtempSync(join(tmpdir(), 'cms-builder-supromail-'));
  const dataPath = join(temporaryRoot, 'build-data.json');
  const outputPath = join(temporaryRoot, 'output');
  const articles = Array.from({ length: 7 }, (_, index) => ({
    title: `CMS article ${index + 1}`,
    slug: `cms-article-${index + 1}`,
    excerpt: `Excerpt for CMS article ${index + 1}.`,
    contentMarkdown: `# CMS article ${index + 1}\n\n## Detail\n\nSafe **CMS** content.`,
    seoTitle: `CMS article ${index + 1}`,
    seoDescription: `SEO description ${index + 1}`,
    canonicalUrl: '',
    categoryName: index % 2 === 0 ? 'Engineering' : 'Product',
    authorName: 'CMS Editor',
    publishedAt: `2026-07-${String(20 - index).padStart(2, '0')}`,
    coverImageUrl: 'https://cdn.example/cover.png',
    isFeatured: index === 0,
  }));
  writeFileSync(dataPath, JSON.stringify({
    site: { name: 'CMS Supromail', domain: 'https://cms.example', blogPath: '/blog' },
    articles,
    sections: [],
  }));

  try {
    const build = spawnSync(process.execPath, [astroCli, 'build'], {
      cwd: builderRoot,
      encoding: 'utf8',
      env: {
        ...process.env,
        CMS_BUILD_DATA_FILE: dataPath,
        CMS_BUILD_OUTPUT_DIR: outputPath,
        CMS_SITE_URL: 'https://cms.example',
        CMS_TEMPLATE_KEY: 'supromail',
        CMS_TEMPLATE_ROOT: templateRoot,
      },
    });
    assert.equal(build.status, 0, `${build.stdout}\n${build.stderr}`);

    const landing = readFileSync(join(outputPath, 'index.html'), 'utf8');
    const blogLanding = readFileSync(join(outputPath, 'articles', 'index.html'), 'utf8');
    const articlesPage = readFileSync(join(outputPath, 'articles', 'all', 'index.html'), 'utf8');
    const articlePage = readFileSync(join(outputPath, 'articles', 'cms-article-1', 'index.html'), 'utf8');
    const paginationPage = readFileSync(join(outputPath, 'articles', 'page', '2', 'index.html'), 'utf8');
    const sitemap = readFileSync(join(outputPath, 'sitemap.xml'), 'utf8');
    const blogSitemap = readFileSync(join(outputPath, 'articles', 'sitemap.xml'), 'utf8');

    assert.match(landing, /^<!DOCTYPE html><html/);
    assert.match(landing, /The CMS Supromail blog/);
    assert.match(landing, /CMS article 1/);
    assert.match(landing, /<style>[\s\S]*\.supromail-site/);
    assert.match(blogLanding, /Notes on messaging <em>you actually own\.<\/em>/);
    assert.match(articlesPage, /All articles/);
    assert.match(articlePage, /Safe <strong>CMS<\/strong> content/);
    assert.match(articlePage, /On this page/);
    assert.match(paginationPage, /CMS article 7/);
    assert.match(landing, /href="blog\/cms-article-1\/"/);
    assert.doesNotMatch(landing, /href="\/(?:_astro|blog)\//);
    assert.match(sitemap, /<loc>https:\/\/cms\.example\/<\/loc>/);
    assert.match(sitemap, /<loc>https:\/\/cms\.example\/blog\/<\/loc>/);
    assert.match(sitemap, /<loc>https:\/\/cms\.example\/blog\/all\/<\/loc>/);
    assert.match(sitemap, /<loc>https:\/\/cms\.example\/blog\/page\/2\/<\/loc>/);
    assert.match(sitemap, /<loc>https:\/\/cms\.example\/blog\/cms-article-1\/<\/loc><lastmod>2026-07-20<\/lastmod>/);
    assert.doesNotMatch(sitemap, /https:\/\/cms\.example\/articles\//);
    assert.match(blogSitemap, /<loc>https:\/\/cms\.example\/blog\/cms-article-1\/<\/loc>/);
    assert.doesNotMatch(blogSitemap, /<loc>https:\/\/cms\.example\/<\/loc>/);
  } finally {
    rmSync(temporaryRoot, { recursive: true, force: true });
  }
});
