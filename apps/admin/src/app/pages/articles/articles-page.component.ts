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

}
