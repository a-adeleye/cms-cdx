export type WorkspacePageKind =
  | 'login'
  | 'dashboard'
  | 'settings'
  | 'landing-page-editor'
  | 'articles'
  | 'media-library'
  | 'ai-assistant'
  | 'builds'
  | 'publishing';

export interface WorkspacePageConfig {
  path: string;
  navLabel: string;
  kind: WorkspacePageKind;
  eyebrow: string;
  title: string;
  primaryAction: {
    label: string;
    path: string;
  };
  secondaryAction?: {
    label: string;
    path: string;
  };
}

export type ArticleStatus = 'draft' | 'review' | 'published' | 'archived';

export type BuildStatus = 'queued' | 'running' | 'success' | 'failed';

export type BuildType = 'preview' | 'published';

export interface AuthSession {
  email: string;
  fullName: string;
  role: 'admin' | 'editor';
}

export interface SiteRecord {
  id: string;
  name: string;
  slug: string;
  domain: string;
  blogPath: string;
  status: 'active' | 'inactive';
  templateKey: string;
  themeConfig: string;
  deployProvider: string;
  deployConfig: string;
  previewDeployProvider: string;
  previewDeployConfig: string;
  aiConfig: string;
  storageConfig: string;
  deploymentWarnings?: string[];
  updatedAt: string;
}

export interface TemplateRecord {
  id: string;
  name: string;
  slug: string;
  updatedAt: string;
}

export interface LandingSectionRecord {
  id: string;
  siteId: string;
  sectionKey: string;
  title: string;
  subtitle: string;
  contentJson: string;
  displayOrder: number;
  isEnabled: boolean;
}

export interface ArticleRecord {
  id: string;
  siteId: string;
  authorId: string;
  categoryId: string;
  title: string;
  slug: string;
  excerpt: string;
  contentMarkdown: string;
  coverImageUrl: string;
  status: ArticleStatus;
  isFeatured: boolean;
  publishedAt: string | null;
  seoTitle: string;
  seoDescription: string;
  canonicalUrl: string;
  generatedByAi: boolean;
  humanReviewed: boolean;
  aiPrompt: string;
  aiModel: string;
  tagIds: string[];
  updatedAt: string;
}

export interface AuthorRecord {
  id: string;
  siteId: string;
  name: string;
  slug: string;
  bio: string;
}

export interface CategoryRecord {
  id: string;
  siteId: string;
  name: string;
  slug: string;
  description: string;
}

export interface TagRecord {
  id: string;
  siteId: string;
  name: string;
  slug: string;
}

export interface MediaAssetRecord {
  id: string;
  siteId: string;
  fileName: string;
  fileUrl: string;
  mimeType: string;
  sizeBytes: number;
  storageProvider: string;
  storageKey: string;
  altText: string;
}

export interface BuildRecord {
  id: string;
  siteId: string;
  status: BuildStatus;
  buildType: BuildType;
  logs: string;
  outputPath: string;
  deployProvider: string;
  deployStatus: string;
  deployUrl: string;
  startedAt: string | null;
  finishedAt: string | null;
}

export interface AdminStateSnapshot {
  authSession: AuthSession | null;
  selectedSiteId: string;
  selectedArticleId: string | null;
  sites: SiteRecord[];
  templates: TemplateRecord[];
  landingSections: LandingSectionRecord[];
  articles: ArticleRecord[];
  authors: AuthorRecord[];
  categories: CategoryRecord[];
  tags: TagRecord[];
  mediaAssets: MediaAssetRecord[];
  builds: BuildRecord[];
}
