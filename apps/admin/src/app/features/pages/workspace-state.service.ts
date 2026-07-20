import { Injectable, computed, inject, signal } from '@angular/core';
import {
  AdminStateSnapshot,
  ArticleRecord,
  ArticleStatus,
  AuthSession,
  AuthorRecord,
  BuildRecord,
  BuildType,
  CategoryRecord,
  LandingSectionRecord,
  MediaAssetRecord,
  SiteRecord,
  TemplateRecord,
  TagRecord,
} from './pages.model';
import { INITIAL_STATE } from './pages.seed';
import { AdminApiService, AISuggestionPayload, AISuggestionResponse } from './admin-api.service';
import { AuthTokenService } from './auth-token.service';

interface ArticleDraftInput {
  id?: string;
  title: string;
  slug: string;
  excerpt: string;
  contentMarkdown: string;
  coverImageUrl: string;
  seoTitle: string;
  seoDescription: string;
  canonicalUrl: string;
  authorId: string;
  categoryId: string;
  tagIds: string[];
  isFeatured: boolean;
  status: ArticleStatus;
}

interface SiteDraftInput {
  name: string;
  slug: string;
  domain: string;
  blogPath: string;
  templateKey: string;
}

interface CategoryDraftInput {
  id?: string;
  name: string;
  description: string;
}

interface AuthorDraftInput {
  id?: string;
  name: string;
  bio: string;
}

interface TagDraftInput {
  id?: string;
  name: string;
}

const EMPTY_SITE: SiteRecord = {
  id: '',
  name: 'No site selected',
  slug: '',
  domain: '',
  blogPath: '/articles',
  description: '',
  contentContext: 'standalone_blog',
  logoMediaId: '',
  logoUrl: '',
  faviconMediaId: '',
  faviconUrl: '',
  status: 'inactive',
  templateKey: 'default-blog',
  themeConfig: '{}',
  deployProvider: '',
  deployConfig: '{}',
  previewDeployProvider: '',
  previewDeployConfig: '{}',
  aiConfig: '{}',
  storageConfig: '{}',
  updatedAt: '',
};

