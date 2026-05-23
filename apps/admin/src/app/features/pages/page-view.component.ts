import { ChangeDetectionStrategy, Component, computed, effect, inject, signal } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { map } from 'rxjs';
import { WORKSPACE_PAGES } from './pages.data';
import { ArticleStatus, WorkspacePageConfig } from './pages.model';
import { SummaryMetric } from './page-view.types';
import { WorkspaceStateService } from './workspace-state.service';

type ArticleFilterOption = {
  value: ArticleStatus | 'all';
  label: string;
};

@Component({
  selector: 'app-page-view',
  templateUrl: './page-view.component.html',
  styleUrl: './page-view.component.css',
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
  host: {
    class: 'page-view',
  },
})
export class PageViewComponent {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);

  readonly page = toSignal(
    this.route.data.pipe(map((data) => (data['page'] as WorkspacePageConfig | undefined) ?? WORKSPACE_PAGES[1])),
    { initialValue: WORKSPACE_PAGES[1] },
  );

  readonly articleFilter = signal<ArticleStatus | 'all'>('all');
  readonly articleFilterOptions: ArticleFilterOption[] = [
    { value: 'all', label: 'All' },
    { value: 'draft', label: 'Draft' },
    { value: 'review', label: 'Review' },
    { value: 'published', label: 'Published' },
    { value: 'archived', label: 'Archived' },
  ];
  readonly previewOpen = signal(false);

  readonly loginForm = this.fb.nonNullable.group({
    email: ['admin@example.com', [Validators.required, Validators.email]],
    password: ['admin123', [Validators.required, Validators.minLength(6)]],
  });

  readonly siteCreateForm = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    slug: ['', [Validators.required, Validators.pattern(/^[a-z0-9-]+$/)]],
    domain: ['', [Validators.required]],
    blogPath: ['/articles', [Validators.required]],
    templateKey: ['default-blog', [Validators.required]],
  });

  readonly siteEditForm = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    slug: ['', [Validators.required, Validators.pattern(/^[a-z0-9-]+$/)]],
    domain: ['', [Validators.required]],
    blogPath: ['/articles', [Validators.required]],
    templateKey: ['default-blog', [Validators.required]],
    status: ['active' as 'active' | 'inactive', [Validators.required]],
  });

  readonly siteSettingsForm = this.fb.nonNullable.group({
    themeConfig: ['', [Validators.required]],
    deployProvider: ['none', [Validators.required]],
    deployConfig: ['', [Validators.required]],
    aiConfig: ['', [Validators.required]],
    storageConfig: ['', [Validators.required]],
  });

  readonly articleForm = this.fb.nonNullable.group({
    id: [''],
    title: ['', [Validators.required, Validators.minLength(3)]],
    excerpt: ['', [Validators.required, Validators.minLength(12)]],
    coverImageUrl: [''],
    seoTitle: ['', [Validators.required]],
    seoDescription: ['', [Validators.required]],
    canonicalUrl: [''],
    authorId: ['', [Validators.required]],
    categoryId: ['', [Validators.required]],
    tagIds: [''],
    isFeatured: [false],
    contentMarkdown: ['', [Validators.required, Validators.minLength(20)]],
  });

  readonly filteredArticles = computed(() => {
    const articles = this.state.articles();
    const filter = this.articleFilter();
    if (filter === 'all') {
      return articles;
    }

    return articles.filter((article) => article.status === filter);
  });

  readonly supportMetrics = computed<SummaryMetric[]>(() => {
    const page = this.page();
    const site = this.state.selectedSite();

    if (page.kind === 'authors') {
      return [
        { label: 'Authors', value: String(this.state.authors().length), detail: 'Contributors available for attribution.' },
        { label: 'Editors', value: '1', detail: 'Editorial operators with review access.' },
      ];
    }

    if (page.kind === 'categories') {
      return [
        { label: 'Categories', value: String(this.state.categories().length), detail: 'Taxonomy groups for the selected site.' },
        { label: 'Articles', value: String(this.state.articles().length), detail: 'Content organized by those categories.' },
      ];
    }

    if (page.kind === 'tags') {
      return [
        { label: 'Tags', value: String(this.state.tags().length), detail: 'Reusable labels for content discovery.' },
        { label: 'Featured', value: String(this.state.articles().filter((article) => article.isFeatured).length), detail: 'Tagged stories promoted on landing pages.' },
      ];
    }

    if (page.kind === 'media-library') {
      return [
        { label: 'Media assets', value: String(this.state.mediaAssets().length), detail: 'Uploaded assets for the current site.' },
        { label: 'Storage provider', value: site?.storageConfig.includes('r2') ? 'R2' : 'MinIO', detail: 'Configured storage backend for uploads.' },
      ];
    }

    if (page.kind === 'ai-assistant') {
      return [
        { label: 'Drafts', value: String(this.state.articles().filter((article) => article.status === 'draft').length), detail: 'AI output must remain draft or review.' },
        { label: 'Review queue', value: String(this.state.articles().filter((article) => article.status === 'review').length), detail: 'Human approval still required.' },
      ];
    }

    if (page.kind === 'builds') {
      return [
        { label: 'Builds', value: String(this.state.builds().length), detail: 'Preview and published build history.' },
        { label: 'Latest status', value: this.state.builds()[0]?.status ?? 'idle', detail: 'Current publish health.' },
      ];
    }

    if (page.kind === 'deployment-settings') {
      return [
        { label: 'Provider', value: site?.deployProvider || 'none', detail: 'Deployment target for the selected site.' },
        { label: 'Template', value: site?.templateKey || 'default-blog', detail: 'Current output template key.' },
      ];
    }

    if (page.kind === 'landing-page-editor') {
      return [
        { label: 'Sections', value: String(this.state.landingSections().length), detail: 'Landing page blocks for the selected site.' },
        { label: 'Enabled', value: String(this.state.landingSections().filter((section) => section.isEnabled).length), detail: 'Live sections included in the build.' },
      ];
    }

    if (page.kind === 'sites') {
      return [
        { label: 'Sites', value: String(this.state.sites().length), detail: 'Managed websites in this CMS.' },
        { label: 'Active', value: String(this.state.sites().filter((entry) => entry.status === 'active').length), detail: 'Sites currently publishing.' },
      ];
    }

    if (page.kind === 'site-settings') {
      return [
        { label: 'Domain', value: site?.domain ?? 'unset', detail: 'Canonical public URL for the selected site.' },
        { label: 'Blog path', value: site?.blogPath ?? '/articles', detail: 'Public article listing path.' },
      ];
    }

    if (page.kind === 'articles' || page.kind === 'article-editor') {
      return [
        { label: 'Drafts', value: String(this.state.articles().filter((article) => article.status === 'draft').length), detail: 'Ready for editing.' },
        { label: 'Published', value: String(this.state.articles().filter((article) => article.status === 'published').length), detail: 'Visible to site visitors.' },
      ];
    }

    return [
      { label: 'Selected site', value: site?.name ?? 'None', detail: 'The site context that drives this admin view.' },
      { label: 'Articles', value: String(this.state.articles().length), detail: 'Content available for the selected site.' },
    ];
  });

  readonly supportHighlights = computed(() => {
    const page = this.page();
    switch (page.kind) {
      case 'settings':
        return ['Shortcut links', 'Content tools', 'Publishing controls'];
      case 'authors':
        return ['Contributor attribution', 'Role management', 'Reusable profiles'];
      case 'categories':
        return ['Taxonomy structure', 'Stable navigation', 'Editorial filtering'];
      case 'tags':
        return ['Reusable labels', 'Search facets', 'Content linking'];
      case 'media-library':
        return ['S3-compatible uploads', 'Alt text validation', 'Asset reuse'];
      case 'ai-assistant':
        return ['Draft only output', 'Human review required', 'Provider-based generation'];
      case 'builds':
        return ['Preview builds', 'Published builds', 'Rollback-friendly history'];
      case 'deployment-settings':
        return ['Provider config', 'Secret references', 'Publish targets'];
      case 'landing-page-editor':
        return ['Section ordering', 'Enable/disable controls', 'Template-driven rendering'];
      default:
        return ['Multi-site context', 'Content ownership', 'Static site output'];
    }
  });

  constructor() {
    effect(() => {
      if (this.state.loading()) {
        return;
      }

      const page = this.page();
      if (!this.state.isAuthenticated() && page.kind !== 'login') {
        void this.router.navigate(['/login']);
        return;
      }

      if (this.state.isAuthenticated() && page.kind === 'login') {
        void this.router.navigate(['/dashboard']);
      }
    });

    effect(() => {
      const site = this.state.selectedSite();
      if (!site) {
        return;
      }

      this.siteEditForm.reset(
        {
          name: site.name,
          slug: site.slug,
          domain: site.domain,
          blogPath: site.blogPath,
          templateKey: site.templateKey,
          status: site.status,
        },
        { emitEvent: false },
      );

      this.siteSettingsForm.reset(
        {
          themeConfig: site.themeConfig,
          deployProvider: site.deployProvider,
          deployConfig: site.deployConfig,
          aiConfig: site.aiConfig,
          storageConfig: site.storageConfig,
        },
        { emitEvent: false },
      );
    });

    effect(() => {
      const article = this.state.selectedArticle();
      if (!article) {
        this.articleForm.reset(
          {
            id: '',
            title: '',
            excerpt: '',
            coverImageUrl: '',
            seoTitle: '',
            seoDescription: '',
            canonicalUrl: '',
            authorId: this.state.authors()[0]?.id ?? '',
            categoryId: this.state.categories()[0]?.id ?? '',
            tagIds: '',
            isFeatured: false,
            contentMarkdown: '',
          },
          { emitEvent: false },
        );
        return;
      }

      this.articleForm.reset(
        {
          id: article.id,
          title: article.title,
          excerpt: article.excerpt,
          coverImageUrl: article.coverImageUrl,
          seoTitle: article.seoTitle,
          seoDescription: article.seoDescription,
          canonicalUrl: article.canonicalUrl,
          authorId: article.authorId,
          categoryId: article.categoryId,
          tagIds: article.tagIds.join(', '),
          isFeatured: article.isFeatured,
          contentMarkdown: article.contentMarkdown,
        },
        { emitEvent: false },
      );
    });
  }

  private reportActionError(message: string, error: unknown): void {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    this.state.reportError(`${message}${detail}`);
  }

  async signIn(): Promise<void> {
    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      return;
    }

    const { email, password } = this.loginForm.getRawValue();
    try {
      await this.state.login(email, password);
      void this.router.navigate(['/dashboard']);
    } catch (error) {
      this.reportActionError('Unable to sign in.', error);
    }
  }

  async createSite(): Promise<void> {
    if (this.siteCreateForm.invalid) {
      this.siteCreateForm.markAllAsTouched();
      return;
    }

    try {
      await this.state.createSite(this.siteCreateForm.getRawValue());
      this.siteCreateForm.reset({
        name: '',
        slug: '',
        domain: '',
        blogPath: '/articles',
        templateKey: 'default-blog',
      });
    } catch (error) {
      this.reportActionError('Unable to create site.', error);
    }
  }

  async saveSite(): Promise<void> {
    if (this.siteEditForm.invalid) {
      this.siteEditForm.markAllAsTouched();
      return;
    }

    const { status, ...site } = this.siteEditForm.getRawValue();
    try {
      await this.state.updateSelectedSite({
        ...site,
        status,
      });
    } catch (error) {
      this.reportActionError('Unable to save site.', error);
    }
  }

  async saveSiteSettings(): Promise<void> {
    if (this.siteSettingsForm.invalid) {
      this.siteSettingsForm.markAllAsTouched();
      return;
    }

    const { themeConfig, deployProvider, deployConfig, aiConfig, storageConfig } = this.siteSettingsForm.getRawValue();
    try {
      await this.state.updateSelectedSite({
        themeConfig,
        deployProvider,
        deployConfig,
        aiConfig,
        storageConfig,
      });
    } catch (error) {
      this.reportActionError('Unable to save site settings.', error);
    }
  }

  async openArticle(articleId: string): Promise<void> {
    try {
      await this.state.selectArticle(articleId);
      void this.router.navigate(['/article-editor']);
    } catch (error) {
      this.reportActionError('Unable to open article.', error);
    }
  }

  openArticleEditor(): void {
    this.state.clearSelectedArticle();
    void this.router.navigate(['/article-editor']);
  }

  async createArticleDraft(): Promise<void> {
    try {
      const article = await this.state.createArticleDraft();
      await this.state.selectArticle(article.id);
      void this.router.navigate(['/article-editor']);
    } catch (error) {
      this.reportActionError('Unable to create article draft.', error);
    }
  }

  async saveArticle(status: ArticleStatus = 'draft'): Promise<void> {
    if (this.articleForm.invalid) {
      this.articleForm.markAllAsTouched();
      return;
    }

    const value = this.articleForm.getRawValue();
    try {
      const article = await this.state.saveArticle({
        id: value.id,
        title: value.title,
        slug: this.buildArticleSlug(value.title),
        excerpt: value.excerpt,
        contentMarkdown: value.contentMarkdown,
        coverImageUrl: value.coverImageUrl,
        seoTitle: value.seoTitle,
        seoDescription: value.seoDescription,
        canonicalUrl: value.canonicalUrl,
        authorId: value.authorId,
        categoryId: value.categoryId,
        tagIds: value.tagIds
          .split(',')
          .map((tag) => tag.trim())
          .filter((tag) => tag.length > 0),
        isFeatured: value.isFeatured,
        status,
      });

      await this.state.selectArticle(article.id);
    } catch (error) {
      this.reportActionError('Unable to save article.', error);
    }
  }

  async triggerPreviewBuild(): Promise<void> {
    try {
      await this.state.triggerBuild('preview');
    } catch (error) {
      this.reportActionError('Unable to start preview build.', error);
    }
  }

  async triggerPublishedBuild(): Promise<void> {
    try {
      await this.state.triggerBuild('published');
    } catch (error) {
      this.reportActionError('Unable to start published build.', error);
    }
  }

  setPreviewOpen(isOpen: boolean): void {
    this.previewOpen.set(isOpen);
  }

  async onSiteSelectionChange(siteId: string): Promise<void> {
    try {
      await this.state.selectSite(siteId);
    } catch (error) {
      this.reportActionError('Unable to switch sites.', error);
    }
  }

  onArticleFilterChange(value: ArticleStatus | 'all'): void {
    this.articleFilter.set(value);
  }

  async toggleSection(sectionId: string): Promise<void> {
    try {
      await this.state.toggleLandingSection(sectionId);
    } catch (error) {
      this.reportActionError('Unable to update landing section.', error);
    }
  }

  async moveSection(sectionId: string, direction: 'up' | 'down'): Promise<void> {
    try {
      await this.state.moveLandingSection(sectionId, direction);
    } catch (error) {
      this.reportActionError('Unable to reorder landing sections.', error);
    }
  }

  async uploadSampleMedia(): Promise<void> {
    try {
      await this.state.uploadMedia('new-asset.jpg', 'https://images.example/new-asset.jpg', 'New uploaded asset');
    } catch (error) {
      this.reportActionError('Unable to upload media.', error);
    }
  }

  private buildArticleSlug(title: string): string {
    return title
      .trim()
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }
}
