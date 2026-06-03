export const TEMPLATE_OPTIONS = [
  { value: 'default-blog', label: 'Default Blog' },
  { value: 'premium-saas', label: 'Premium SaaS' },
];

export const DEPLOY_PROVIDER_OPTIONS = [
  { value: 'none', label: 'None' },
  { value: 'netlify', label: 'Netlify' },
  { value: 'cloudflare', label: 'Cloudflare' },
  { value: 'firebase', label: 'Firebase' },
  { value: 's3', label: 'S3 / MinIO' },
];

const DEPLOY_CONFIG_TEMPLATES: Record<string, Record<string, string>> = {
  none: {
    provider: 'none',
  },
  netlify: {
    provider: 'netlify',
    siteId: '',
    tokenSecretRef: '',
    branch: 'main',
  },
  cloudflare: {
    provider: 'cloudflare',
    accountId: '',
    projectName: '',
    apiTokenSecretRef: '',
  },
  firebase: {
    provider: 'firebase',
    projectId: '',
    siteId: '',
    tokenSecretRef: '',
  },
  s3: {
    provider: 's3',
    bucket: '',
    region: '',
    prefix: '',
  },
};

export function defaultDeployConfigTemplate(provider: string): string {
  return JSON.stringify(DEPLOY_CONFIG_TEMPLATES[provider] ?? DEPLOY_CONFIG_TEMPLATES['none'], null, 2);
}

export function shouldReplaceDeployConfigTemplate(provider: string, value: string): boolean {
  const template = DEPLOY_CONFIG_TEMPLATES[provider] ?? DEPLOY_CONFIG_TEMPLATES['none'];
  if (!value.trim()) {
    return true;
  }

  try {
    const parsed = JSON.parse(value) as Record<string, unknown>;
    const keys = Object.keys(parsed);
    const templateKeys = Object.keys(template);
    if (keys.length !== templateKeys.length) {
      return keys.length === 0;
    }

    return templateKeys.every((key) => parsed[key] === template[key]);
  } catch {
    return true;
  }
}

export function isDeployConfigTemplate(value: string): boolean {
  if (!value.trim()) {
    return true;
  }

  try {
    const parsed = JSON.parse(value) as Record<string, unknown>;
    return Object.values(DEPLOY_CONFIG_TEMPLATES).some((template) => {
      const templateKeys = Object.keys(template);
      const keys = Object.keys(parsed);
      if (keys.length !== templateKeys.length) {
        return false;
      }
      return templateKeys.every((key) => parsed[key] === template[key]);
    });
  } catch {
    return true;
  }
}
