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
  categorySlug: string;
  author: string;
  authorInitials: string;
  publishedLabel: string;
  publishedLongLabel: string;
  publishedIso: string;
  readingTime: number;
  imageUrl?: string;
  visual: VisualVariant;
  tone: VisualTone;
  isFeatured: boolean;
};

export type VisualVariant = 'routing' | 'shield' | 'inbox' | 'signal' | 'mail' | 'nodes';

export type VisualTone = 'light' | 'alt' | 'dark';

export type SiteViewModel = {
  name: string;
  blogPath: string;
};

const featuredContent = `Every messaging platform sells the same promise: hand us your traffic and we will handle deliverability. It is a good pitch, and for the first few thousand messages it works. The trouble starts when your volume grows enough that the shortcuts underneath the promise begin to show.

Those shortcuts have names. Shared IP pools. Pooled short codes. Rotating sender IDs. Each one is a way of borrowing reputation you did not build, from a pool you do not control, alongside senders you will never meet.

## What renting actually buys you

When you send through a traditional provider, you are not really buying delivery. You are buying access to a reputation asset the provider owns and rations. Your messages leave from their IPs, carry their sending domain in the return path, and inherit whatever mailbox providers currently think of that infrastructure.

That works fine while the pool is healthy. It is a reasonable trade when you are sending your first campaign and have no reputation of your own to speak of. But it comes with a structural catch: **you cannot improve the asset, and you cannot leave with it.**

> Reputation you rent is reputation you have to re-earn every time you switch providers.

Two years of clean sending, low complaint rates, and careful list hygiene should compound into something. On a shared pool, it mostly does not. The pool's reputation moves as an aggregate, and your contribution to it is a rounding error.

## The shared pool problem

Here is the failure mode we see most often. A team sends transactional email — receipts, password resets, delivery notifications — on a shared pool. Complaint rates are near zero. Then, over a single week, open rates fall off a cliff and password resets start landing in spam.

Nothing changed on their side. What changed was that another tenant on the same pool ran a large acquisition campaign against a purchased list, and the mailbox providers responded to the IP range, not to the individual sender.

> The tell is timing: a deliverability drop that starts on a specific day, affects every mailbox provider at once, and has no corresponding change in your own sending behaviour.

Diagnosing this from inside a rented setup is close to impossible. You can see your bounce codes, but you cannot see the pool. You file a support ticket, wait, and are eventually moved to a different pool — which is to say, a different set of strangers.

## What owned sending changes

Owned sending inverts the arrangement. The domain is yours. The DKIM key is generated against your domain and lives in your DNS. The relay is one you chose, and can change without touching the reputation you have built, because the reputation is attached to \`mail.yourcompany.com\` and not to someone's IP block.

Three things follow from that:

- **Isolation.** Your delivery outcomes reflect your sending behaviour and nothing else. When something breaks, the cause is in your own data.
- **Portability.** Switching infrastructure becomes an operational change rather than a reputation reset.
- **Diagnosability.** Every bounce, defer, and complaint maps to a message you sent, a domain you own, and a route you chose.

Setting up a verified sending domain in Supromail is three DNS records and a verification check:

\`\`\`
curl -X POST https://api.supromail.com/v1/domains \\
  -H "Authorization: Bearer $SUPROMAIL_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{ "domain": "mail.yourcompany.com" }'

# returns the DKIM, SPF, and DMARC records to publish
\`\`\`

Once the records propagate, the domain moves to \`verified\` and every send from it is signed with a key only you control.

## The same argument for SMS

The email version of this story is well understood. The SMS version is newer, and the economics are sharper, because on SMS you are not only renting reputation — you are also renting numbers and paying a per-message toll on top of the carrier rate.

Pooled short codes have the same aggregation problem as shared IPs: filtering decisions land on the number, and the number is shared. Sending through your own SIMs means the number your customer sees is one you control, and the answer rate reflects your relationship with them rather than the pool's.

### Failover is the counter-argument, and it is solvable

The honest objection to owned SMS is resilience. One phone is a single point of failure in a way a carrier gateway is not. That is real, and it is why device pools and routing rules exist: when a device goes offline, the next healthy device in the pool picks up the queue, and the send continues.

## Where to start

You do not have to move everything at once, and you should not. The migration that works looks like this:

1. Verify one sending domain and move a low-risk transactional stream to it.
2. Watch bounce and complaint rates for two weeks against your existing baseline.
3. Move the remaining transactional traffic, then marketing.
4. For SMS, pair two devices before routing any production traffic, so failover exists from day one.

The point is not that managed providers are bad. It is that the asset they are managing on your behalf should belong to you — so that two years of careful sending compounds into something you can keep.`;

