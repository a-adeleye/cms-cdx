function pathSegments(pathname: string): string[] {
  return pathname.split('/').filter(Boolean);
}

// Generates a static-site-safe link. It works when the site is served from a
// domain root and when the CMS serves an artifact beneath /deployments/.
export function relativeSiteHref(currentPathname: string, target: string): string {
  const destination = new URL(target, 'https://template.invalid');
  const currentSegments = pathSegments(currentPathname);
  const currentDirectory = currentPathname.endsWith('/') ? currentSegments : currentSegments.slice(0, -1);
  const targetSegments = pathSegments(destination.pathname);
  let shared = 0;

  while (shared < currentDirectory.length && shared < targetSegments.length && currentDirectory[shared] === targetSegments[shared]) {
    shared += 1;
  }

  const relativeSegments = [
    ...Array.from({ length: currentDirectory.length - shared }, () => '..'),
    ...targetSegments.slice(shared),
  ];
  const path = relativeSegments.length === 0 ? './' : `${relativeSegments.join('/')}${destination.pathname.endsWith('/') ? '/' : ''}`;
  return `${path}${destination.search}${destination.hash}`;
}
