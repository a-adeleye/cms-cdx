export type TemplateSite = {
  name?: string;
  domain?: string;
  blogPath?: string;
};

export type TemplateArticle = {
  title: string;
  slug: string;
  excerpt?: string;
  contentMarkdown?: string;
  seoTitle?: string;
  seoDescription?: string;
  canonicalUrl?: string;
  category?: string | { name?: string };
  author?: string | { name?: string };
  publishedAt?: string;
  published_at?: string;
  readingTime?: number | string;
  reading_time?: number | string;
  imageUrl?: string;
  featuredImage?: string;
  coverImageUrl?: string;
  authorName?: string;
  categoryName?: string;
  isFeatured?: boolean;
};

export type ArticleViewModel = {
  title: string;
  slug: string;
  excerpt: string;
  contentMarkdown: string;
  seoTitle: string;
  seoDescription: string;
  canonicalUrl: string;
  category: string;
  author: string;
  publishedLabel: string;
  readingTime: number;
  imageUrl?: string;
  visual: VisualVariant;
  isFeatured: boolean;
};

export type VisualVariant =
  | 'silhouette'
  | 'key'
  | 'message'
  | 'identity'
  | 'product'
  | 'fingerprint'
  | 'architecture'
  | 'email'
  | 'brokers'
  | 'minimalism'
  | 'lock';

export type SiteViewModel = {
  name: string;
  blogPath: string;
};

const articleContent = `In today's digital world, nearly every conversation leaves a trace. Emails are scanned, messages are logged, and data is stored far longer than most people realize. Encrypted communication changes that.

It ensures that only the intended recipient can read the message — not your provider, not hackers, not even governments. But who really needs it? The answer may surprise you.

## 1. Individuals Protecting Their Personal Lives

From dating apps to medical consultations, we share deeply personal information online. Without encryption, this data can be exposed in breaches or sold to third parties.

> Privacy is not about hiding something. It's about protecting the right to be left alone.

Encrypted messaging ensures your conversations stay between you and the people you trust.

## 2. Businesses Handling Sensitive Information

Companies deal with confidential data every day — finances, client details, product plans, and more. A single leak can damage trust and cost millions.

> Encrypted communication helps teams share, collaborate, and operate securely.

Whether you're a startup or an enterprise, encryption is now a business essential.

## 3. Journalists, Activists, and Communities at Risk

For those who speak truth to power, privacy is safety. Encrypted communication protects sources, organizes movements, and ensures critical information reaches the right people.

> In high-risk environments, encryption can be the difference between exposure and protection.

## The Bottom Line

Encrypted communication isn't just for tech experts or security professionals. It's for anyone who values privacy, security, and control over their words and data.`;

export const featuredDesignArticle: ArticleViewModel = {
  title: 'Who Needs Encrypted Communication? Scenarios Where Privacy Matters Most',
  slug: 'who-needs-encrypted-communication',
  excerpt:
    'Explore real cases where encrypted communication is essential — from protecting sensitive business information to safeguarding personal privacy.',
  contentMarkdown: articleContent,
  seoTitle: 'Who Needs Encrypted Communication?',
  seoDescription: 'Explore the people and situations where private communication matters most.',
  canonicalUrl: '',
  category: 'Secure Communication',
  author: 'Rene Carter',
  publishedLabel: 'Feb 7, 2026',
  readingTime: 5,
  visual: 'lock',
  isFeatured: true,
};

