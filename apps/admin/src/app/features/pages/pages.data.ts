import { WorkspacePageConfig } from './pages.model';

export const WORKSPACE_PAGES: WorkspacePageConfig[] = [
  {
    path: 'login',
    navLabel: 'Login',
    kind: 'login',
    eyebrow: 'Access',
    title: 'Login',
    primaryAction: { label: 'Open dashboard', path: '/dashboard' },
  },
  {
    path: 'dashboard',
    navLabel: 'Dashboard',
    kind: 'dashboard',
    eyebrow: 'Overview',
    title: 'Dashboard',
    primaryAction: { label: 'Review articles', path: '/content/articles' },
    secondaryAction: { label: 'Manage sites', path: '/configuration/sites' },
  },
  {
    path: 'publishing',
    navLabel: 'Publishing',
    kind: 'publishing',
    eyebrow: 'Delivery',
    title: 'Publishing',
    primaryAction: { label: 'Open articles', path: '/content/articles' },
  },
  {
    path: 'landing-page-editor',
    navLabel: 'Landing Page Editor',
    kind: 'landing-page-editor',
    eyebrow: 'Content design',
    title: 'Landing Page Editor',
    primaryAction: { label: 'Edit articles', path: '/content/articles' },
  },
  {
    path: 'articles',
    navLabel: 'Articles',
    kind: 'articles',
    eyebrow: 'Editorial',
    title: 'Articles',
    primaryAction: { label: 'Open editor', path: '/content/articles/new' },
  },
  {
    path: 'media-library',
    navLabel: 'Media',
    kind: 'media-library',
    eyebrow: 'Assets',
    title: 'Media',
    primaryAction: { label: 'Add media', path: '/content/media' },
  },
];
