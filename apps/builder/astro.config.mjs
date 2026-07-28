import { fileURLToPath } from 'node:url';
import { defineConfig } from 'astro/config';

const defaultTemplateRoot = fileURLToPath(new URL('../../packages/templates/supromail', import.meta.url));
const templateRoot = process.env.CMS_TEMPLATE_ROOT ?? defaultTemplateRoot;
const isSupromailBuild = process.env.CMS_TEMPLATE_KEY === 'supromail';

export default defineConfig({
  site: process.env.CMS_SITE_URL ?? 'http://localhost:4321',
  outDir: process.env.CMS_BUILD_OUTPUT_DIR ?? './dist',
  build: {
    inlineStylesheets: isSupromailBuild ? 'always' : 'auto',
  },
  vite: {
    resolve: {
      alias: {
        '@cms-template': templateRoot,
      },
    },
  },
});

