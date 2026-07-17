import { DEPLOY_PROVIDER_OPTIONS, defaultDeployConfigTemplate } from './site-config-options';

describe('site deployment configuration options', () => {
  it('offers Cloudflare Pages with its supported configuration fields', () => {
    expect(DEPLOY_PROVIDER_OPTIONS).toContain({ value: 'cloudflare_pages', label: 'Cloudflare Pages' });
    expect(JSON.parse(defaultDeployConfigTemplate('cloudflare_pages'))).toEqual({
      projectName: '',
      productionBranch: 'main',
    });
  });
});
