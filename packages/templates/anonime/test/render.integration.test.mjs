import assert from 'node:assert/strict';
import { cpSync, mkdtempSync, readFileSync, rmSync } from 'node:fs';
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
    assert.equal(readFileSync(join(outputDirectory, 'landing', 'index.html'), 'utf8').includes('Thoughts that'), true);
    assert.equal(readFileSync(join(outputDirectory, 'articles', 'index.html'), 'utf8').includes('All Articles'), true);
    assert.equal(readFileSync(join(outputDirectory, 'article', 'index.html'), 'utf8').includes('On this page'), true);
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
