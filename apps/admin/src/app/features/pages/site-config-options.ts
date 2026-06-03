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
    serviceAccountSecretRef: '',
  },
  s3: {
    provider: 's3',
    bucket: '',
    region: '',
    prefix: '',
  },
};

const AI_CONFIG_TEMPLATE = {
  provider: '',
  model: '',
  tone: '',
  brand_context: '',
};

const STORAGE_CONFIG_TEMPLATE = {
  provider: '',
  bucket: '',
  region: '',
  prefix: '',
  public_url: '',
};

export function defaultDeployConfigTemplate(provider: string): string {
  return JSON.stringify(DEPLOY_CONFIG_TEMPLATES[provider] ?? DEPLOY_CONFIG_TEMPLATES['none'], null, 2);
}

export function defaultAiConfigTemplate(): string {
  return JSON.stringify(AI_CONFIG_TEMPLATE, null, 2);
}

export function defaultStorageConfigTemplate(): string {
  return JSON.stringify(STORAGE_CONFIG_TEMPLATE, null, 2);
}

export function isJsonTemplate(value: string): boolean {
  if (!value.trim()) {
    return true;
  }

  try {
    const parsed = JSON.parse(value) as Record<string, unknown>;
    return (
      matchesTemplate(parsed, DEPLOY_CONFIG_TEMPLATES['none']) ||
      matchesTemplate(parsed, DEPLOY_CONFIG_TEMPLATES['netlify']) ||
      matchesTemplate(parsed, DEPLOY_CONFIG_TEMPLATES['cloudflare']) ||
      matchesTemplate(parsed, DEPLOY_CONFIG_TEMPLATES['firebase']) ||
      matchesTemplate(parsed, DEPLOY_CONFIG_TEMPLATES['s3']) ||
      matchesTemplate(parsed, AI_CONFIG_TEMPLATE) ||
      matchesTemplate(parsed, STORAGE_CONFIG_TEMPLATE)
    );
  } catch {
    return true;
  }
}

function matchesTemplate(parsed: Record<string, unknown>, template: Record<string, string>): boolean {
  const templateKeys = Object.keys(template);
  const keys = Object.keys(parsed);
  if (keys.length !== templateKeys.length) {
    return false;
  }
  return templateKeys.every((key) => parsed[key] === template[key]);
}
