import { ChangeDetectionStrategy, Component, computed, effect, inject } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { filter, map, startWith } from 'rxjs';
import { WORKSPACE_PAGES } from './pages.data';
import { WorkspacePageConfig } from './pages.model';
import { SummaryMetric } from './page-view.types';
import { WorkspaceStateService } from './workspace-state.service';

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
  readonly currentUrl = toSignal(
    this.router.events.pipe(
      filter((event): event is NavigationEnd => event instanceof NavigationEnd),
      startWith(null),
      map(() => this.router.url),
    ),
    { initialValue: this.router.url },
  );
  readonly isSettingsRoot = computed(() => this.page().kind === 'settings' && this.currentUrl() === '/settings');

  readonly loginForm = this.fb.nonNullable.group({
    email: ['admin@example.com', [Validators.required, Validators.email]],
    password: ['admin123', [Validators.required, Validators.minLength(6)]],
  });

  readonly supportMetrics = computed<SummaryMetric[]>(() => {
    const page = this.page();
    const site = this.state.selectedSite();

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

  async openArticle(articleId: string): Promise<void> {
    try {
      await this.state.selectArticle(articleId);
      void this.router.navigate(['/articles/editor']);
    } catch (error) {
      this.reportActionError('Unable to open article.', error);
    }
  }

  async createArticleDraft(): Promise<void> {
    try {
      const article = await this.state.createArticleDraft();
      await this.state.selectArticle(article.id);
      void this.router.navigate(['/articles/editor']);
    } catch (error) {
      this.reportActionError('Unable to create article draft.', error);
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

  async onSiteSelectionChange(siteId: string): Promise<void> {
    try {
      await this.state.selectSite(siteId);
    } catch (error) {
      this.reportActionError('Unable to switch sites.', error);
    }
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
}
