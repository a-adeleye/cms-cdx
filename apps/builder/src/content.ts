import { readFileSync, readdirSync } from 'node:fs';
import path from 'node:path';

export type BuildArticle = {
  title: string;
  slug: string;
  excerpt: string;
  contentMarkdown: string;
  seoTitle: string;
  seoDescription: string;
  canonicalUrl: string;
  authorName: string;
  categoryName: string;
  tagNames: string[];
  publishedAt: string;
  isFeatured: boolean;
  imageUrl?: string;
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
  template: 'anonime' | 'generic';
  theme: 'emerald' | 'ocean';
  seoTitle: string;
  seoDescription: string;
};

export type BuildData = { site: BuildSite; articles: BuildArticle[]; sections: BuildSection[] };

const slugPattern = /^[a-z0-9]+(?:-[a-z0-9]+)*$/;
const allowedSectionTypes = new Set(['hero', 'featured_articles', 'topic_tabs', 'latest_articles', 'mission_cta']);

const fail = (message: string): never => { throw new Error(`Invalid CMS content: ${message}`); };
const asObject = (value: unknown, label: string): Record<string, any> => {
  if (!value || typeof value !== 'object' || Array.isArray(value)) fail(`${label} must be an object.`);
  return value as Record<string, any>;
};
const asString = (value: unknown, label: string, max = 10_000): string => {
  if (typeof value !== 'string' || !value.trim() || value.length > max) fail(`${label} must be a non-empty string.`);
  return value;
};
const asSlug = (value: unknown, label: string): string => {
  const slug = asString(value, label, 120);
  if (!slugPattern.test(slug)) fail(`${label} must be a safe slug.`);
  return slug;
};
const asUrl = (value: unknown, label: string, allowEmpty = true): string => {
  if (allowEmpty && (value === '' || value == null)) return '';
  const url = asString(value, label, 2_000);
  if (url.startsWith('/') && !url.startsWith('//')) return url;
  try {
    const parsed = new URL(url);
    if (!['http:', 'https:'].includes(parsed.protocol) || parsed.username || parsed.password) fail(`${label} is unsafe.`);
  } catch { fail(`${label} must be root-relative or HTTP(S).`); }
  return url;
};

const readJson = (filePath: string): unknown => {
  try { return JSON.parse(readFileSync(filePath, 'utf8')); }
  catch { return fail(`${path.basename(filePath)} is not valid JSON.`); }
};

const readCollection = (root: string, name: string): Record<string, any>[] => {
  const directory = path.join(root, name);
  return readdirSync(directory, { withFileTypes: true })
    .filter((entry) => entry.isFile() && entry.name.endsWith('.json'))
    .sort((a, b) => a.name.localeCompare(b.name))
    .map((entry) => asObject(readJson(path.join(directory, entry.name)), `${name}/${entry.name}`));
};

