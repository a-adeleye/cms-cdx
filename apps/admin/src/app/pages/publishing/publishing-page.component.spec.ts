import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { PublishingPageComponent } from './publishing-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('PublishingPageComponent', () => {
  let fixture: ComponentFixture<PublishingPageComponent>;
  let fakeState: WorkspaceStateService;

  beforeEach(async () => {
    const selectedSite = {
      id: 'site-example',
      name: 'Example Site',
      slug: 'example',
      domain: 'https://example.test',
      blogPath: '/articles',
      status: 'active' as const,
      templateKey: 'default-blog',
      themeConfig: '{}',
      deployProvider: 'netlify',
      deployConfig: '{}',
      previewDeployProvider: 'cloudflare',
      previewDeployConfig: '{}',
      aiConfig: '{}',
      storageConfig: '{}',
      updatedAt: '2026-05-23T00:00:00.000Z',
    };

    fakeState = {
      selectedSite: jasmine.createSpy('selectedSite').and.returnValue(selectedSite),
      articles: jasmine.createSpy('articles').and.returnValue([
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
      ]),
      error: jasmine.createSpy('error').and.returnValue(null),
      reportError: jasmine.createSpy('reportError'),
      triggerBuild: jasmine.createSpy('triggerBuild').and.resolveTo({
        id: 'build-1',
        siteId: 'site-example',
        status: 'success',
        buildType: 'preview',
        logs: '',
        outputPath: 'dist/preview/site',
        deployProvider: 'cloudflare',
        deployStatus: 'deployed',
        deployUrl: '',
        startedAt: null,
        finishedAt: null,
      }),
    } as unknown as WorkspaceStateService;

    await TestBed.configureTestingModule({
      imports: [CommonModule, RouterTestingModule, PublishingPageComponent],
      providers: [{ provide: WorkspaceStateService, useValue: fakeState }],
    }).compileComponents();

    fixture = TestBed.createComponent(PublishingPageComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('renders the selected article list and build targets', () => {
    expect(fixture.nativeElement.textContent).toContain('Publishing');
    expect(fixture.nativeElement.textContent).toContain('First article');
    expect(fixture.nativeElement.textContent).toContain('cloudflare');
    expect(fixture.nativeElement.textContent).toContain('default-blog');
  });

  it('selects all articles and triggers preview builds', async () => {
    const selectAll = fixture.nativeElement.querySelector('.publishing-select-all input') as HTMLInputElement | null;
    selectAll?.click();
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fixture.nativeElement.textContent).toContain('1 selected of 1');

    const previewButton = Array.from(fixture.nativeElement.querySelectorAll('button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.trim() === 'Trigger preview build',
    );
    previewButton?.click();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fakeState.triggerBuild).toHaveBeenCalledWith('preview', ['article-1']);
    expect(fixture.nativeElement.textContent).toContain('Preview build started successfully.');
  });
});
