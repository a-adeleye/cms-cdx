export function externalSiteUrl(domain: string | null | undefined): string {
  const value = domain?.trim();
  if (!value) {
    return '#';
  }

  return /^https?:\/\//i.test(value) ? value : `https://${value}`;
}

export function canonicalSiteUrl(domain: string | null | undefined): string {
  const url = externalSiteUrl(domain);
  return url === '#' ? '' : `${url.replace(/\/$/, '')}/`;
}

export function siteHostname(domain: string | null | undefined): string {
  return domain?.trim().replace(/^https?:\/\//i, '').replace(/\/.*$/, '') || 'example.com';
}
