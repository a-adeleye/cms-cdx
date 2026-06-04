import { AdminStateSnapshot } from './pages.model';

export const INITIAL_STATE: AdminStateSnapshot = {
  authSession: null,
  selectedSiteId: '',
  selectedArticleId: null,
  sites: [],
  templates: [],
  landingSections: [],
  articles: [],
  authors: [],
  categories: [],
  tags: [],
  mediaAssets: [],
  builds: [],
};
