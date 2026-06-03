import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-publishing-page',
  templateUrl: './publishing-page.component.html',
  styleUrl: '../../features/pages/page-view.component.css',
  standalone: true,
  imports: [CommonModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PublishingPageComponent {
  readonly state = inject(WorkspaceStateService);
  readonly successMessage = signal<string | null>(null);
  readonly selectedArticleIds = signal<string[]>([]);

  readonly selectedArticleIdSet = computed(() => new Set(this.selectedArticleIds()));
  readonly selectedCount = computed(() => this.selectedArticleIds().length);
  readonly deploymentWarnings = computed(() => this.state.selectedSite().deploymentWarnings ?? []);
  readonly articles = computed(() => this.state.articles());
  readonly selectedArticles = computed(() => this.articles().filter((article) => this.selectedArticleIdSet().has(article.id)));
  readonly isAllSelected = computed(() => {
    const articles = this.articles();
    return articles.length > 0 && this.selectedCount() === articles.length;
  });

  readonly summaryMetrics = computed(() => [
    { label: 'Articles', value: String(this.articles().length), detail: 'All articles for the selected site.' },
    { label: 'Selected', value: String(this.selectedCount()), detail: 'Queued for your next build action.' },
    { label: 'Published', value: String(this.articles().filter((article) => article.status === 'published').length), detail: 'Live content that can be deployed.' },
  ]);

  async triggerPreviewBuild(): Promise<void> {
    try {
      const build = await this.state.triggerBuild('preview', this.selectedArticleIds());
      this.successMessage.set(this.buildSuccessMessage('Preview build', build.deployProvider, build.deployUrl));
    } catch (error) {
      this.reportActionError('Unable to start preview build.', error);
    }
  }

  async triggerPublishedBuild(): Promise<void> {
    try {
      for (const article of this.selectedArticles()) {
        if (article.status !== 'published') {
          await this.state.setArticleStatus(article.id, 'published');
        }
      }

      const build = await this.state.triggerBuild('published', this.selectedArticleIds());
      this.successMessage.set(this.buildSuccessMessage('Published build', build.deployProvider, build.deployUrl));
    } catch (error) {
      this.reportActionError('Unable to start published build.', error);
    }
  }

  toggleArticleSelection(articleId: string): void {
    const selectedIds = new Set(this.selectedArticleIds());
    if (selectedIds.has(articleId)) {
      selectedIds.delete(articleId);
    } else {
      selectedIds.add(articleId);
    }
    this.selectedArticleIds.set(Array.from(selectedIds));
  }

  toggleAllArticles(): void {
    if (this.isAllSelected()) {
      this.selectedArticleIds.set([]);
      return;
    }

    this.selectedArticleIds.set(this.articles().map((article) => article.id));
  }

  isArticleSelected(articleId: string): boolean {
    return this.selectedArticleIdSet().has(articleId);
  }

  private reportActionError(message: string, error: unknown): void {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    this.state.reportError(`${message}${detail}`);
  }

  private buildSuccessMessage(prefix: string, provider: string, deployUrl: string): string {
    const trimmedProvider = provider.trim();
    const trimmedUrl = deployUrl.trim();
    if (trimmedProvider && trimmedUrl) {
      return `${prefix} deployed to ${trimmedProvider}. Open ${trimmedUrl}`;
    }
    if (trimmedProvider) {
      return `${prefix} deployed to ${trimmedProvider}.`;
    }
    if (trimmedUrl) {
      return `${prefix} deployed successfully. Open ${trimmedUrl}`;
    }
    return `${prefix} deployed successfully.`;
  }
}
