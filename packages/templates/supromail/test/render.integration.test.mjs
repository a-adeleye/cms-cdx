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
  const temporaryRoot = mkdtempSync(join(tmpdir(), 'supromail-template-'));
  const outputDirectory = join(temporaryRoot, 'dist');
  const isolatedTemplate = join(temporaryRoot, 'supromail');
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
          SUPROMAIL_TEST_OUTPUT: outputDirectory,
          SUPROMAIL_TEST_CACHE: join(temporaryRoot, 'astro-cache'),
          SUPROMAIL_TEST_VITE_CACHE: join(temporaryRoot, 'vite-cache'),
        },
      },
    );
    assert.equal(build.status, 0, `${build.stdout}\n${build.stderr}`);

    const html = readFileSync(join(outputDirectory, 'index.html'), 'utf8');
    const landingHtml = readFileSync(join(outputDirectory, 'landing', 'index.html'), 'utf8');
    const articlesHtml = readFileSync(join(outputDirectory, 'articles', 'index.html'), 'utf8');
    const articleHtml = readFileSync(join(outputDirectory, 'article', 'index.html'), 'utf8');
    const singleArticleLandingHtml = readFileSync(join(outputDirectory, 'single-article', 'index.html'), 'utf8');

    assert.match(landingHtml, /Notes on messaging <em>you actually own\.<\/em>/);
    assert.match(landingHtml, /Latest writing/);
    assert.match(articlesHtml, /All articles/);
    assert.match(articleHtml, /On this page/);
    assert.match(singleArticleLandingHtml, /The only published article/);
    assert.doesNotMatch(singleArticleLandingHtml, /Latest writing/);
    assert.doesNotMatch(singleArticleLandingHtml, /No articles published yet\./);

    for (const renderedPage of [landingHtml, articlesHtml, articleHtml, singleArticleLandingHtml]) {
      assert.match(renderedPage, /Start sending from what you own/);
      assert.match(renderedPage, /Start sending for free/);
      assert.doesNotMatch(renderedPage, /Get the next one by email/);
    }

    for (const renderedPage of [landingHtml, articlesHtml, articleHtml]) {
      assert.match(renderedPage, /href="https:\/\/supromail\.com\/#platform"/);
      assert.match(renderedPage, /href="https:\/\/supromail\.com\/pricing"/);
      assert.match(renderedPage, /href="https:\/\/app\.supromail\.com"/);
      assert.match(renderedPage, /class="container footer-wordmark">Supromail</);
      assert.match(renderedPage, /One workspace for every conversation your business owns\./);
      // Internal blog routes stay relative so a static artifact also works
      // when the CMS serves it below its local /deployments mount.
      assert.doesNotMatch(renderedPage, /href="\/blog\//);
    }

    // The listing page keeps the client-side category filter wired to real data.
    assert.match(articlesHtml, /data-filter="deliverability"/);
    assert.match(articlesHtml, /data-category="engineering"/);
    assert.match(articlesHtml, /id="emptyState"/);

    // The read page renders the design fixture end to end.
    assert.match(articleHtml, /<h1>Why owned sending beats rented reputation<\/h1>/);
    assert.match(articleHtml, /<h2 id="what-owned-sending-changes">/);
    assert.match(articleHtml, /href="#what-owned-sending-changes"/);
    assert.match(articleHtml, /href="\.\.\/blog\/sms-failover-across-devices\/"/);
    assert.match(articleHtml, /id="readProgress"/);
    assert.match(articleHtml, /<code>mail\.yourcompany\.com<\/code>/);
    assert.doesNotMatch(landingHtml, /id="readProgress"/);

    // Hostile Markdown stays inert on the security fixture.
    assert.doesNotMatch(html, /<script>alert\(1\)<\/script>/);
    assert.doesNotMatch(html, /<script>alert\(2\)<\/script>/);
    assert.doesNotMatch(html, /javascript:alert/);
    assert.doesNotMatch(html, /Why owned sending beats rented reputation/);
    assert.match(html, /&lt;script&gt;alert\(1\)&lt;\/script&gt;/);
    assert.match(html, /&lt;script&gt;alert\(2\)&lt;\/script&gt;/);
    assert.match(html, /<strong>Bold<\/strong>/);
    assert.match(html, /href="https:\/\/example\.com" rel="noopener noreferrer"/);
    assert.match(html, /<h2 id="safe-heading">Safe heading<\/h2>/);
  } finally {
    rmSync(temporaryRoot, { recursive: true, force: true });
  }
});
