import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { ArticlesPageComponent } from './articles-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('ArticlesPageComponent', () => {
  let fixture: ComponentFixture<ArticlesPageComponent>;
  let router: Router;

  const fakeState = {
    selectedSite: () => ({ name: 'Example Site', blogPath: '/blog' }),
    articles: () => [
      {
        id: 'article-1',
        siteId: 'site-example',
        authorId: 'author-1',
        categoryId: 'category-1',
        title: 'First article',
        slug: 'first-article',
        excerpt: 'A short excerpt.',
        contentMarkdown: '# First article',
        coverImageUrl: '',
        status: 'draft' as const,
        isFeatured: false,
        publishedAt: null,
        seoTitle: 'First article',
        seoDescription: 'First article description',
        canonicalUrl: '',
        generatedByAi: false,
        humanReviewed: false,
        aiPrompt: '',
        aiModel: '',
        tagIds: [],
        updatedAt: '2026-05-23T00:00:00.000Z',
      },
    ],
    categories: () => [{ id: 'category-1', siteId: 'site-example', name: 'Editorial', slug: 'editorial', description: '' }],
    authors: () => [{ id: 'author-1', siteId: 'site-example', name: 'Ada Author', slug: 'ada-author', bio: '' }],
    error: () => null,
    clearError: jasmine.createSpy('clearError'),
    reportError: jasmine.createSpy('reportError'),
    clearSelectedArticle: jasmine.createSpy('clearSelectedArticle'),
    selectArticle: jasmine.createSpy('selectArticle').and.resolveTo(),
    toggleFeatured: jasmine.createSpy('toggleFeatured').and.resolveTo(),
    setArticleStatus: jasmine.createSpy('setArticleStatus').and.resolveTo(),
    deleteArticle: jasmine.createSpy('deleteArticle').and.resolveTo(),
  };

  beforeEach(async () => {
    fakeState.clearSelectedArticle.calls.reset();
    fakeState.selectArticle.calls.reset();
    fakeState.deleteArticle.calls.reset();
    spyOn(globalThis, 'confirm').and.returnValue(true);
    await TestBed.configureTestingModule({
      imports: [CommonModule, RouterTestingModule, ArticlesPageComponent],
      providers: [{ provide: WorkspaceStateService, useValue: fakeState }],
    }).compileComponents();

    fixture = TestBed.createComponent(ArticlesPageComponent);
    router = TestBed.inject(Router);
    spyOn(router, 'navigate').and.resolveTo(true);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('renders articles in a single responsive table', () => {
    expect(fixture.nativeElement.querySelector('.articles-table-card')).toBeTruthy();
    expect(fixture.nativeElement.querySelector('.two-column')).toBeNull();
    expect(fixture.nativeElement.textContent).toContain('First article');
    expect(fixture.nativeElement.textContent).toContain('Ada Author');
    expect(fixture.nativeElement.textContent).toContain('Editorial');
  });

  it('navigates to the editor when creating a new article', async () => {
    fixture.componentInstance.startArticle();

    await fixture.whenStable();

    expect(fakeState.clearSelectedArticle).toHaveBeenCalled();
    expect(router.navigate).toHaveBeenCalledWith(['/content/articles/new']);
    expect(fakeState.selectArticle).not.toHaveBeenCalled();
  });

  it('opens an article from the table', async () => {
    const articleButton = fixture.nativeElement.querySelector('.article-title-cell') as HTMLButtonElement | null;
    articleButton?.click();

    await fixture.whenStable();
    fixture.detectChanges();

    expect(fakeState.selectArticle).not.toHaveBeenCalled();
    expect(router.navigate).toHaveBeenCalledWith(['/content/articles', 'article-1', 'edit']);
  });

  it('selects all visible articles with the header checkbox', () => {
    const headerCheckbox = fixture.nativeElement.querySelector('thead input[type="checkbox"]') as HTMLInputElement;

    headerCheckbox.click();
    fixture.detectChanges();

    expect(fixture.componentInstance.selectedArticleIds()).toEqual(['article-1']);

    headerCheckbox.click();
    fixture.detectChanges();

    expect(fixture.componentInstance.selectedArticleIds()).toEqual([]);
  });

  it('deletes an article from its row action', async () => {
    fixture.componentInstance.deleteArticle({ id: 'article-1', title: 'First article' });

    await fixture.whenStable();

    expect(fakeState.deleteArticle).toHaveBeenCalledWith('article-1');
  });

  it('only shows the clear button when a filter is active', () => {
    expect(fixture.nativeElement.querySelector('.articles-clear')).toBeNull();

    fixture.componentInstance.onSearch('foo');
    fixture.detectChanges();

    expect(fixture.nativeElement.querySelector('.articles-clear')).toBeTruthy();

    fixture.componentInstance.clearFilters();
    fixture.detectChanges();

    expect(fixture.nativeElement.querySelector('.articles-clear')).toBeNull();
    expect(fixture.componentInstance.articleSearch()).toBe('');
  });
});
