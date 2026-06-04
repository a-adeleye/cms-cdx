import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, effect, inject, signal } from '@angular/core';
import { AbstractControl, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
import { ArticleStatus } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-article-editor-page',
  templateUrl: './article-editor-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './article-editor-page.component.css'],
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ArticleEditorPageComponent {
  private readonly fb = inject(FormBuilder);
  private readonly router = inject(Router);
  readonly state = inject(WorkspaceStateService);

  readonly previewOpen = signal(false);
  readonly coverImageMode = signal<'url' | 'upload'>('url');
  readonly uploadingCoverImage = signal(false);
  readonly coverImageFileName = signal('');
  readonly coverImageError = signal<string | null>(null);
  readonly feedback = createPageActionFeedback();

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
    tagIds: [[] as string[]],
    isFeatured: [false],
    status: ['draft' as ArticleStatus, [Validators.required]],
    contentMarkdown: ['', [Validators.required, Validators.minLength(20)]],
  });

  constructor() {
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
            tagIds: [],
            isFeatured: false,
            status: 'draft',
            contentMarkdown: '',
          },
          { emitEvent: false },
        );
        this.coverImageMode.set('url');
        this.coverImageFileName.set('');
        this.coverImageError.set(null);
        this.feedback.clear();
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
          tagIds: article.tagIds,
          isFeatured: article.isFeatured,
          status: article.status,
          contentMarkdown: article.contentMarkdown,
        },
        { emitEvent: false },
      );
      this.coverImageMode.set(article.coverImageUrl ? 'url' : 'upload');
      this.coverImageFileName.set('');
      this.coverImageError.set(null);
      this.feedback.clear();
    });
  }

  setPreviewOpen(isOpen: boolean): void {
    this.previewOpen.set(isOpen);
  }

  setCoverImageMode(mode: 'url' | 'upload'): void {
    this.coverImageMode.set(mode);
    this.coverImageError.set(null);
  }

  async onCoverImageSelected(event: Event): Promise<void> {
    const target = event.target as HTMLInputElement | null;
    const file = target?.files?.[0];
    if (!file) {
      return;
    }

    if (!file.type.startsWith('image/')) {
      this.coverImageError.set('Please choose an image file.');
      if (target) {
        target.value = '';
      }
      return;
    }

    this.uploadingCoverImage.set(true);
    this.coverImageError.set(null);

    try {
      const media = await this.state.uploadMediaFile(file, this.articleForm.controls.title.value || 'Article cover');
      this.articleForm.controls.coverImageUrl.setValue(media.fileUrl);
      this.coverImageMode.set('url');
      this.coverImageFileName.set(file.name);
    } catch (error) {
      const detail = error instanceof Error && error.message ? error.message : 'Upload failed.';
      this.coverImageError.set(detail);
    } finally {
      this.uploadingCoverImage.set(false);
      if (target) {
        target.value = '';
      }
    }
  }

  async saveArticle(): Promise<void> {
    if (this.articleForm.invalid) {
      this.articleForm.markAllAsTouched();
      this.feedback.error('Please fix the highlighted fields before saving.');
      return;
    }

    const value = this.articleForm.getRawValue();
    try {
      this.feedback.loading('Saving article...');
      this.state.clearError();
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
        tagIds: value.tagIds,
        isFeatured: value.isFeatured,
        status: value.status,
      });

      await this.state.selectArticle(article.id);
      void this.router.navigate(['/articles', article.id], {
        state: { flashMessage: 'Article saved successfully.' },
      });
      this.state.clearError();
    } catch (error) {
      this.feedback.error(this.buildErrorMessage('Unable to save article.', error));
    }
  }

  controlHasError(control: AbstractControl | null, errorName: string): boolean {
    return Boolean(control?.hasError(errorName) && (control?.touched || control?.dirty));
  }

  private buildArticleSlug(title: string): string {
    return title
      .trim()
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }

  private buildErrorMessage(message: string, error: unknown): string {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    return `${message}${detail}`;
  }
}