export function loadBuildData(): BuildData {
  const legacyFile = process.env.CMS_BUILD_DATA_FILE;
  if (legacyFile) return JSON.parse(readFileSync(legacyFile, 'utf8')) as BuildData;

  const defaultRoot = path.resolve(process.cwd(), 'content');
  const root = process.env.CMS_CONTENT_ROOT ? path.resolve(process.env.CMS_CONTENT_ROOT) : defaultRoot;
  const siteInput = asObject(readJson(path.join(root, 'site.json')), 'site.json');
  const homeInput = asObject(readJson(path.join(root, 'home.json')), 'home.json');
  const siteSeo = asObject(siteInput.seo, 'site.seo');
  const template = asString(siteInput.template, 'site.template') as BuildSite['template'];
  const theme = asString(siteInput.theme, 'site.theme') as BuildSite['theme'];
  if (!['anonime', 'generic'].includes(template)) fail('site.template is unsupported.');
  if (!['emerald', 'ocean'].includes(theme)) fail('site.theme is unsupported.');
  const blogPath = asUrl(siteInput.blogPath, 'site.blogPath', false);

  const authors = new Map(readCollection(root, 'authors').map((item) => [asSlug(item.slug, 'author.slug'), asString(item.name, 'author.name', 200)]));
  const categories = new Map(readCollection(root, 'categories').map((item) => [asSlug(item.slug, 'category.slug'), asString(item.name, 'category.name', 200)]));
  const tags = new Map(readCollection(root, 'tags').map((item) => [asSlug(item.slug, 'tag.slug'), asString(item.name, 'tag.name', 200)]));

  const seenSlugs = new Set<string>();
  const articles = readCollection(root, 'articles').map((item): BuildArticle & { status: string } => {
    const slug = asSlug(item.slug, 'article.slug');
    if (seenSlugs.has(slug)) fail(`article slug "${slug}" is duplicated.`);
    seenSlugs.add(slug);
    const author = asSlug(item.author, `article ${slug} author`);
    const category = asSlug(item.category, `article ${slug} category`);
    if (!authors.has(author) || !categories.has(category)) fail(`article "${slug}" has a missing author or category.`);
    const tagSlugs = Array.isArray(item.tags) ? item.tags.map((tag) => asSlug(tag, `article ${slug} tag`)) : [];
    if (tagSlugs.some((tag) => !tags.has(tag))) fail(`article "${slug}" has a missing tag.`);
    const status = asString(item.status, `article ${slug} status`);
    if (!['draft', 'published'].includes(status)) fail(`article "${slug}" has an invalid status.`);
    const publishedAt = asString(item.publishedAt, `article ${slug} publishedAt`, 10);
    if (!/^\d{4}-\d{2}-\d{2}$/.test(publishedAt) || Number.isNaN(Date.parse(`${publishedAt}T00:00:00Z`))) fail(`article "${slug}" has an invalid date.`);
    const seo = asObject(item.seo, `article ${slug} seo`);
    return {
      status,
      slug,
      title: asString(item.title, `article ${slug} title`, 200),
      excerpt: asString(item.excerpt, `article ${slug} excerpt`, 1_000),
      contentMarkdown: asString(item.body, `article ${slug} body`, 100_000),
      seoTitle: asString(seo.title, `article ${slug} seo.title`, 200),
      seoDescription: asString(seo.description, `article ${slug} seo.description`, 1_000),
      canonicalUrl: asUrl(seo.canonicalUrl, `article ${slug} canonicalUrl`),
      authorName: authors.get(author)!, categoryName: categories.get(category)!,
      tagNames: tagSlugs.map((tag) => tags.get(tag)!), publishedAt,
      isFeatured: item.featured === true,
      imageUrl: asUrl(item.coverImage, `article ${slug} coverImage`) || undefined,
    };
  }).filter((article) => article.status === 'published')
    .sort((a, b) => b.publishedAt.localeCompare(a.publishedAt) || a.slug.localeCompare(b.slug));

  if (!Array.isArray(homeInput.sections)) fail('home.sections must be an array.');
  const sections = homeInput.sections.map((raw, index): BuildSection => {
    const section = asObject(raw, `home.sections[${index}]`);
    const type = asString(section.type, `home.sections[${index}].type`, 80);
    if (!allowedSectionTypes.has(type)) fail(`home.sections[${index}] has an unsupported type.`);
    return {
      sectionKey: type,
      title: typeof section.heading === 'string' ? section.heading : (typeof section.eyebrow === 'string' ? section.eyebrow : type),
      subtitle: typeof section.body === 'string' ? section.body : '',
      contentJson: JSON.stringify(section),
      isEnabled: section.enabled !== false,
    };
  });

  return {
    site: {
      name: asString(siteInput.name, 'site.name', 200), domain: asUrl(siteInput.domain, 'site.domain', false),
      blogPath, template, theme,
      seoTitle: asString(siteSeo.title, 'site.seo.title', 200),
      seoDescription: asString(siteSeo.description, 'site.seo.description', 1_000),
    },
    articles,
    sections,
  };
}
