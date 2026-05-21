export async function GET() {
  return new Response(`User-agent: *\nAllow: /\nSitemap: /sitemap.xml\n`, {
    headers: { 'Content-Type': 'text/plain' },
  });
}

