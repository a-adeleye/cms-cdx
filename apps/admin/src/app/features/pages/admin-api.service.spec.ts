import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { AdminApiService } from './admin-api.service';
import { AuthTokenService } from './auth-token.service';

describe('AdminApiService article persistence', () => {
  let service: AdminApiService;
  let http: HttpTestingController;

  const article = {
    id: 'article-1',
    siteId: 'site-1',
    title: 'Existing article',
    slug: 'existing-article',
    excerpt: 'Existing article excerpt.',
    contentMarkdown: '# Existing article\n\nBody copy.',
    coverImageUrl: '',
    seoTitle: 'Existing article',
    seoDescription: 'Existing article description.',
    canonicalUrl: '',
    authorId: 'author-1',
    categoryId: 'category-1',
    tags: '',
    isFeatured: false,
    status: 'draft' as const,
    publishedAt: null,
    generatedByAi: false,
    humanReviewed: false,
    aiPrompt: '',
    aiModel: '',
    updatedAt: '2026-07-14T00:00:00.000Z',
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [
        AdminApiService,
        { provide: AuthTokenService, useValue: { getToken: () => 'test-token' } },
      ],
    });
    service = TestBed.inject(AdminApiService);
    http = TestBed.inject(HttpTestingController);
  });

  afterEach(() => http.verify());

  it('PATCHes an existing article when an id is present', async () => {
    const result = service.upsertArticle('site-1', article);
    const request = http.expectOne('/api/v1/articles/article-1');
    expect(request.request.method).toBe('PATCH');
    expect(request.request.body.id).toBe('article-1');
    request.flush(article);
    await expectAsync(result).toBeResolvedTo(article);
  });

  it('POSTs a new article without sending an id', async () => {
    const { id: _id, ...draft } = article;
    const result = service.upsertArticle('site-1', { ...draft, id: '' });
    const request = http.expectOne('/api/v1/sites/site-1/articles');
    expect(request.request.method).toBe('POST');
    expect(request.request.body.id).toBeUndefined();
    request.flush(article);
    await expectAsync(result).toBeResolvedTo(article);
  });

  it('sends destructive site and history requests to their scoped endpoints', async () => {
    const deleteSite = service.deleteSite('site-1');
    const siteRequest = http.expectOne('/api/v1/sites/site-1');
    expect(siteRequest.request.method).toBe('DELETE');
    siteRequest.flush(null);

    const clearHistory = service.clearBuildHistory('site-1');
    const historyRequest = http.expectOne('/api/v1/sites/site-1/builds');
    expect(historyRequest.request.method).toBe('DELETE');
    historyRequest.flush(null);

    await expectAsync(deleteSite).toBeResolved();
    await expectAsync(clearHistory).toBeResolved();
  });

  it('updates, replaces, and deletes scoped media assets', async () => {
    const updated = service.updateMediaAsset('site-1', 'media-1', 'Accessible image description');
    const updateRequest = http.expectOne('/api/v1/sites/site-1/media/media-1');
    expect(updateRequest.request.method).toBe('PATCH');
    expect(updateRequest.request.body).toEqual({ altText: 'Accessible image description' });
    updateRequest.flush({ id: 'media-1' });

    const replacement = service.replaceMediaFile('site-1', 'media-1', new File(['image'], 'replacement.png', { type: 'image/png' }), 'Replacement');
    const replacementRequest = http.expectOne('/api/v1/sites/site-1/media/media-1');
    expect(replacementRequest.request.method).toBe('PUT');
    expect(replacementRequest.request.body instanceof FormData).toBeTrue();
    replacementRequest.flush({ id: 'media-1' });

    const deleted = service.deleteMediaAsset('site-1', 'media-1');
    const deleteRequest = http.expectOne('/api/v1/sites/site-1/media/media-1');
    expect(deleteRequest.request.method).toBe('DELETE');
    deleteRequest.flush(null);

    await expectAsync(updated).toBeResolved();
    await expectAsync(replacement).toBeResolved();
    await expectAsync(deleted).toBeResolved();
  });
});
