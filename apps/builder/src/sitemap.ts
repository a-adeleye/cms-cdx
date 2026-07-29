import type { BuildArticle, BuildSite } from './content';

const articlesPerPage = 6;

export function renderSitemap(site: BuildSite, articles: BuildArticle[], includeLanding: boolean) {
  const baseURL = new URL(site.domain || process.env.CMS_SITE_URL || 'http://localhost:8081');
  const blogPath = withTrailingSlash(site.blogPath || '/articles');
  const paths = new Set([blogPath, `${blogPath}all/`]);
  if (includeLanding) {
    paths.add('/');
  }

  for (let page = 2; page <= Math.ceil(articles.length / articlesPerPage); page += 1) {
    paths.add(`${blogPath}page/${page}/`);
  }

  const urls = [
    ...[...paths].map((path) => ({ location: absoluteURL(baseURL, path) })),
    ...articles.map((article) => ({
      location: absoluteURL(baseURL, `${blogPath}${encodeURIComponent(article.slug)}/`),
      lastModified: article.publishedAt,
    })),
  ];

  const entries = urls.map(({ location, lastModified }) => {
    const lastModifiedTag = lastModified ? `<lastmod>${escapeXML(lastModified)}</lastmod>` : '';
    return `  <url><loc>${escapeXML(location)}</loc>${lastModifiedTag}</url>`;
  });
  return `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${entries.join('\n')}\n</urlset>`;
}

function withTrailingSlash(path: string) {
  const normalized = `/${path.trim().replace(/^\/+|\/+$/g, '')}`;
  return normalized === '/' ? normalized : `${normalized}/`;
}

function absoluteURL(baseURL: URL, path: string) {
  return new URL(path, baseURL).toString();
}

function escapeXML(value: string) {
  return value.replace(/[<>&'\"]/g, (character) => ({
    '<': '&lt;',
    '>': '&gt;',
    '&': '&amp;',
    "'": '&apos;',
    '"': '&quot;',
  })[character] ?? character);
}