export const designArticles: ArticleViewModel[] = [
  designArticle(
    'The Problem with “Free”: Why Ad-Driven Services Put Your Privacy at Risk',
    'the-problem-with-free',
    'Uncover the hidden costs of free online services that rely on advertising and data collection.',
    'Privacy Basics',
    'Feb 2, 2026',
    4,
    'silhouette',
  ),
  designArticle(
    'Leaving No Trace: How to Communicate Anonymously and Securely with Anonime',
    'leaving-no-trace',
    'Discover how Anonime enables private, anonymous communication through end-to-end encryption.',
    'Guides',
    'Jan 29, 2026',
    5,
    'key',
  ),
  designArticle(
    'Reclaim Your Digital Footprint: The Power of Ephemeral Messaging',
    'reclaim-your-digital-footprint',
    'Tired of your data living online forever? Discover how ephemeral messaging helps you stay in control.',
    'Security',
    'Jan 27, 2026',
    4,
    'message',
  ),
  designArticle(
    'Protect Your Online Identity: The Power of Private Personas',
    'protect-your-online-identity',
    'Discover how creating distinct online personas shields your real identity from tracking and profiling.',
    'Identity',
    'Jan 23, 2026',
    10,
    'identity',
  ),
  designArticle(
    'What’s New in Anonime: Updates That Put You First',
    'whats-new-in-anonime',
    'Explore the latest features and improvements built for your privacy.',
    'Product',
    'Jan 21, 2026',
    4,
    'product',
  ),
  designArticle(
    'Welcome to the Anonime Blog',
    'welcome-to-anonime',
    'Introducing the Anonime blog — your source for privacy insights, security tips, and tools to take back control of your digital life.',
    'Welcome',
    'Jan 20, 2026',
    2,
    'fingerprint',
  ),
  designArticle(
    'Beyond Encryption: Zero-Knowledge Architecture and Why It Matters for Your Privacy',
    'beyond-encryption',
    'A practical guide to systems that never need to see your private information.',
    'Privacy Basics',
    'Jan 24, 2026',
    5,
    'architecture',
  ),
  designArticle(
    'Why Your Email Address Is Your Biggest Privacy Liability',
    'email-address-privacy-liability',
    'Why a familiar identifier can quietly connect your activity across the web.',
    'Privacy Basics',
    'Jan 22, 2026',
    3,
    'email',
  ),
  designArticle(
    'The Hidden Risks of Data Brokers',
    'hidden-risks-of-data-brokers',
    'Understand how personal profiles are assembled, traded, and used.',
    'Security',
    'Jan 18, 2026',
    6,
    'brokers',
  ),
  designArticle(
    'Digital Minimalism: Owning Less of Yourself Online',
    'digital-minimalism',
    'A calmer approach to reducing the personal information you leave behind.',
    'Guides',
    'Jan 16, 2026',
    7,
    'minimalism',
  ),
];

function designArticle(
  title: string,
  slug: string,
  excerpt: string,
  category: string,
  publishedLabel: string,
  readingTime: number,
  visual: VisualVariant,
): ArticleViewModel {
  return {
    title,
    slug,
    excerpt,
    contentMarkdown: '',
    seoTitle: title,
    seoDescription: excerpt,
    canonicalUrl: '',
    category,
    author: 'Rene Carter',
    publishedLabel,
    readingTime,
    visual,
    isFeatured: false,
  };
}

export class TemplateContentError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'TemplateContentError';
  }
}

export function resolveSite(value?: unknown): SiteViewModel {
  if (value === undefined) return { name: 'anonime', blogPath: '/articles' };
  const site = requireRecord(value, 'Site configuration');
  const configuredName = readOptionalString(site, 'name');
  const blogPath = readOptionalString(site, 'blogPath');
  return {
    name: configuredName && configuredName !== 'Example Site' ? configuredName : 'anonime',
    blogPath: normalizeBlogPath(blogPath),
  };
}

export function normalizeBlogPath(value?: string): string {
  const path = value?.trim() || '/articles';
  if (path === '/' || !/^\/[a-z0-9/_-]*$/i.test(path) || path.startsWith('//')) {
    return '/articles';
  }
  return path.replace(/\/+$/, '');
}

export function normalizePositiveInteger(value: unknown, fallback: number): number {
  if (value === undefined) return fallback;
  if (typeof value !== 'number' || !Number.isSafeInteger(value) || value < 1) {
    throw new TemplateContentError('Pagination values must be positive integers.');
  }
  return value;
}

export function resolveArticles(value?: unknown): ArticleViewModel[] {
  if (value === undefined) return designArticles;
  if (!Array.isArray(value)) {
    throw new TemplateContentError('Articles must be provided as an array.');
  }
  return value.map((article, index) => adaptArticle(article, index));
}

export function resolveFeaturedArticle(value?: unknown): ArticleViewModel {
  return value === undefined ? featuredDesignArticle : adaptArticle(value, 10, true);
}

