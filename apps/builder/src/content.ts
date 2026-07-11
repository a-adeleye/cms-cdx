import { readFileSync } from 'node:fs';

export type BuildArticle = {
  title: string;
  slug: string;
  excerpt: string;
  contentMarkdown: string;
  seoTitle: string;
  seoDescription: string;
  canonicalUrl: string;
};

export type BuildSection = {
  sectionKey: string;
  title: string;
  subtitle: string;
  contentJson: string;
  isEnabled: boolean;
};

export type BuildSite = {
  name: string;
  domain: string;
  blogPath: string;
};

type BuildData = {
  site: BuildSite;
  articles: BuildArticle[];
  sections: BuildSection[];
};

const fallbackData: BuildData = {
  site: { name: 'Example Site', domain: 'http://localhost:4321', blogPath: '/articles' },
  articles: [],
  sections: [
    {
      sectionKey: 'hero',
      title: 'Build once, deploy many',
      subtitle: 'CMS-backed static sites',
      contentJson: '{}',
      isEnabled: true,
    },
  ],
};

export function loadBuildData(): BuildData {
  const filePath = process.env.CMS_BUILD_DATA_FILE;
  if (!filePath) {
    return fallbackData;
  }

  return JSON.parse(readFileSync(filePath, 'utf8')) as BuildData;
}
