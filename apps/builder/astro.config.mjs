import { defineConfig } from 'astro/config';

export default defineConfig({
  site: process.env.CMS_SITE_URL ?? 'http://localhost:4321',
  outDir: process.env.CMS_BUILD_OUTPUT_DIR ?? './dist',
});

