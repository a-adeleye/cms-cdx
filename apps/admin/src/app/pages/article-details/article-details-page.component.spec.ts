import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router, convertToParamMap } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { ArticleDetailsPageComponent } from './article-details-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('ArticleDetailsPageComponent', () => {
  let fixture: ComponentFixture<ArticleDetailsPageComponent>;
  let router: Router;

  const fakeState = {
    selectedArticle: () => ({
      id: 'article-1',
      siteId: 'site-example',
      authorId: 'author-1',
      categoryId: 'category-1',
      title: 'Saved article',
      slug: 'saved-article',
      excerpt: 'A short excerpt.',
      contentMarkdown: '# Saved article',
      coverImageUrl: '',
      status: 'draft' as const,
      isFeatured: false,
      publishedAt: null,
      seoTitle: 'Saved article',
      seoDescription: 'Saved article description',
      canonicalUrl: '',
      generatedByAi: false,
      humanReviewed: false,
      aiPrompt: '',
      aiModel: '',
      tagIds: ['tag-1'],
      updatedAt: '2026-06-02T00:00:00.000Z',
    }),
    error: () => null,
    authors: () => [{ id: 'author-1', siteId: 'site-example', name: 'Author', slug: 'author', bio: '' }],
    categories: () => [{ id: 'category-1', siteId: 'site-example', name: 'Category', slug: 'category', description: '' }],
    tags: () => [{ id: 'tag-1', siteId: 'site-example', name: 'Tag One', slug: 'tag-one' }],
    selectArticle: jasmine.createSpy('selectArticle').and.resolveTo(),
    deleteArticle: jasmine.createSpy('deleteArticle').and.resolveTo(),
    clearError: jasmine.createSpy('clearError'),
    reportError: jasmine.createSpy('reportError'),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CommonModule, RouterTestingModule, ArticleDetailsPageComponent],
      providers: [
        { provide: WorkspaceStateService, useValue: fakeState },
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              paramMap: convertToParamMap({ articleId: 'article-1' }),
            },
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(ArticleDetailsPageComponent);
    router = TestBed.inject(Router);
    spyOn(router, 'navigate').and.resolveTo(true);
    spyOn(window, 'confirm').and.returnValue(true);
    window.history.replaceState({ flashMessage: 'Article saved successfully.' }, '');

    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('shows the saved message and article details', () => {
    expect(fixture.nativeElement.textContent).toContain('Article saved successfully.');
    expect(fixture.nativeElement.textContent).toContain('Saved article');
    expect(fakeState.selectArticle).toHaveBeenCalledWith('article-1');
  });

  it('deletes the article and returns to the list with feedback', async () => {
    const deleteButton = Array.from(fixture.nativeElement.querySelectorAll('.hero-actions button') as NodeListOf<HTMLButtonElement>).find((button) =>
      button.textContent?.includes('Delete article'),
    );
    deleteButton?.click();

    await fixture.whenStable();

    expect(fakeState.deleteArticle).toHaveBeenCalledWith('article-1');
    expect(router.navigate).toHaveBeenCalledWith(['/articles'], {
      state: { flashMessage: 'Article deleted successfully.' },
    });
  });
});
