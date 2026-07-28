import { defineConfig } from 'astro/config';

const testOutput = process.env.SUPROMAIL_TEST_OUTPUT;
const testCache = process.env.SUPROMAIL_TEST_CACHE;
const testViteCache = process.env.SUPROMAIL_TEST_VITE_CACHE;

export default defineConfig({
  outDir: testOutput ?? './dist',
  cacheDir: testCache ?? './node_modules/.astro',
  vite: { cacheDir: testViteCache ?? './node_modules/.vite' },
});