const EMPTY_STATE: AdminStateSnapshot = {
  ...INITIAL_STATE,
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

@Injectable({
  providedIn: 'root',
})
export class WorkspaceStateService {
  private readonly api = inject(AdminApiService);
  private readonly tokenStore = inject(AuthTokenService);
  private readonly state = signal<AdminStateSnapshot>(EMPTY_STATE);
  private readonly loadingState = signal(true);
  private readonly errorState = signal<string | null>(null);

  readonly authSession = computed(() => this.state().authSession);
  readonly isAuthenticated = computed(() => Boolean(this.state().authSession));
  readonly loading = computed(() => this.loadingState());
  readonly error = computed(() => this.errorState());
  readonly sites = computed(() => this.state().sites.slice().sort((left, right) => left.name.localeCompare(right.name)));
  readonly templates = computed(() => this.state().templates.slice().sort((left, right) => left.name.localeCompare(right.name)));
  readonly selectedSiteId = computed(() => this.state().selectedSiteId);
  readonly selectedSite = computed(() => this.findSite(this.state().selectedSiteId) ?? this.state().sites[0] ?? EMPTY_SITE);
  readonly landingSections = computed(() =>
    this.state()
      .landingSections.filter((section) => section.siteId === this.selectedSite().id)
      .slice()
      .sort((left, right) => left.displayOrder - right.displayOrder),
  );
  readonly articles = computed(() =>
    this.state()
      .articles.filter((article) => article.siteId === this.selectedSite().id)
      .slice()
      .sort((left, right) => right.updatedAt.localeCompare(left.updatedAt)),
  );
  readonly authors = computed(() => this.state().authors.filter((author) => author.siteId === this.selectedSite().id));
  readonly categories = computed(() => this.state().categories.filter((category) => category.siteId === this.selectedSite().id));
  readonly tags = computed(() => this.state().tags.filter((tag) => tag.siteId === this.selectedSite().id));
  readonly mediaAssets = computed(() => this.state().mediaAssets.filter((asset) => asset.siteId === this.selectedSite().id));
  readonly builds = computed(() => this.state().builds.filter((build) => build.siteId === this.selectedSite().id));
  readonly selectedArticleId = computed(() => this.state().selectedArticleId);
  readonly selectedArticle = computed(() =>
    this.articles().find((article) => article.id === this.state().selectedArticleId) ?? null,
  );

  readonly dashboardStats = computed(() => {
    const articles = this.articles();
    const published = articles.filter((article) => article.status === 'published').length;
    const review = articles.filter((article) => article.status === 'review').length;
    const drafts = articles.filter((article) => article.status === 'draft').length;
    const builds = this.builds();
    const latestBuild = builds[0];

    return [
      { label: 'Published', value: String(published), detail: 'Live articles' },
      { label: 'Drafts', value: String(drafts), detail: 'Work in progress' },
      { label: 'Review queue', value: String(review), detail: 'Needs review' },
      { label: 'Deployments', value: latestBuild?.status === 'failed' ? 'Attention' : 'Healthy', detail: latestBuild ? 'Latest build recorded' : 'Ready to deploy' },
      { label: 'Media assets', value: String(this.mediaAssets().length), detail: 'Images and documents' },
    ];
  });

  constructor() {
    void this.bootstrap();
  }

  clearError(): void {
    this.errorState.set(null);
  }

  reportError(message: string): void {
    this.errorState.set(message);
  }

  async bootstrap(): Promise<void> {
    this.loadingState.set(true);
    this.errorState.set(null);

    try {
      const token = this.tokenStore.getToken();
      if (!token) {
        this.state.set(EMPTY_STATE);
        return;
      }

      this.tokenStore.setToken(token);
      const session = await this.api.me();
      const workspace = await this.api.loadWorkspace();
      this.applyWorkspace(workspace);
      this.state.update((state) => ({
        ...state,
        authSession: {
          email: session.user.email,
          fullName: session.user.fullName,
          role: session.user.role,
        },
      }));
    } catch {
      this.tokenStore.clear();
      this.state.set(EMPTY_STATE);
      this.errorState.set('Unable to restore your session. Please sign in again.');
    } finally {
      this.loadingState.set(false);
    }
  }

  async login(email: string, password: string): Promise<void> {
    this.errorState.set(null);
    const response = await this.api.login(email, password);
    this.tokenStore.setToken(response.token);
    await this.refreshSession();
  }

  async logout(): Promise<void> {
    this.errorState.set(null);
    try {
      if (this.tokenStore.getToken()) {
        await this.api.logout();
      }
    } finally {
      this.tokenStore.clear();
      this.state.set(EMPTY_STATE);
    }
  }

  async loadWorkspace(siteId?: string): Promise<void> {
    const workspace = await this.api.loadWorkspace(siteId || undefined);
    this.applyWorkspace(workspace);
  }

  async selectSite(siteId: string): Promise<void> {
    if (!siteId || siteId === this.state().selectedSiteId) {
      return;
    }

    await this.loadWorkspace(siteId);
  }

  async createSite(input: SiteDraftInput): Promise<SiteRecord> {
    const site = await this.api.createSite({
      name: input.name,
      slug: input.slug,
      domain: input.domain,
      blogPath: input.blogPath,
      description: '',
      contentContext: 'standalone_blog',
      logoMediaId: '',
      faviconMediaId: '',
      status: 'active',
      templateKey: input.templateKey,
      themeConfig: '{"tone":"professional"}',
      deployProvider: '',
      deployConfig: '{}',
      previewDeployProvider: '',
      previewDeployConfig: '{}',
      aiConfig: '{}',
      storageConfig: '{}',
    });

    await this.loadWorkspace(site.id);
    return site;
  }

  async updateSite(siteId: string, patch: Partial<SiteRecord>): Promise<SiteRecord> {
    const current = this.findSite(siteId) ?? this.selectedSite();
    const site = await this.api.updateSite(siteId, {
      name: patch.name ?? current.name,
      slug: patch.slug ?? current.slug,
      domain: patch.domain ?? current.domain,
      blogPath: patch.blogPath ?? current.blogPath,
      description: patch.description ?? current.description ?? '',
      contentContext: patch.contentContext ?? current.contentContext,
      logoMediaId: patch.logoMediaId ?? current.logoMediaId ?? '',
      faviconMediaId: patch.faviconMediaId ?? current.faviconMediaId ?? '',
      status: patch.status ?? current.status,
      templateKey: patch.templateKey ?? current.templateKey,
      themeConfig: patch.themeConfig ?? current.themeConfig,
      deployProvider: patch.deployProvider ?? current.deployProvider,
      deployConfig: patch.deployConfig ?? current.deployConfig,
      previewDeployProvider: patch.previewDeployProvider ?? current.previewDeployProvider,
      previewDeployConfig: patch.previewDeployConfig ?? current.previewDeployConfig,
      aiConfig: patch.aiConfig ?? current.aiConfig,
      storageConfig: patch.storageConfig ?? current.storageConfig,
    });

    await this.loadWorkspace(siteId);
    return site;
  }

  async updateSelectedSite(patch: Partial<SiteRecord>): Promise<void> {
    const site = this.selectedSite();
    if (!site.id) {
      return;
    }

    await this.updateSite(site.id, patch);
  }

  async createArticleDraft(): Promise<ArticleRecord> {
    const site = this.selectedSite();
    const article = await this.api.upsertArticle(site.id, {
      title: 'Untitled article',
      slug: `untitled-article-${Date.now()}`,
      excerpt: '',
      contentMarkdown: '# Untitled article\n\nStart writing here.',
      coverImageUrl: '',
      seoTitle: '',
      seoDescription: '',
      canonicalUrl: '',
      authorId: this.authors()[0]?.id ?? '',
      categoryId: this.categories()[0]?.id ?? '',
      tagIds: [],
      isFeatured: false,
      status: 'draft',
    });

    await this.loadWorkspace(site.id);
    this.state.update((state) => ({ ...state, selectedArticleId: article.id }));
    return article;
  }

  async saveCategory(input: CategoryDraftInput): Promise<CategoryRecord> {
    const site = this.selectedSite();
    if (!site.id) {
      throw new Error('No site selected.');
    }

    const category = input.id
      ? await this.api.updateCategory(site.id, input.id, {
          name: input.name,
          description: input.description,
        })
      : await this.api.createCategory(site.id, {
          name: input.name,
          description: input.description,
        });

    await this.loadWorkspace(site.id);
    return category;
  }

  async saveAuthor(input: AuthorDraftInput): Promise<AuthorRecord> {
    const site = this.selectedSite();
    if (!site.id) {
      throw new Error('No site selected.');
    }

    const author = input.id
      ? await this.api.updateAuthor(site.id, input.id, {
          name: input.name,
          bio: input.bio,
        })
      : await this.api.createAuthor(site.id, {
          name: input.name,
          bio: input.bio,
        });

    await this.loadWorkspace(site.id);
    return author;
  }

  async deleteAuthor(authorId: string): Promise<void> {
    const site = this.selectedSite();
    if (!site.id) {
      throw new Error('No site selected.');
    }

    await this.api.deleteAuthor(site.id, authorId);
    await this.loadWorkspace(site.id);
  }

  async deleteCategory(categoryId: string): Promise<void> {
    const site = this.selectedSite();
    if (!site.id) {
      throw new Error('No site selected.');
    }

    await this.api.deleteCategory(site.id, categoryId);
    await this.loadWorkspace(site.id);
  }

  async saveTag(input: TagDraftInput): Promise<TagRecord> {
    const site = this.selectedSite();
    if (!site.id) {
      throw new Error('No site selected.');
    }

    const tag = input.id
      ? await this.api.updateTag(site.id, input.id, {
          name: input.name,
        })
      : await this.api.createTag(site.id, {
          name: input.name,
        });

    await this.loadWorkspace(site.id);
    return tag;
  }

  async deleteTag(tagId: string): Promise<void> {
    const site = this.selectedSite();
    if (!site.id) {
      throw new Error('No site selected.');
    }

    await this.api.deleteTag(site.id, tagId);
    await this.loadWorkspace(site.id);
  }

  async saveArticle(input: ArticleDraftInput): Promise<ArticleRecord> {
    const site = this.selectedSite();
    const article = await this.api.upsertArticle(site.id, {
      id: input.id?.trim() || undefined,
      title: input.title,
      slug: input.slug,
      excerpt: input.excerpt,
      contentMarkdown: input.contentMarkdown,
      coverImageUrl: input.coverImageUrl,
      seoTitle: input.seoTitle,
      seoDescription: input.seoDescription,
      canonicalUrl: input.canonicalUrl,
      authorId: input.authorId,
      categoryId: input.categoryId,
      tagIds: input.tagIds,
      isFeatured: input.isFeatured,
      status: input.status,
    });

    await this.loadWorkspace(site.id);
    this.state.update((state) => ({ ...state, selectedArticleId: article.id }));
    return article;
  }

  async generateAISuggestion(input: AISuggestionPayload): Promise<AISuggestionResponse> {
    const site = this.selectedSite();
    if (!site.id) {
      throw new Error('No site selected.');
    }
    return this.api.generateAISuggestion(site.id, input);
  }

  async deleteArticle(articleId: string): Promise<void> {
    const site = this.selectedSite();
    if (!site.id) {
      throw new Error('No site selected.');
    }

    await this.api.deleteArticle(articleId);
    await this.loadWorkspace(site.id);
    this.state.update((state) => ({
      ...state,
      selectedArticleId: state.selectedArticleId === articleId ? null : state.selectedArticleId,
    }));
  }

  async selectArticle(articleId: string): Promise<void> {
    if (!this.articles().some((article) => article.id === articleId)) {
      return;
    }

    this.state.update((state) => ({
      ...state,
      selectedArticleId: articleId,
    }));
  }

  clearSelectedArticle(): void {
    this.state.update((state) => ({
      ...state,
      selectedArticleId: null,
    }));
  }

  async setArticleStatus(articleId: string, status: ArticleStatus): Promise<void> {
    const article = this.articles().find((entry) => entry.id === articleId);
    if (!article) {
      return;
    }

    await this.updateArticle(article, {
      id: article.id,
      title: article.title,
      slug: article.slug,
      excerpt: article.excerpt,
      contentMarkdown: article.contentMarkdown,
      coverImageUrl: article.coverImageUrl,
      seoTitle: article.seoTitle,
      seoDescription: article.seoDescription,
      canonicalUrl: article.canonicalUrl,
      authorId: article.authorId,
      categoryId: article.categoryId,
      tagIds: article.tagIds,
      isFeatured: article.isFeatured,
      status,
    });
  }

  async toggleFeatured(articleId: string): Promise<void> {
    const article = this.articles().find((entry) => entry.id === articleId);
    if (!article) {
      return;
    }

    await this.updateArticle(article, {
      id: article.id,
      title: article.title,
      slug: article.slug,
      excerpt: article.excerpt,
      contentMarkdown: article.contentMarkdown,
      coverImageUrl: article.coverImageUrl,
      seoTitle: article.seoTitle,
      seoDescription: article.seoDescription,
      canonicalUrl: article.canonicalUrl,
      authorId: article.authorId,
      categoryId: article.categoryId,
      tagIds: article.tagIds,
      isFeatured: !article.isFeatured,
      status: article.status,
    });
  }

  private async updateArticle(article: ArticleRecord, input: ArticleDraftInput): Promise<void> {
    await this.api.updateArticle(article.id, {
      id: input.id?.trim() || undefined,
      title: input.title,
      slug: input.slug,
      excerpt: input.excerpt,
      contentMarkdown: input.contentMarkdown,
      coverImageUrl: input.coverImageUrl,
      seoTitle: input.seoTitle,
      seoDescription: input.seoDescription,
      canonicalUrl: input.canonicalUrl,
      authorId: input.authorId,
      categoryId: input.categoryId,
      tagIds: input.tagIds,
      isFeatured: input.isFeatured,
      status: input.status,
    });

    await this.loadWorkspace(this.selectedSite().id);
  }

  async toggleLandingSection(sectionId: string): Promise<void> {
    const site = this.selectedSite();
    const section = this.landingSections().find((entry) => entry.id === sectionId);
    if (!section || !site.id) {
      return;
    }

    await this.api.updateLandingSection(site.id, sectionId, { isEnabled: !section.isEnabled });
    await this.loadWorkspace(site.id);
  }

  async moveLandingSection(sectionId: string, direction: 'up' | 'down'): Promise<void> {
    const site = this.selectedSite();
    const orderedSections = this.landingSections();
    const currentIndex = orderedSections.findIndex((section) => section.id === sectionId);
    const targetIndex = direction === 'up' ? currentIndex - 1 : currentIndex + 1;
    if (!site.id || currentIndex < 0 || targetIndex < 0 || targetIndex >= orderedSections.length) {
      return;
    }

    const nextSections = orderedSections.slice();
    [nextSections[currentIndex], nextSections[targetIndex]] = [nextSections[targetIndex], nextSections[currentIndex]];
    await this.api.reorderLandingSections(site.id, nextSections.map((section) => section.id));
    await this.loadWorkspace(site.id);
  }

  async triggerBuild(buildType: BuildType, articleIds: string[] = []): Promise<BuildRecord> {
    const site = this.selectedSite();
    try {
      return await this.api.createBuild(site.id, buildType, articleIds);
    } finally {
      await this.loadWorkspace(site.id);
    }
  }

  async uploadMediaFile(file: File, altText: string): Promise<MediaAssetRecord> {
    const site = this.selectedSite();
    if (!site.id) {
      throw new Error('No site selected.');
    }

    const media = await this.api.uploadMediaFile(site.id, file, altText);
    this.state.update((state) => ({
      ...state,
      mediaAssets: [media, ...state.mediaAssets.filter((asset) => asset.id !== media.id)],
    }));
    return media;
  }

  private async refreshSession(): Promise<void> {
    const session = await this.api.me();
    const workspace = await this.api.loadWorkspace();
    this.applyWorkspace(workspace);
    this.state.update((state) => ({
      ...state,
      authSession: {
        email: session.user.email,
        fullName: session.user.fullName,
        role: session.user.role,
      },
    }));
  }

  private applyWorkspace(workspace: Awaited<ReturnType<AdminApiService['loadWorkspace']>>): void {
    this.state.set({
      authSession: this.state().authSession,
      selectedSiteId: workspace.selectedSiteId,
      selectedArticleId: workspace.selectedArticleId || null,
      sites: workspace.sites,
      templates: workspace.templates,
      landingSections: workspace.landingSections,
      articles: workspace.articles,
      authors: workspace.authors,
      categories: workspace.categories,
      tags: workspace.tags,
      mediaAssets: workspace.mediaAssets,
      builds: workspace.builds,
    });

    if (workspace.user) {
      this.state.update((state) => ({
        ...state,
        authSession: {
          email: workspace.user.email,
          fullName: workspace.user.fullName,
          role: workspace.user.role,
        },
      }));
    }
  }

  private findSite(siteId: string): SiteRecord | undefined {
    return this.state().sites.find((site) => site.id === siteId);
  }
}
