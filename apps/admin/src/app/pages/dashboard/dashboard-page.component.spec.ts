import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { DashboardPageComponent } from './dashboard-page.component';

describe('DashboardPageComponent', () => {
  let fixture: ComponentFixture<DashboardPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CommonModule, DashboardPageComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(DashboardPageComponent);
    fixture.componentRef.setInput('dashboardStats', [
      { label: 'Drafts', value: '3', detail: 'Ready for editing.' },
      { label: 'Published', value: '8', detail: 'Visible to site visitors.' },
    ]);
    fixture.componentRef.setInput('selectedSite', {
      id: 'site-example',
      name: 'Example Site',
      slug: 'example',
      domain: 'https://example.test',
      blogPath: '/articles',
      status: 'active',
      templateKey: 'default-blog',
      themeConfig: '{}',
      deployProvider: 'netlify',
      deployConfig: '{}',
      aiConfig: '{}',
      storageConfig: '{}',
      updatedAt: '2026-05-23T00:00:00.000Z',
    });
    fixture.componentRef.setInput('recentArticles', [
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
        status: 'draft',
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
    ]);
    fixture.detectChanges();
  });

  it('renders the selected site and recent article list', () => {
    expect(fixture.nativeElement.textContent).toContain('Drafts');
    expect(fixture.nativeElement.textContent).toContain('Example Site');
    expect(fixture.nativeElement.textContent).toContain('First article');
  });

  it('emits the article id when a recent article is clicked', () => {
    spyOn(fixture.componentInstance.articleSelected, 'emit');

    const articleButton = fixture.nativeElement.querySelector('.list-item') as HTMLButtonElement | null;
    articleButton?.click();

    expect(fixture.componentInstance.articleSelected.emit).toHaveBeenCalledWith('article-1');
  });
});
