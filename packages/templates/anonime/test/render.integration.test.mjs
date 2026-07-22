import assert from 'node:assert/strict';
import { cpSync, mkdtempSync, readFileSync, rmSync, symlinkSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { basename, dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';
import test from 'node:test';

const testDirectory = dirname(fileURLToPath(import.meta.url));
const templateRoot = dirname(testDirectory);
const repositoryRoot = fileURLToPath(new URL('../../../../', import.meta.url));
const astroCli = join(repositoryRoot, 'node_modules', 'astro', 'astro.js');

test('Astro renders formatting while keeping hostile Markdown inert', () => {
  const temporaryRoot = mkdtempSync(join(tmpdir(), 'anonime-template-'));
  const outputDirectory = join(temporaryRoot, 'dist');
  const isolatedTemplate = join(temporaryRoot, 'anonime');
  cpSync(templateRoot, isolatedTemplate, {
    recursive: true,
    filter: (source) => !['.astro', 'dist', 'node_modules'].includes(basename(source)),
  });
  symlinkSync(join(repositoryRoot, 'node_modules'), join(isolatedTemplate, 'node_modules'), 'junction');
  const isolatedFixture = join(isolatedTemplate, 'test', 'fixture');

  try {
    const build = spawnSync(
      process.execPath,
      [astroCli, 'build'],
      {
        cwd: isolatedFixture,
        encoding: 'utf8',
        env: {
          ...process.env,
          ANONIME_TEST_OUTPUT: outputDirectory,
          ANONIME_TEST_CACHE: join(temporaryRoot, 'astro-cache'),
          ANONIME_TEST_VITE_CACHE: join(temporaryRoot, 'vite-cache'),
        },
      },
    );
    assert.equal(build.status, 0, `${build.stdout}\n${build.stderr}`);

    const html = readFileSync(join(outputDirectory, 'index.html'), 'utf8');
    const landingHtml = readFileSync(join(outputDirectory, 'landing', 'index.html'), 'utf8');
    const articlesHtml = readFileSync(join(outputDirectory, 'articles', 'index.html'), 'utf8');
    const articleHtml = readFileSync(join(outputDirectory, 'article', 'index.html'), 'utf8');
    assert.equal(landingHtml.includes('Thoughts that'), true);
    assert.equal(articlesHtml.includes('All Articles'), true);
    assert.equal(articleHtml.includes('On this page'), true);
    for (const renderedHero of [landingHtml, articlesHtml]) {
      assert.match(renderedHero, /class="anonime-hero-art"/);
      assert.match(renderedHero, /anonime-blog-hero-white\.webp/);
      assert.match(renderedHero, /anonime-blog-hero-black\.webp/);
    }
    for (const renderedPage of [landingHtml, articlesHtml, articleHtml]) {
      assert.match(renderedPage, /href="https:\/\/anonime\.io\/#top"/);
      assert.match(renderedPage, /href="https:\/\/anonime\.io\/#how-it-works"/);
      assert.match(renderedPage, /href="https:\/\/anonime\.io\/pricing"/);
      assert.match(renderedPage, /href="https:\/\/anonime\.io\/blog"/);
      assert.match(renderedPage, /data-anonime-theme-control/);
      assert.match(renderedPage, /Built with <strong>privacy by design\.<\/strong>/);
      assert.match(renderedPage, /Hosted with privacy-first infra/);
      assert.match(renderedPage, /© 2026 Anonime, Inc\. All rights reserved\./);
    }
    assert.match(landingHtml, /localStorage\.setItem\(storageKey, theme\)/);
    assert.doesNotMatch(html, /<script>alert\(1\)<\/script>/);
    assert.doesNotMatch(html, /javascript:alert/);
    assert.doesNotMatch(html, /Beyond Encryption: Zero-Knowledge Architecture/);
    assert.match(html, /&lt;script&gt;alert\(1\)&lt;\/script&gt;/);
    assert.match(html, /<strong>Bold<\/strong>/);
    assert.match(html, /href="https:\/\/example\.com" rel="noopener noreferrer"/);
    assert.match(html, /<h2 id="safe-heading">Safe heading<\/h2>/);
  } finally {
    rmSync(temporaryRoot, { recursive: true, force: true });
  }
});
