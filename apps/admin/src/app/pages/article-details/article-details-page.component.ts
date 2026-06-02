import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, OnInit, inject, signal } from '@angular/core';
import { Router, RouterModule, ActivatedRoute } from '@angular/router';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-article-details-page',
  templateUrl: './article-details-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css'],
  standalone: true,
  imports: [CommonModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ArticleDetailsPageComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  readonly state = inject(WorkspaceStateService);
  readonly successMessage = signal<string | null>(null);

  async ngOnInit(): Promise<void> {
    const articleId = this.route.snapshot.paramMap.get('articleId')?.trim();
    if (!articleId) {
      void this.router.navigate(['/articles']);
      return;
    }

    const flashMessage = window.history.state?.flashMessage;
    if (typeof flashMessage === 'string' && flashMessage.trim()) {
      this.successMessage.set(flashMessage);
    }

    await this.state.selectArticle(articleId);
    if (!this.state.selectedArticle()) {
      void this.router.navigate(['/articles']);
    }
  }

  async deleteArticle(): Promise<void> {
    const article = this.state.selectedArticle();
    if (!article) {
      return;
    }

    if (!window.confirm(`Delete "${article.title}"? This cannot be undone.`)) {
      return;
    }

    try {
      this.state.clearError();
      await this.state.deleteArticle(article.id);
      void this.router.navigate(['/articles'], {
        state: { flashMessage: 'Article deleted successfully.' },
      });
    } catch (error) {
      const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
      this.state.reportError(`Unable to delete article.${detail}`);
    }
  }

  async editArticle(): Promise<void> {
    const articleId = this.state.selectedArticle()?.id;
    if (!articleId) {
      return;
    }

    await this.state.selectArticle(articleId);
    void this.router.navigate(['/articles/editor']);
  }

  authorName(authorId: string): string {
    return this.state.authors().find((author) => author.id === authorId)?.name ?? '—';
  }

  categoryName(categoryId: string): string {
    return this.state.categories().find((category) => category.id === categoryId)?.name ?? '—';
  }

  tagName(tagId: string): string {
    return this.state.tags().find((tag) => tag.id === tagId)?.name ?? tagId;
  }
}
