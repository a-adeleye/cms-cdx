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
    error: () => null,
    clearError: jasmine.createSpy('clearError'),
    reportError: jasmine.createSpy('reportError'),
    clearSelectedArticle: jasmine.createSpy('clearSelectedArticle'),
    selectArticle: jasmine.createSpy('selectArticle').and.resolveTo(),
    toggleFeatured: jasmine.createSpy('toggleFeatured').and.resolveTo(),
    setArticleStatus: jasmine.createSpy('setArticleStatus').and.resolveTo(),
  };

  beforeEach(async () => {
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

  it('renders articles in a single column list', () => {
    expect(fixture.nativeElement.querySelector('.panel')).toBeTruthy();
    expect(fixture.nativeElement.querySelector('.two-column')).toBeNull();
    expect(fixture.nativeElement.textContent).toContain('First article');
  });

  it('navigates to the editor when creating a new article', async () => {
    const button = fixture.nativeElement.querySelector('.hero-actions .button') as HTMLButtonElement | null;
    button?.click();

    await fixture.whenStable();

    expect(fakeState.clearSelectedArticle).toHaveBeenCalled();
    expect(router.navigate).toHaveBeenCalledWith(['/articles/editor']);
  });

  it('shows success feedback after updating an article', async () => {
    const featureButton = fixture.nativeElement.querySelector('.list-card .button-secondary:nth-child(2)') as HTMLButtonElement | null;
    featureButton?.click();

    await fixture.whenStable();
    fixture.detectChanges();

    expect(fakeState.toggleFeatured).toHaveBeenCalledWith('article-1');
    expect(fixture.nativeElement.textContent).toContain('Featured state updated successfully.');
  });
});
