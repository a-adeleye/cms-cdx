import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router, RouterModule } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { ArticleEditorPageComponent } from './article-editor-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('ArticleEditorPageComponent', () => {
  let fixture: ComponentFixture<ArticleEditorPageComponent>;
  let router: Router;

  const fakeState = {
    selectedArticle: () => ({
      id: 'article-1',
      siteId: 'site-example',
      authorId: 'author-1',
      categoryId: 'category-1',
      title: 'Example article',
      slug: 'example-article',
      excerpt: 'Example excerpt for the editor test.',
      contentMarkdown: '# Example article',
      coverImageUrl: '',
      status: 'draft' as const,
      isFeatured: false,
      publishedAt: null,
      seoTitle: 'Example article',
      seoDescription: 'Example description',
      canonicalUrl: '',
      generatedByAi: false,
      humanReviewed: false,
      aiPrompt: '',
      aiModel: '',
      tagIds: [],
      updatedAt: '2026-05-23T00:00:00.000Z',
    }),
    error: () => null,
    authors: () => [{ id: 'author-1', siteId: 'site-example', name: 'Author', slug: 'author', bio: '' }],
    categories: () => [{ id: 'category-1', siteId: 'site-example', name: 'Category', slug: 'category', description: '' }],
    tags: () => [
      { id: 'tag-1', siteId: 'site-example', name: 'Tag One', slug: 'tag-one' },
      { id: 'tag-2', siteId: 'site-example', name: 'Tag Two', slug: 'tag-two' },
    ],
    reportError: jasmine.createSpy('reportError'),
    clearError: jasmine.createSpy('clearError'),
    saveArticle: jasmine.createSpy('saveArticle').and.resolveTo({ id: 'article-1' }),
    selectArticle: jasmine.createSpy('selectArticle').and.resolveTo(),
    uploadMediaFile: jasmine.createSpy('uploadMediaFile').and.resolveTo({
      id: 'media-1',
      siteId: 'site-example',
      fileName: 'cover.jpg',
      fileUrl: 'https://cdn.example/cover.jpg',
      mimeType: 'image/jpeg',
      sizeBytes: 1024,
      storageProvider: 'minio',
      storageKey: 'site-example/media/cover.jpg',
      altText: 'Example article',
    }),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CommonModule, RouterModule, RouterTestingModule, ArticleEditorPageComponent],
      providers: [{ provide: WorkspaceStateService, useValue: fakeState }],
    }).compileComponents();

    fixture = TestBed.createComponent(ArticleEditorPageComponent);
    router = TestBed.inject(Router);
    spyOn(router, 'navigate').and.resolveTo(true);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('renders a two-column form with the markdown body spanning both columns', () => {
    const form = fixture.nativeElement.querySelector('.article-form-grid') as HTMLElement | null;
    expect(form).toBeTruthy();
    expect(fixture.nativeElement.querySelector('.article-layout')).toBeTruthy();
    expect(fixture.nativeElement.querySelector('.article-markdown.article-span-2')).toBeTruthy();
    expect(fixture.nativeElement.querySelector('.article-layout .two-column')).toBeNull();
  });

  it('uploads a cover image file and populates the cover image url', async () => {
    const file = new File(['fake-image'], 'cover.jpg', { type: 'image/jpeg' });
    const input = fixture.nativeElement.querySelector('input[type="file"]') as HTMLInputElement | null;
    Object.defineProperty(input, 'files', {
      value: [file],
    });

    input?.dispatchEvent(new Event('change'));
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fakeState.uploadMediaFile).toHaveBeenCalled();
    expect(fixture.nativeElement.querySelector('input[formcontrolname="coverImageUrl"]')?.value).toBe('https://cdn.example/cover.jpg');
  });

  it('navigates to the article details page after a successful save', async () => {
    fixture.componentInstance.articleForm.controls.contentMarkdown.setValue('# Example article\n\nBody text for the save test.');
    const saveButton = fixture.nativeElement.querySelector('.hero-actions .button-primary') as HTMLButtonElement | null;
    saveButton?.click();

    await fixture.whenStable();

    expect(fakeState.saveArticle).toHaveBeenCalled();
    expect(router.navigate).toHaveBeenCalledWith(['/articles', 'article-1'], {
      state: { flashMessage: 'Article saved successfully.' },
    });
  });

  it('renders tags as a multi-select from the workspace tag list', () => {
    const select = fixture.nativeElement.querySelector('select[formcontrolname="tagIds"]') as HTMLSelectElement | null;
    expect(select).toBeTruthy();
    expect(select?.multiple).toBeTrue();
    expect(select?.options.length).toBe(2);
    expect(select?.options[0].textContent?.trim()).toBe('Tag One');
    expect(select?.options[1].textContent?.trim()).toBe('Tag Two');
  });

  it('shows validation feedback instead of silently doing nothing on invalid save', async () => {
    fakeState.saveArticle.calls.reset();
    fixture.componentInstance.articleForm.controls.title.setValue('');
    fixture.componentInstance.articleForm.controls.contentMarkdown.setValue('');

    const saveButton = fixture.nativeElement.querySelector('.hero-actions .button-primary') as HTMLButtonElement | null;
    saveButton?.click();

    await fixture.whenStable();
    fixture.detectChanges();

    expect(fakeState.saveArticle).not.toHaveBeenCalled();
    expect(fixture.nativeElement.textContent).toContain('Please fix the highlighted fields before saving.');
    expect(fixture.nativeElement.textContent).toContain('Title is required.');
    expect(fixture.nativeElement.textContent).toContain('Markdown body is required.');
  });
});