export const featuredDesignArticle: ArticleViewModel = designArticle(
  'Why owned sending beats rented reputation',
  'why-owned-sending-beats-rented-reputation',
  'Shared IP pools make your delivery rate a function of someone else’s worst customer. Here is what changes when the domains, the phones, and the routes are yours.',
  'Deliverability',
  'Ada Ajayi',
  '2026-06-12',
  8,
  'nodes',
  'light',
  { contentMarkdown: featuredContent, isFeatured: true },
);

export const designArticles: ArticleViewModel[] = [
  designArticle(
    'A practical guide to SMS failover across devices',
    'sms-failover-across-devices',
    'Device pools, health checks, and the routing rules that keep sends moving when a phone drops offline.',
    'Engineering',
    'Tunde Okafor',
    '2026-06-02',
    6,
    'routing',
    'alt',
  ),
  designArticle(
    'DKIM, SPF, and DMARC without the headache',
    'dkim-spf-and-dmarc-without-the-headache',
    'The three records every sending domain needs, what each one actually proves, and how to verify them fast.',
    'Deliverability',
    'Ada Ajayi',
    '2026-05-28',
    9,
    'mail',
    'light',
  ),
  designArticle(
    'Designing webhook consumers that never lose an event',
    'webhook-consumers-that-never-lose-an-event',
    'Signature verification, retries, and the queue-first pattern that keeps delivery state consistent.',
    'Engineering',
    'Tunde Okafor',
    '2026-05-19',
    7,
    'inbox',
    'dark',
  ),
  designArticle(
    'The hidden cost of per-message platform markup',
    'the-hidden-cost-of-per-message-markup',
    'We priced a year of SMS through three gateways and against direct carrier rates. The gap compounds.',
    'Product',
    'Nkem Eze',
    '2026-05-11',
    5,
    'signal',
    'light',
  ),
  designArticle(
    'Local SIM routing: what it is and when to use it',
    'local-sim-routing',
    'Sending in-country changes both cost and answer rates. Here is how to decide which traffic to route locally.',
    'Product',
    'Nkem Eze',
    '2026-04-24',
    6,
    'routing',
    'alt',
  ),
  designArticle(
    'How we think about workspace isolation',
    'how-we-think-about-workspace-isolation',
    'Scoped API keys, per-workspace log boundaries, and why we keep customer message history off shared storage.',
    'Security',
    'Ada Ajayi',
    '2026-04-15',
    4,
    'shield',
    'dark',
  ),
  designArticle(
    'Migrating from a legacy SMS gateway in a weekend',
    'migrating-from-a-legacy-sms-gateway',
    'A step-by-step cutover plan: dual-send, verify, shift traffic, then decommission.',
    'Guides',
    'Nkem Eze',
    '2026-04-03',
    7,
    'nodes',
    'light',
  ),
];

function designArticle(
  title: string,
  slug: string,
  excerpt: string,
  category: string,
  author: string,
  publishedAt: string,
  readingTime: number,
  visual: VisualVariant,
  tone: VisualTone,
  overrides: Partial<ArticleViewModel> = {},
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
    categorySlug: slugifyCategory(category),
    author,
    authorInitials: deriveInitials(author),
    publishedLabel: formatPublishedDate(publishedAt, 'short'),
    publishedLongLabel: formatPublishedDate(publishedAt, 'long'),
    publishedIso: publishedAt,
    readingTime,
    visual,
    tone,
    isFeatured: false,
    ...overrides,
  };
}

