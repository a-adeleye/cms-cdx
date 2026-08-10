import { defineConfig } from 'astro/config';

export default defineConfig({
  devToolbar: { enabled: false },
  site: process.env.CMS_SITE_URL ?? 'http://localhost:4321',
  outDir: process.env.CMS_BUILD_OUTPUT_DIR ?? './dist',
});

