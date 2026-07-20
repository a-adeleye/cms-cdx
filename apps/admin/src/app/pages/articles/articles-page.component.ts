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
  styleUrls: ['../../features/pages/page-view.component.css', './articles-page.component.css'],
  standalone: true,
  imports: [CommonModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ArticlesPageComponent {
  private readonly router = inject(Router);
  readonly state = inject(WorkspaceStateService);
  readonly feedback = createPageActionFeedback();

  readonly articleFilter = signal<ArticleStatus | 'all'>('all');
  readonly categoryFilter = signal('all');
  readonly authorFilter = signal('all');
  readonly articleSearch = signal('');
  readonly selectedArticleIds = signal<string[]>([]);
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
    const category = this.categoryFilter();
    const author = this.authorFilter();
    const search = this.articleSearch().trim().toLowerCase();
    return articles.filter((article) => {
      const matchesStatus = filter === 'all' || article.status === filter;
      const matchesCategory = category === 'all' || article.categoryId === category;
      const matchesAuthor = author === 'all' || article.authorId === author;
      const matchesSearch = !search || `${article.title} ${article.excerpt} ${article.slug}`.toLowerCase().includes(search);
      return matchesStatus && matchesCategory && matchesAuthor && matchesSearch;
    });
  });

  readonly hasActiveFilters = computed(
    () =>
      this.articleFilter() !== 'all' ||
      this.categoryFilter() !== 'all' ||
      this.authorFilter() !== 'all' ||
      this.articleSearch().trim().length > 0,
  );

  readonly allVisibleSelected = computed(() => {
    const visible = this.filteredArticles();
    return visible.length > 0 && visible.every((article) => this.selectedArticleIds().includes(article.id));
  });

  readonly someVisibleSelected = computed(() => {
    const visible = this.filteredArticles();
    return visible.some((article) => this.selectedArticleIds().includes(article.id)) && !this.allVisibleSelected();
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

  onSearch(value: string): void {
    this.articleSearch.set(value);
  }

  clearFilters(): void {
    this.articleSearch.set('');
    this.articleFilter.set('all');
    this.categoryFilter.set('all');
    this.authorFilter.set('all');
  }

  toggleSelection(articleId: string): void {
    const selected = new Set(this.selectedArticleIds());
    selected.has(articleId) ? selected.delete(articleId) : selected.add(articleId);
    this.selectedArticleIds.set([...selected]);
  }

  toggleAll(): void {
    const visibleIds = this.filteredArticles().map((article) => article.id);
    this.selectedArticleIds.set(this.selectedArticleIds().length === visibleIds.length ? [] : visibleIds);
  }

  isSelected(articleId: string): boolean {
    return this.selectedArticleIds().includes(articleId);
  }

  authorName(authorId: string): string {
    return this.state.authors().find((author) => author.id === authorId)?.name ?? 'Unassigned';
  }

  categoryName(categoryId: string): string {
    return this.state.categories().find((category) => category.id === categoryId)?.name ?? 'Uncategorized';
  }

  startArticle(): void {
    this.state.clearSelectedArticle();
    this.state.clearError();
    this.feedback.clear();
    void this.router.navigate(['/content/articles/new']);
  }

  openArticle(articleId: string): void {
    this.state.clearError();
    this.feedback.clear();
    void this.router.navigate(['/content/articles', articleId, 'edit']);
  }

  async deleteArticle(article: { id: string; title: string }): Promise<void> {
    if (!this.confirmDelete(article.title)) {
      return;
    }

    try {
      this.state.clearError();
      this.feedback.loading('Deleting article…');
      await this.state.deleteArticle(article.id);
      this.selectedArticleIds.set(this.selectedArticleIds().filter((id) => id !== article.id));
      this.feedback.success('Article deleted.');
    } catch (error) {
      this.reportActionError('Unable to delete article.', error);
    }
  }

  async deleteSelected(): Promise<void> {
    const ids = this.selectedArticleIds();
    if (!ids.length || !this.confirmDelete(`${ids.length} selected article${ids.length === 1 ? '' : 's'}`)) {
      return;
    }

    try {
      this.state.clearError();
      this.feedback.loading('Deleting articles…');
      for (const id of ids) {
        await this.state.deleteArticle(id);
      }
      this.selectedArticleIds.set([]);
      this.feedback.success('Selected articles deleted.');
    } catch (error) {
      this.reportActionError('Unable to delete selected articles.', error);
    }
  }

  private confirmDelete(label: string): boolean {
    return typeof globalThis.confirm !== 'function' || globalThis.confirm(`Delete ${label}? This cannot be undone.`);
  }

  private reportActionError(message: string, error: unknown): void {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    this.feedback.error(`${message}${detail}`);
  }
}
