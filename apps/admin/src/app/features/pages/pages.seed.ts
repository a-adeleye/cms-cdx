import { AdminStateSnapshot } from './pages.model';

export const INITIAL_STATE: AdminStateSnapshot = {
  authSession: null,
  selectedSiteId: '',
  selectedArticleId: null,
  sites: [],
  landingSections: [],
  articles: [],
  authors: [],
  categories: [],
  tags: [],
  mediaAssets: [],
  builds: [],
};