export class TemplateContentError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'TemplateContentError';
  }
}

export function resolveSite(value?: unknown): SiteViewModel {
  if (value === undefined) return { name: 'Supromail', blogPath: '/blog' };
  const site = requireRecord(value, 'Site configuration');
  const configuredName = readOptionalString(site, 'name');
  const blogPath = readOptionalString(site, 'blogPath');
  return {
    name: configuredName && configuredName !== 'Example Site' ? configuredName : 'Supromail',
    blogPath: normalizeBlogPath(blogPath),
  };
}

export function normalizeBlogPath(value?: string): string {
  const path = value?.trim() || '/blog';
  if (path === '/' || !/^\/[a-z0-9/_-]*$/i.test(path) || path.startsWith('//')) {
    return '/blog';
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
  return value === undefined ? featuredDesignArticle : adaptArticle(value, 0, true);
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
  const category = readOptionalString(article, 'categoryName') || readName(article.category) || 'Messaging';
  const author = readOptionalString(article, 'authorName') || readName(article.author) || 'Supromail team';
  const publishedAt = readPublishedAt(article);
  return {
    title,
    slug,
    excerpt,
    contentMarkdown,
    seoTitle: readOptionalString(article, 'seoTitle') || title,
    seoDescription: readOptionalString(article, 'seoDescription') || excerpt,
    canonicalUrl: safeMediaUrl(article.canonicalUrl) || '',
    category,
    categorySlug: slugifyCategory(category),
    author,
    authorInitials: deriveInitials(author),
    publishedLabel: formatPublishedDate(publishedAt, 'short'),
    publishedLongLabel: formatPublishedDate(publishedAt, 'long'),
    publishedIso: /^\d{4}-\d{2}-\d{2}$/.test(publishedAt) ? publishedAt : '',
    readingTime: readTime(article, contentMarkdown),
    imageUrl: safeMediaUrl(article.imageUrl ?? article.featuredImage ?? article.coverImageUrl),
    visual: visualVariants[index % visualVariants.length],
    tone: visualTones[index % visualTones.length],
    isFeatured: article.isFeatured === true,
  };
}

const visualVariants: VisualVariant[] = ['nodes', 'routing', 'mail', 'inbox', 'signal', 'shield'];

const visualTones: VisualTone[] = ['light', 'alt', 'dark'];

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

export function slugifyCategory(value: string): string {
  return value.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '') || 'general';
}

export function deriveInitials(value: string): string {
  const words = value.trim().split(/\s+/).filter(Boolean);
  const initials = words.slice(0, 2).map((word) => word[0]).join('');
  return initials.toUpperCase() || 'SM';
}

function readPublishedAt(article: Record<string, unknown>): string {
  const value = article.publishedAt ?? article.published_at;
  if (value === undefined || value === null || value === '') return '';
  if (typeof value !== 'string') throw new TemplateContentError('publishedAt must be a string.');
  return value.trim();
}

function formatPublishedDate(value: string, style: 'short' | 'long'): string {
  if (!value) return '';
  const parsed = /^\d{4}-\d{2}-\d{2}$/.test(value) ? new Date(`${value}T00:00:00Z`) : new Date(value);
  if (Number.isNaN(parsed.getTime())) return '';
  return new Intl.DateTimeFormat('en-GB', {
    day: '2-digit',
    month: style === 'long' ? 'long' : 'short',
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

export function collectCategories(articles: ArticleViewModel[]): Array<{ name: string; slug: string }> {
  const seen = new Map<string, string>();
  for (const article of articles) {
    if (!seen.has(article.categorySlug)) seen.set(article.categorySlug, article.category);
  }
  return Array.from(seen, ([slug, name]) => ({ slug, name }));
}
