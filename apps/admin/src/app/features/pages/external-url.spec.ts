import { canonicalSiteUrl, externalSiteUrl, siteHostname } from './external-url';

describe('external URL helpers', () => {
  it('preserves complete URLs and adds HTTPS to bare domains', () => {
    expect(externalSiteUrl('https://anonime.io')).toBe('https://anonime.io');
    expect(externalSiteUrl('anonime.io')).toBe('https://anonime.io');
    expect(externalSiteUrl('')).toBe('#');
  });

  it('normalizes canonical URLs and display hostnames', () => {
    expect(canonicalSiteUrl('https://anonime.io/blog/')).toBe('https://anonime.io/blog/');
    expect(siteHostname('https://anonime.io/blog')).toBe('anonime.io');
    expect(siteHostname(null)).toBe('example.com');
  });
});