export function adaptArticle(value: unknown, index: number, requiresContent = false): ArticleViewModel {
  const article = requireRecord(value, 'Article');
  const title = readRequiredString(article, 'title', 'Article title');
  const slug = normalizeSlug(readRequiredString(article, 'slug', 'Article slug'));
  const contentMarkdown = readOptionalString(article, 'contentMarkdown');
  if (requiresContent && !contentMarkdown) {
    throw new TemplateContentError('Article contentMarkdown is required for the read page.');
  }
  const excerpt = readOptionalString(article, 'excerpt') || deriveExcerpt(contentMarkdown);
  const visual = visualVariants[index % visualVariants.length];
  return {
    title,
    slug,
    excerpt,
    contentMarkdown,
    seoTitle: readOptionalString(article, 'seoTitle') || title,
    seoDescription: readOptionalString(article, 'seoDescription') || excerpt,
    canonicalUrl: safeMediaUrl(article.canonicalUrl) || '',
    category: readOptionalString(article, 'categoryName') || readName(article.category) || 'Privacy',
    author: readOptionalString(article, 'authorName') || readName(article.author) || 'Editorial team',
    publishedLabel: formatPublishedDate(article.publishedAt ?? article.published_at),
    readingTime: readTime(article, contentMarkdown),
    imageUrl: safeMediaUrl(article.imageUrl ?? article.featuredImage ?? article.coverImageUrl),
    visual,
    isFeatured: article.isFeatured === true,
  };
}

const visualVariants: VisualVariant[] = [
  'silhouette',
  'key',
  'message',
  'identity',
  'product',
  'fingerprint',
  'architecture',
  'email',
  'brokers',
  'minimalism',
  'lock',
];

function requireRecord(value: unknown, label: string): Record<string, unknown> {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) {
    throw new TemplateContentError(`${label} must be an object.`);
  }
  return value as Record<string, unknown>;
}

function readRequiredString(record: Record<string, unknown>, key: string, label: string): string {
  const value = readOptionalString(record, key);
  if (!value) throw new TemplateContentError(`${label} is required.`);
  return value;
}

function readOptionalString(record: Record<string, unknown>, key: string): string {
  const value = record[key];
  if (value === undefined || value === null) return '';
  if (typeof value !== 'string') {
    throw new TemplateContentError(`${key} must be a string.`);
  }
  return value.trim();
}

function readName(value: unknown): string {
  if (typeof value === 'string') return value.trim();
  if (value === undefined || value === null) return '';
  const record = requireRecord(value, 'Named content');
  return readOptionalString(record, 'name');
}

function normalizeSlug(value: string): string {
  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(value)) {
    throw new TemplateContentError('Article slug must contain lowercase letters, numbers, and single hyphens only.');
  }
  return value;
}

function formatPublishedDate(value: unknown): string {
  if (value === undefined || value === null || value === '') return '';
  if (typeof value !== 'string') throw new TemplateContentError('publishedAt must be a string.');
  const parsed = /^\d{4}-\d{2}-\d{2}$/.test(value) ? new Date(`${value}T00:00:00Z`) : new Date(value);
  if (Number.isNaN(parsed.getTime())) return '';
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(parsed);
}

function readTime(article: Record<string, unknown>, markdown: string): number {
  const value = article.readingTime ?? article.reading_time;
  if (value !== undefined && typeof value !== 'number' && typeof value !== 'string') {
    throw new TemplateContentError('readingTime must be a number or numeric string.');
  }
  const configured = Number(value);
  if (Number.isFinite(configured) && configured > 0) {
    return Math.ceil(configured);
  }
  const words = markdown.trim().split(/\s+/).filter(Boolean).length;
  return words ? Math.max(1, Math.ceil(words / 220)) : 1;
}

function deriveExcerpt(markdown: string): string {
  return markdown
    .replace(/<[^>]*>/g, ' ')
    .replace(/[#>*_`\[\]()~-]/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
    .slice(0, 180);
}

export function safeMediaUrl(value?: unknown): string | undefined {
  if (typeof value !== 'string') return undefined;
  const url = value.trim();
  if (!url) return undefined;
  if (/^https?:\/\/[^\s]+$/i.test(url)) {
    try {
      const parsed = new URL(url);
      return parsed.username || parsed.password ? undefined : url;
    } catch {
      return undefined;
    }
  }
  if (/^\/(?!\/)[^\s]*$/.test(url)) return url;
  return undefined;
}
