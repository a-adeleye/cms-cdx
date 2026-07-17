import { TemplateRecord } from './pages.model';

export interface SelectOption {
  value: string;
  label: string;
}

export const DEFAULT_TEMPLATE_SLUG = 'default-blog';

export const DEPLOY_PROVIDER_OPTIONS = [
  { value: 'none', label: 'Local build only' },
  { value: 'cloudflare_pages', label: 'Cloudflare Pages' },
  { value: 'firebase', label: 'Firebase' },
  { value: 'git_repository', label: 'Repository branch' },
];

export function templateSelectOptions(templates: TemplateRecord[], currentValue = ''): SelectOption[] {
  const options = templates.map((template) => ({
    value: template.slug,
    label: template.name,
  }));

  const fallbackValue = currentValue.trim() || DEFAULT_TEMPLATE_SLUG;
  if (!options.some((option) => option.value === fallbackValue)) {
    options.unshift({
      value: fallbackValue,
      label: humanizeTemplateSlug(fallbackValue),
    });
  }

  return options;
}

const DEPLOY_CONFIG_TEMPLATES: Record<string, Record<string, string>> = {
  none: {
    provider: 'none',
  },
  cloudflare_pages: {
    projectName: '',
    productionBranch: 'main',
  },
  firebase: {
    siteId: '',
    serviceAccountSecretRef: '',
    publicUrl: '',
  },
  git_repository: {
    repositoryUrl: '',
    branch: 'main',
    contentPath: 'public/blog',
    tokenSecretRef: '',
    publicUrl: '',
  },
};

export function defaultDeployConfigTemplate(provider: string): string {
  return JSON.stringify(DEPLOY_CONFIG_TEMPLATES[provider] ?? DEPLOY_CONFIG_TEMPLATES['none'], null, 2);
}

function humanizeTemplateSlug(value: string): string {
  return value
    .split('-')
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(' ');
}
