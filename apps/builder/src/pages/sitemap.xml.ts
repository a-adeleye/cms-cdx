import { loadBuildData } from '../content';
import { renderSitemap } from '../sitemap';

export async function GET() {
  const { site, articles } = loadBuildData();
  return new Response(renderSitemap(site, articles, true), {
    headers: { 'Content-Type': 'application/xml' },
  });
}
