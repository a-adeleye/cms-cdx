import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { Router } from '@angular/router';
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

  onArticleFilterChange(value: ArticleStatus | 'all'): void {
    this.articleFilter.set(value);
  }

  async startArticle(): Promise<void> {
    try {
      const article = await this.state.createArticleDraft();
      await this.state.selectArticle(article.id);
      void this.router.navigate(['/articles/editor']);
    } catch (error) {
      this.reportActionError('Unable to create article draft.', error);
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

  async toggleFeatured(articleId: string): Promise<void> {
    try {
      await this.state.toggleFeatured(articleId);
    } catch (error) {
      this.reportActionError('Unable to update featured state.', error);
    }
  }

  async setStatus(articleId: string, status: ArticleStatus): Promise<void> {
    try {
      await this.state.setArticleStatus(articleId, status);
    } catch (error) {
      this.reportActionError('Unable to update article status.', error);
    }
  }

  private reportActionError(message: string, error: unknown): void {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    this.state.reportError(`${message}${detail}`);
  }
}
