import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router, RouterModule, convertToParamMap } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { ArticleEditorPageComponent } from './article-editor-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('ArticleEditorPageComponent', () => {
  let fixture: ComponentFixture<ArticleEditorPageComponent>;
  let router: Router;
  let resolveSaveArticle: ((value: { id: string }) => void) | null;

  const articleRecord = () => ({
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
    tags: '',
    updatedAt: '2026-05-23T00:00:00.000Z',
  });
  let selectedArticleValue: ReturnType<typeof articleRecord> | null;

  const fakeState = {
    loading: () => false,
    selectedSite: () => ({
      id: 'site-example',
      aiConfig: '{"provider":"google","model":"gemini-2.5-flash","apiKeySecretRef":"GEMINI_API_KEY","masterPrompt":"Write useful Anonime articles."}',
    }),
    selectedArticle: () => selectedArticleValue,
    articles: () => [articleRecord()],
    error: () => null,
    authors: () => [{ id: 'author-1', siteId: 'site-example', name: 'Author', slug: 'author', bio: '' }],
    categories: () => [{ id: 'category-1', siteId: 'site-example', name: 'Category', slug: 'category', description: '' }],
    reportError: jasmine.createSpy('reportError'),
    clearError: jasmine.createSpy('clearError'),
    clearSelectedArticle: jasmine.createSpy('clearSelectedArticle').and.callFake(() => {
      selectedArticleValue = null;
    }),
    saveArticle: jasmine.createSpy('saveArticle').and.callFake(
      () =>
        new Promise<{ id: string }>((resolve) => {
          resolveSaveArticle = resolve;
        }),
    ),
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
    generateAISuggestion: jasmine.createSpy('generateAISuggestion').and.resolveTo({
      suggestion: '## Suggested revision\n\nThis is the suggested article content.',
      model: 'gpt-4.1-mini',
    }),
    generateAIArticleDraft: jasmine.createSpy('generateAIArticleDraft').and.resolveTo({
      title: 'How email aliases protect privacy',
      slug: 'email-aliases-protect-privacy',
      category: 'Category',
      featured: true,
      tags: ['tag one'],
      seoTitle: 'How email aliases protect privacy',
      metaDescription: 'Learn how email aliases help protect your inbox privacy.',
      canonicalUrl: 'https://anonime.io/blog/email-aliases-protect-privacy',
      excerpt: 'A practical explanation of email aliases and privacy.',
      contentMarkdown: '# How email aliases protect privacy\n\nA complete AI draft.',
      coverImage: { fileName: 'ai-email-aliases.png', fileUrl: 'https://cdn.example/ai-email-aliases.png' },
      model: 'gemini-2.5-flash',
    }),
  };

  beforeEach(async () => {
    resolveSaveArticle = null;
    selectedArticleValue = articleRecord();
    fakeState.clearSelectedArticle.calls.reset();
    fakeState.selectArticle.calls.reset();
    fakeState.reportError.calls.reset();
    fakeState.selectArticle.and.callFake(async (articleId: string) => {
      selectedArticleValue = articleId === 'article-1' ? articleRecord() : null;
    });
    await TestBed.configureTestingModule({
      imports: [CommonModule, RouterModule, RouterTestingModule, ArticleEditorPageComponent],
      providers: [
        { provide: WorkspaceStateService, useValue: fakeState },
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              data: { editorMode: 'edit' },
              paramMap: convertToParamMap({ articleId: 'article-1' }),
            },
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(ArticleEditorPageComponent);
    router = TestBed.inject(Router);
    spyOn(router, 'navigate').and.resolveTo(true);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('renders a two-column editor workspace with the document and publishing tools', () => {
    const form = fixture.nativeElement.querySelector('.editor-workspace') as HTMLElement | null;
    expect(form).toBeTruthy();
    expect(fixture.nativeElement.querySelector('.editor-document')).toBeTruthy();
    expect(fixture.nativeElement.querySelector('.editor-content-field')).toBeTruthy();
    expect(fixture.nativeElement.querySelector('.editor-sidebar')).toBeTruthy();
  });

  it('keeps SEO fields and counters together in the sidebar', () => {
    const fields = fixture.nativeElement.querySelectorAll('.editor-seo-field');
    expect(fields.length).toBe(2);
    expect(fields[0].querySelector('#article-seo-title')).toBeTruthy();
    expect(fields[0].querySelector('#article-seo-title-help')).toBeTruthy();
    expect(fields[1].querySelector('#article-seo-description')).toBeTruthy();
    expect(fields[1].querySelector('#article-seo-description-help')).toBeTruthy();
  });

  it('keeps AI, preview, and history within the article workflow', () => {
    expect(fixture.nativeElement.textContent).toContain('Content');
    expect(fixture.nativeElement.textContent).toContain('Media');
    expect(fixture.nativeElement.textContent).toContain('Publishing');
    expect(fixture.nativeElement.textContent).toContain('SEO');

    const aiButton = Array.from(fixture.nativeElement.querySelectorAll('.ui-tabs button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.includes('AI Assistant'),
    );
    aiButton?.click();
    fixture.detectChanges();

    expect(fixture.nativeElement.textContent).toContain('AI writing assistance');
    expect(fixture.nativeElement.textContent).toContain('Configure AI provider');
    expect(fixture.nativeElement.querySelector('a[href="/configuration/ai"]')).toBeTruthy();
  });

  it('enables a configured AI provider and applies its generated suggestion', async () => {
    const aiButton = Array.from(fixture.nativeElement.querySelectorAll('.ui-tabs button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.includes('AI Assistant'),
    );
    aiButton?.click();
    fixture.detectChanges();

    const generateButton = Array.from(fixture.nativeElement.querySelectorAll('.editor-ai-panel > button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.includes('Generate suggestion'),
    ) ?? null;
    expect(generateButton?.disabled).toBeFalse();

    const prompt = fixture.nativeElement.querySelector('.editor-ai-panel textarea') as HTMLTextAreaElement;
    prompt.value = 'Make the introduction clearer';
    prompt.dispatchEvent(new Event('input'));
    generateButton?.click();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fakeState.generateAISuggestion).toHaveBeenCalledWith(jasmine.objectContaining({
      instruction: 'Make the introduction clearer',
      title: 'Example article',
    }));
    expect(fixture.nativeElement.textContent).toContain('Suggested Markdown');

    const useSuggestion = Array.from(fixture.nativeElement.querySelectorAll('.editor-ai-result button') as NodeListOf<HTMLButtonElement>)[0];
    useSuggestion?.click();
    expect(fixture.componentInstance.articleForm.controls.contentMarkdown.value).toContain('Suggested revision');
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
    expect(fixture.componentInstance.articleForm.controls.coverImageUrl.value).toBe('https://cdn.example/cover.jpg');
  });

  it('keeps the editor open and reports success after a successful save', async () => {
    fixture.componentInstance.articleForm.controls.contentMarkdown.setValue('# Example article\n\nBody text for the save test.');
    const saveButton = Array.from(fixture.nativeElement.querySelectorAll('.hero-actions button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.trim() === 'Save draft',
    ) ?? null;
    saveButton?.click();
    fixture.detectChanges();

    expect(saveButton?.disabled).toBeTrue();
    expect(fixture.nativeElement.textContent).toContain('Saving article...');

    resolveSaveArticle?.({ id: 'article-1' });

    await fixture.whenStable();

    expect(fakeState.saveArticle).toHaveBeenCalled();
    fixture.detectChanges();
    expect(fakeState.selectArticle).toHaveBeenCalledWith('article-1');
    expect(router.navigate).toHaveBeenCalledWith(['/content/articles', 'article-1', 'edit'], { replaceUrl: true });
    expect(fixture.nativeElement.textContent).toContain('Article saved successfully.');
  });

  it('renders tags as a free-text input', () => {
    const input = fixture.nativeElement.querySelector('input[formcontrolname="tags"]') as HTMLInputElement | null;
    expect(input).toBeTruthy();
    expect(input?.type).toBe('text');
  });

  it('shows validation feedback instead of silently doing nothing on invalid save', async () => {
    fakeState.saveArticle.calls.reset();
    fixture.componentInstance.articleForm.controls.title.setValue('');
    fixture.componentInstance.articleForm.controls.excerpt.setValue('');
    fixture.componentInstance.articleForm.controls.contentMarkdown.setValue('');
    fixture.componentInstance.articleForm.controls.seoTitle.setValue('');
    fixture.componentInstance.articleForm.controls.seoDescription.setValue('');
    fixture.componentInstance.articleForm.controls.authorId.setValue('');
    fixture.componentInstance.articleForm.controls.categoryId.setValue('');

    const saveButton = Array.from(fixture.nativeElement.querySelectorAll('.hero-actions button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.trim() === 'Save draft',
    ) ?? null;
    saveButton?.click();

    await fixture.whenStable();
    fixture.detectChanges();

    expect(fakeState.saveArticle).not.toHaveBeenCalled();
    expect(fixture.nativeElement.textContent).toContain('Please fix the highlighted fields before saving.');
    expect(fixture.nativeElement.textContent).toContain('Title is required.');
    expect(fixture.nativeElement.textContent).toContain('Excerpt is required.');
    expect(fixture.nativeElement.textContent).toContain('Article body is required.');
    expect(fixture.nativeElement.textContent).toContain('SEO title is required.');
    expect(fixture.nativeElement.textContent).toContain('Meta description is required.');
    expect(fixture.nativeElement.textContent).toContain('Author is required.');
    expect(fixture.nativeElement.textContent).toContain('Category is required.');
    expect(fixture.nativeElement.querySelector('#article-title')?.getAttribute('aria-invalid')).toBe('true');
  });

  it('clears a previously selected article on the explicit new-article route', async () => {
    fixture.destroy();
    selectedArticleValue = articleRecord();
    fakeState.clearSelectedArticle.calls.reset();
    const routeStub = TestBed.inject(ActivatedRoute) as unknown as {
      snapshot: { data: Record<string, string>; paramMap: ReturnType<typeof convertToParamMap> };
    };
    routeStub.snapshot.data = { editorMode: 'create' };
    routeStub.snapshot.paramMap = convertToParamMap({});

    const createFixture = TestBed.createComponent(ArticleEditorPageComponent);
    createFixture.detectChanges();
    await createFixture.whenStable();
    createFixture.detectChanges();

    expect(fakeState.clearSelectedArticle).toHaveBeenCalled();
    expect(createFixture.componentInstance.articleForm.controls.id.value).toBe('');
    expect(createFixture.componentInstance.articleForm.controls.title.value).toBe('');
    expect(createFixture.nativeElement.textContent).toContain('New article');
  });

  it('generates and applies a complete AI draft without requiring a topic', async () => {
    fixture.destroy();
    const routeStub = TestBed.inject(ActivatedRoute) as unknown as {
      snapshot: { data: Record<string, string>; paramMap: ReturnType<typeof convertToParamMap> };
    };
    routeStub.snapshot.data = { editorMode: 'create' };
    routeStub.snapshot.paramMap = convertToParamMap({});
    const createFixture = TestBed.createComponent(ArticleEditorPageComponent);
    createFixture.detectChanges();
    await createFixture.whenStable();
    createFixture.detectChanges();

    const writerButton = Array.from(createFixture.nativeElement.querySelectorAll('.hero-actions button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.trim() === 'AI Writer',
    );
    writerButton?.click();
    createFixture.detectChanges();
    const generateButton = Array.from(createFixture.nativeElement.querySelectorAll('.editor-ai-panel button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.includes('Generate full article'),
    );
    generateButton?.click();
    await createFixture.whenStable();
    createFixture.detectChanges();

    expect(fakeState.generateAIArticleDraft).toHaveBeenCalledWith({ topic: '' });
    expect(createFixture.componentInstance.articleForm.controls.title.value).toBe('How email aliases protect privacy');
    expect(createFixture.componentInstance.articleForm.controls.slug.value).toBe('email-aliases-protect-privacy');
    expect(createFixture.componentInstance.articleForm.controls.contentMarkdown.value).toContain('A complete AI draft.');
    expect(createFixture.componentInstance.articleForm.controls.coverImageUrl.value).toBe('https://cdn.example/ai-email-aliases.png');
    expect(createFixture.componentInstance.articleForm.controls.categoryId.value).toBe('category-1');
    expect(createFixture.componentInstance.articleForm.controls.tags.value).toBe('tag one');
  });

  it('rejects an edit route whose article id is not in the current workspace', async () => {
    fixture.destroy();
    fakeState.reportError.calls.reset();
    fakeState.clearSelectedArticle.calls.reset();
    (router.navigate as jasmine.Spy).calls.reset();
    const routeStub = TestBed.inject(ActivatedRoute) as unknown as {
      snapshot: { data: Record<string, string>; paramMap: ReturnType<typeof convertToParamMap> };
    };
    routeStub.snapshot.data = { editorMode: 'edit' };
    routeStub.snapshot.paramMap = convertToParamMap({ articleId: 'missing-article' });

    const invalidFixture = TestBed.createComponent(ArticleEditorPageComponent);
    invalidFixture.detectChanges();
    await invalidFixture.whenStable();

    expect(fakeState.clearSelectedArticle).toHaveBeenCalled();
    expect(fakeState.reportError).toHaveBeenCalledWith('The requested article could not be found.');
    expect(router.navigate).toHaveBeenCalledWith(['/content/articles']);
  });
});
