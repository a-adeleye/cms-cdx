import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { DashboardPageComponent } from './dashboard-page.component';

describe('DashboardPageComponent', () => {
  let fixture: ComponentFixture<DashboardPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CommonModule, RouterTestingModule, DashboardPageComponent],
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
      previewDeployProvider: 'cloudflare',
      previewDeployConfig: '{}',
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
    fixture.componentRef.setInput('recentBuilds', []);
    fixture.detectChanges();
  });

  it('renders the current site and recent article list', () => {
    expect(fixture.nativeElement.textContent).toContain('Drafts');
    expect(fixture.nativeElement.textContent).toContain('Example Site');
    expect(fixture.nativeElement.textContent).toContain('First article');
    expect(fixture.nativeElement.textContent).toContain('Deployment status');
    const newArticleLink = fixture.nativeElement.querySelector('.hero-actions a') as HTMLAnchorElement | null;
    expect(newArticleLink?.getAttribute('href')).toContain('/content/articles/new');
  });

  it('emits the article id when a recent article is clicked', () => {
    spyOn(fixture.componentInstance.articleSelected, 'emit');

    const articleButton = fixture.nativeElement.querySelector('.dashboard-article-title') as HTMLButtonElement | null;
    articleButton?.click();

    expect(fixture.componentInstance.articleSelected.emit).toHaveBeenCalledWith('article-1');
  });

  it('renders failed and empty deployment states without claiming they are healthy', () => {
    fixture.componentRef.setInput('recentBuilds', [
      {
        id: 'failed-production-build',
        siteId: 'site-example',
        status: 'failed',
        buildType: 'published',
        logs: 'Deployment failed.',
        outputPath: '',
        deployProvider: 'netlify',
        deployStatus: 'failed',
        deployUrl: '',
        startedAt: '2026-05-23T10:00:00.000Z',
        finishedAt: '2026-05-23T10:01:00.000Z',
      },
    ]);
    fixture.detectChanges();

    const environments = fixture.nativeElement.querySelectorAll('.deployment-environment') as NodeListOf<HTMLElement>;
    expect(environments[0].textContent).toContain('Failed');
    expect(environments[0].textContent).not.toContain('Healthy');
    expect(environments[1].textContent).toContain('Not deployed');
    expect(environments[1].textContent).not.toContain('Healthy');
    expect(fixture.nativeElement.textContent).toContain('Deployment failed');
    expect(fixture.nativeElement.textContent).not.toContain('Deployment succeeded');
    expect(fixture.nativeElement.textContent).not.toContain('Open Site');
  });
});
