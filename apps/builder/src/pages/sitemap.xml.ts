export async function GET() {
  return new Response(
    `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`,
    {
      headers: { 'Content-Type': 'application/xml' },
    },
  );
}

