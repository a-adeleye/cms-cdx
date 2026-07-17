import { CommonModule, NgOptimizedImage } from '@angular/common';
import { ChangeDetectionStrategy, Component, ElementRef, ViewChild, effect, inject, signal } from '@angular/core';
import { AbstractControl, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
import { ArticleStatus } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { RichTextEditorComponent } from './rich-text-editor.component';

type EditorTool = 'write' | 'ai' | 'preview' | 'history';
type EditorMode = 'create' | 'edit';

@Component({
  selector: 'app-article-editor-page',
  templateUrl: './article-editor-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './article-editor-page.component.css'],
  standalone: true,
  imports: [CommonModule, NgOptimizedImage, ReactiveFormsModule, RouterModule, RichTextEditorComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ArticleEditorPageComponent {
  @ViewChild('previewDrawer') private previewDrawer?: ElementRef<HTMLElement>;

  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private routeInitializationStarted = false;
  private previewTrigger: HTMLElement | null = null;
  readonly state = inject(WorkspaceStateService);

  readonly editorMode: EditorMode = this.route.snapshot.data['editorMode'] === 'edit' ? 'edit' : 'create';
  readonly routeReady = signal(false);
  readonly activeTool = signal<EditorTool>('write');
  readonly coverImageMode = signal<'url' | 'upload'>('url');
  readonly uploadingCoverImage = signal(false);
  readonly coverImageFileName = signal('');
  readonly coverImageError = signal<string | null>(null);
  readonly coverImageLoadFailed = signal(false);
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
      if (this.state.loading() || this.routeInitializationStarted) {
        return;
      }

      this.routeInitializationStarted = true;
      void this.initializeFromRoute();
    });

    effect(() => {
      if (!this.routeReady()) {
        return;
      }

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
        this.coverImageLoadFailed.set(false);
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
      this.coverImageLoadFailed.set(false);
      this.feedback.clear();
    });
  }

  setActiveTool(tool: EditorTool, trigger?: EventTarget | null): void {
    const closingPreview = this.activeTool() === 'preview' && tool !== 'preview';
    if (tool === 'preview' && trigger instanceof HTMLElement) {
      this.previewTrigger = trigger;
    }
    this.activeTool.set(tool);

    if (tool === 'preview') {
      setTimeout(() => this.previewDrawer?.nativeElement.focus());
    } else if (closingPreview) {
      const restoreTarget = this.previewTrigger;
      this.previewTrigger = null;
      setTimeout(() => restoreTarget?.focus());
    }
  }

  handlePreviewKeydown(event: KeyboardEvent): void {
    if (event.key === 'Escape') {
      event.preventDefault();
      this.setActiveTool('write');
      return;
    }

    if (event.key !== 'Tab') {
      return;
    }

    const drawer = this.previewDrawer?.nativeElement;
    const focusable = drawer ? Array.from(drawer.querySelectorAll<HTMLElement>('button:not([disabled]), a[href], input:not([disabled]), [tabindex]:not([tabindex="-1"])')) : [];
    if (!drawer || focusable.length === 0) {
      event.preventDefault();
      drawer?.focus();
      return;
    }

    const first = focusable[0];
    const last = focusable[focusable.length - 1];
    if (event.shiftKey && document.activeElement === drawer) {
      event.preventDefault();
      last.focus();
      return;
    }
    if ((event.shiftKey && document.activeElement === first) || (!event.shiftKey && document.activeElement === last)) {
      event.preventDefault();
      (event.shiftKey ? last : first).focus();
    }
  }

  async publishArticle(): Promise<void> {
    this.articleForm.controls.status.setValue('published');
    await this.saveArticle();
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
      this.coverImageLoadFailed.set(false);
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
      await this.router.navigate(['/content/articles', article.id, 'edit'], { replaceUrl: true });
      this.feedback.success('Article saved successfully.');
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

  private async initializeFromRoute(): Promise<void> {
    this.state.clearError();

    if (this.editorMode === 'create') {
      this.state.clearSelectedArticle();
      this.routeReady.set(true);
      return;
    }

    const articleId = this.route.snapshot.paramMap.get('articleId')?.trim() ?? '';
    const routeArticle = this.state.articles().find((article) => article.id === articleId);
    if (!articleId || !routeArticle) {
      this.state.clearSelectedArticle();
      this.state.reportError('The requested article could not be found.');
      await this.router.navigate(['/content/articles']);
      return;
    }

    await this.state.selectArticle(articleId);
    if (this.state.selectedArticle()?.id !== articleId) {
      this.state.clearSelectedArticle();
      this.state.reportError('The requested article could not be opened.');
      await this.router.navigate(['/content/articles']);
      return;
    }

    this.routeReady.set(true);
  }

  private buildErrorMessage(message: string, error: unknown): string {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    return `${message}${detail}`;
  }
}
