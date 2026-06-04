import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { Router } from '@angular/router';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
import { ArticleStatus } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

type ArticleFilterOption = {
  value: ArticleStatus | 'all';
  label: string;
};

@Component({
  selector: 'app-articles-page',
  templateUrl: './articles-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css'],
  standalone: true,
  imports: [CommonModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ArticlesPageComponent {
  private readonly router = inject(Router);
  readonly state = inject(WorkspaceStateService);
  readonly feedback = createPageActionFeedback();

  readonly articleFilter = signal<ArticleStatus | 'all'>('all');
  readonly articleFilterOptions: ArticleFilterOption[] = [
    { value: 'all', label: 'All' },
    { value: 'draft', label: 'Draft' },
    { value: 'review', label: 'Review' },
    { value: 'published', label: 'Published' },
    { value: 'archived', label: 'Archived' },
  ];

  readonly filteredArticles = computed(() => {
    const articles = this.state.articles();
    const filter = this.articleFilter();
    if (filter === 'all') {
      return articles;
    }

    return articles.filter((article) => article.status === filter);
  });

  constructor() {
    const flashMessage = window.history.state?.flashMessage;
    if (typeof flashMessage === 'string' && flashMessage.trim()) {
      this.feedback.success(flashMessage);
    }
  }

  onArticleFilterChange(value: ArticleStatus | 'all'): void {
    this.articleFilter.set(value);
  }

  async startArticle(): Promise<void> {
    this.state.clearSelectedArticle();
    this.state.clearError();
    this.feedback.clear();
    void this.router.navigate(['/articles/editor']);
  }

  async openArticle(articleId: string): Promise<void> {
    try {
      this.state.clearError();
      await this.state.selectArticle(articleId);
      this.feedback.clear();
      void this.router.navigate(['/articles/editor']);
    } catch (error) {
      this.feedback.error(this.buildErrorMessage('Unable to open article.', error));
    }
  }

  async toggleFeatured(articleId: string): Promise<void> {
    try {
      this.state.clearError();
      this.feedback.loading('Updating featured state...');
      await this.state.toggleFeatured(articleId);
      this.feedback.success('Featured state updated successfully.');
    } catch (error) {
      this.feedback.error(this.buildErrorMessage('Unable to update featured state.', error));
    }
  }

  async setStatus(articleId: string, status: ArticleStatus): Promise<void> {
    try {
      this.state.clearError();
      this.feedback.loading('Updating article status...');
      await this.state.setArticleStatus(articleId, status);
      this.feedback.success('Article status updated successfully.');
    } catch (error) {
      this.feedback.error(this.buildErrorMessage('Unable to update article status.', error));
    }
  }

  private buildErrorMessage(message: string, error: unknown): string {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    return `${message}${detail}`;
  }
}
