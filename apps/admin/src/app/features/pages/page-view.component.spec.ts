import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { of } from 'rxjs';
import { DashboardPageComponent } from '../../pages/dashboard/dashboard-page.component';
import { PageViewComponent } from './page-view.component';
import { SettingsPageComponent } from '../../pages/settings/settings-page.component';
import { WorkspaceStateService } from './workspace-state.service';

const settingsPage = {
  path: 'settings',
  navLabel: 'Settings',
  kind: 'settings' as const,
  eyebrow: 'Administration',
  title: 'Settings',
  primaryAction: { label: 'Manage sites', path: '/sites' },
};

const articlesPage = {
  path: 'articles',
  navLabel: 'Articles',
  kind: 'articles' as const,
  eyebrow: 'Editorial',
  title: 'Articles',
  primaryAction: { label: 'Open editor', path: '/article-editor' },
};

const articleEditorPage = {
  path: 'article-editor',
  navLabel: 'Article Editor',
  kind: 'article-editor' as const,
  eyebrow: 'Writing',
  title: 'Article Editor',
  primaryAction: { label: 'Save draft', path: '/articles' },
};

describe('PageViewComponent', () => {
  let fixture: ComponentFixture<PageViewComponent>;

  const loginPage = {
    path: 'login',
    navLabel: 'Login',
    kind: 'login' as const,
    eyebrow: 'Access',
    title: 'Login',
    primaryAction: { label: 'Open dashboard', path: '/dashboard' },
  };

  const fakeState = {
    loading: () => false,
    isAuthenticated: () => false,
    error: () => null,
    selectedSite: () => null,
    selectedSiteId: () => 'site-example',
    selectedArticle: () => null,
    sites: () => [],
    authSession: () => null,
    authors: () => [],
    categories: () => [],
    articles: () => [],
    tags: () => [],
    mediaAssets: () => [],
    builds: () => [],
    dashboardStats: () => [],
    landingSections: () => [],
    reportError: jasmine.createSpy('reportError'),
    login: jasmine.createSpy('login').and.resolveTo(),
    logout: jasmine.createSpy('logout').and.resolveTo(),
    selectSite: jasmine.createSpy('selectSite').and.resolveTo(),
    selectArticle: jasmine.createSpy('selectArticle').and.resolveTo(),
    clearSelectedArticle: jasmine.createSpy('clearSelectedArticle'),
    createSite: jasmine.createSpy('createSite').and.resolveTo(),
    updateSelectedSite: jasmine.createSpy('updateSelectedSite').and.resolveTo(),
    createArticleDraft: jasmine.createSpy('createArticleDraft').and.resolveTo({ id: 'article-1' }),
    saveArticle: jasmine.createSpy('saveArticle').and.resolveTo({ id: 'article-1' }),
    triggerBuild: jasmine.createSpy('triggerBuild').and.resolveTo(),
    toggleLandingSection: jasmine.createSpy('toggleLandingSection').and.resolveTo(),
    moveLandingSection: jasmine.createSpy('moveLandingSection').and.resolveTo(),
    uploadMedia: jasmine.createSpy('uploadMedia').and.resolveTo(),
    toggleFeatured: jasmine.createSpy('toggleFeatured').and.resolveTo(),
    setArticleStatus: jasmine.createSpy('setArticleStatus').and.resolveTo(),
    updateSelectedSiteSettings: jasmine.createSpy('updateSelectedSiteSettings').and.resolveTo(),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [PageViewComponent],
      imports: [CommonModule, ReactiveFormsModule, RouterTestingModule, DashboardPageComponent, SettingsPageComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({ page: loginPage }),
          },
        },
        { provide: WorkspaceStateService, useValue: fakeState },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(PageViewComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('renders only the sign-in form on the login page', () => {
    expect(fixture.nativeElement.querySelector('.hero')).toBeNull();
    expect(fixture.nativeElement.querySelector('.two-column')).toBeNull();
    expect(fixture.nativeElement.querySelectorAll('.panel').length).toBe(1);
    expect(fixture.nativeElement.querySelector('form')).toBeTruthy();
  });
});

describe('Settings page', () => {
  let fixture: ComponentFixture<PageViewComponent>;

  const fakeState = {
    loading: () => false,
    isAuthenticated: () => true,
    error: () => null,
    selectedSite: () => ({
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
      aiConfig: '{}',
      storageConfig: '{}',
      updatedAt: '2026-05-23T00:00:00.000Z',
    }),
    selectedSiteId: () => 'site-example',
    selectedArticle: () => null,
    sites: () => [
      {
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
        aiConfig: '{}',
        storageConfig: '{}',
        updatedAt: '2026-05-23T00:00:00.000Z',
      },
    ],
    authSession: () => null,
    authors: () => [],
    categories: () => [],
    articles: () => [],
    tags: () => [],
    mediaAssets: () => [],
    builds: () => [],
    dashboardStats: () => [],
    landingSections: () => [],
    reportError: jasmine.createSpy('reportError'),
    login: jasmine.createSpy('login').and.resolveTo(),
    logout: jasmine.createSpy('logout').and.resolveTo(),
    selectSite: jasmine.createSpy('selectSite').and.resolveTo(),
    selectArticle: jasmine.createSpy('selectArticle').and.resolveTo(),
    clearSelectedArticle: jasmine.createSpy('clearSelectedArticle'),
    createSite: jasmine.createSpy('createSite').and.resolveTo(),
    updateSelectedSite: jasmine.createSpy('updateSelectedSite').and.resolveTo(),
    createArticleDraft: jasmine.createSpy('createArticleDraft').and.resolveTo({ id: 'article-1' }),
    saveArticle: jasmine.createSpy('saveArticle').and.resolveTo({ id: 'article-1' }),
    triggerBuild: jasmine.createSpy('triggerBuild').and.resolveTo(),
    toggleLandingSection: jasmine.createSpy('toggleLandingSection').and.resolveTo(),
    moveLandingSection: jasmine.createSpy('moveLandingSection').and.resolveTo(),
    uploadMedia: jasmine.createSpy('uploadMedia').and.resolveTo(),
    toggleFeatured: jasmine.createSpy('toggleFeatured').and.resolveTo(),
    setArticleStatus: jasmine.createSpy('setArticleStatus').and.resolveTo(),
    updateSelectedSiteSettings: jasmine.createSpy('updateSelectedSiteSettings').and.resolveTo(),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [PageViewComponent],
      imports: [CommonModule, ReactiveFormsModule, RouterTestingModule, DashboardPageComponent, SettingsPageComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({ page: settingsPage }),
          },
        },
        { provide: WorkspaceStateService, useValue: fakeState },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(PageViewComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('shows settings shortcuts for the administration areas', () => {
    expect(fixture.nativeElement.querySelector('h2')?.textContent).not.toContain('Site context');
    expect(fixture.nativeElement.textContent).toContain('Shortcut links');
    expect(fixture.nativeElement.textContent).toContain('Sites');
    expect(fixture.nativeElement.textContent).toContain('Deployment settings');
    expect(fixture.nativeElement.querySelectorAll('.settings-link').length).toBeGreaterThan(0);
  });
});

describe('Articles page', () => {
  let fixture: ComponentFixture<PageViewComponent>;
  let router: Router;

  const fakeState = {
    loading: () => false,
    isAuthenticated: () => true,
    error: () => null,
    selectedSite: () => ({
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
      aiConfig: '{}',
      storageConfig: '{}',
      updatedAt: '2026-05-23T00:00:00.000Z',
    }),
    selectedSiteId: () => 'site-example',
    selectedArticle: () => null,
    sites: () => [],
    authSession: () => null,
    authors: () => [],
    categories: () => [],
    articles: () => [],
    tags: () => [],
    mediaAssets: () => [],
    builds: () => [],
    dashboardStats: () => [],
    landingSections: () => [],
    reportError: jasmine.createSpy('reportError'),
    login: jasmine.createSpy('login').and.resolveTo(),
    logout: jasmine.createSpy('logout').and.resolveTo(),
    selectSite: jasmine.createSpy('selectSite').and.resolveTo(),
    selectArticle: jasmine.createSpy('selectArticle').and.resolveTo(),
    clearSelectedArticle: jasmine.createSpy('clearSelectedArticle'),
    createSite: jasmine.createSpy('createSite').and.resolveTo(),
    updateSelectedSite: jasmine.createSpy('updateSelectedSite').and.resolveTo(),
    createArticleDraft: jasmine.createSpy('createArticleDraft').and.resolveTo({ id: 'article-1' }),
    saveArticle: jasmine.createSpy('saveArticle').and.resolveTo({ id: 'article-1' }),
    triggerBuild: jasmine.createSpy('triggerBuild').and.resolveTo(),
    toggleLandingSection: jasmine.createSpy('toggleLandingSection').and.resolveTo(),
    moveLandingSection: jasmine.createSpy('moveLandingSection').and.resolveTo(),
    uploadMedia: jasmine.createSpy('uploadMedia').and.resolveTo(),
    toggleFeatured: jasmine.createSpy('toggleFeatured').and.resolveTo(),
    setArticleStatus: jasmine.createSpy('setArticleStatus').and.resolveTo(),
    updateSelectedSiteSettings: jasmine.createSpy('updateSelectedSiteSettings').and.resolveTo(),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [PageViewComponent],
      imports: [CommonModule, ReactiveFormsModule, RouterTestingModule.withRoutes([
        { path: 'articles', component: PageViewComponent },
        { path: 'article-editor', component: PageViewComponent },
      ]), DashboardPageComponent, SettingsPageComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({ page: articlesPage }),
          },
        },
        { provide: WorkspaceStateService, useValue: fakeState },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(PageViewComponent);
    router = TestBed.inject(Router);
    spyOn(router, 'navigate').and.resolveTo(true);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('navigates to the article editor when creating a new article', async () => {
    expect(fixture.nativeElement.textContent).toContain('All');
    expect(fixture.nativeElement.textContent).toContain('Draft');
    expect(fixture.nativeElement.textContent).toContain('Published');
    expect(fixture.nativeElement.textContent).not.toContain('Open editor');

    const newArticleButton = fixture.nativeElement.querySelector('.toolbar .button-primary') as HTMLButtonElement | null;

    expect(newArticleButton).toBeTruthy();

    newArticleButton?.click();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fakeState.clearSelectedArticle).toHaveBeenCalled();
    expect(router.navigate).toHaveBeenCalledWith(['/article-editor']);
  });
});

describe('Article editor page', () => {
  let fixture: ComponentFixture<PageViewComponent>;

  const fakeState = {
    loading: () => false,
    isAuthenticated: () => true,
    error: () => 'Unable to save article. validation error: invalid tag id "not-a-uuid"',
    selectedSite: () => ({
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
      aiConfig: '{}',
      storageConfig: '{}',
      updatedAt: '2026-05-23T00:00:00.000Z',
    }),
    selectedSiteId: () => 'site-example',
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
    sites: () => [],
    authSession: () => null,
    authors: () => [{ id: 'author-1', siteId: 'site-example', name: 'Author', slug: 'author', bio: '' }],
    categories: () => [{ id: 'category-1', siteId: 'site-example', name: 'Category', slug: 'category', description: '' }],
    articles: () => [],
    tags: () => [],
    mediaAssets: () => [],
    builds: () => [],
    dashboardStats: () => [],
    landingSections: () => [],
    reportError: jasmine.createSpy('reportError'),
    login: jasmine.createSpy('login').and.resolveTo(),
    logout: jasmine.createSpy('logout').and.resolveTo(),
    selectSite: jasmine.createSpy('selectSite').and.resolveTo(),
    selectArticle: jasmine.createSpy('selectArticle').and.resolveTo(),
    createSite: jasmine.createSpy('createSite').and.resolveTo(),
    updateSelectedSite: jasmine.createSpy('updateSelectedSite').and.resolveTo(),
    createArticleDraft: jasmine.createSpy('createArticleDraft').and.resolveTo({ id: 'article-1' }),
    saveArticle: jasmine.createSpy('saveArticle').and.resolveTo({ id: 'article-1' }),
    triggerBuild: jasmine.createSpy('triggerBuild').and.resolveTo(),
    toggleLandingSection: jasmine.createSpy('toggleLandingSection').and.resolveTo(),
    moveLandingSection: jasmine.createSpy('moveLandingSection').and.resolveTo(),
    uploadMedia: jasmine.createSpy('uploadMedia').and.resolveTo(),
    toggleFeatured: jasmine.createSpy('toggleFeatured').and.resolveTo(),
    setArticleStatus: jasmine.createSpy('setArticleStatus').and.resolveTo(),
    updateSelectedSiteSettings: jasmine.createSpy('updateSelectedSiteSettings').and.resolveTo(),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [PageViewComponent],
      imports: [CommonModule, ReactiveFormsModule, RouterTestingModule, DashboardPageComponent, SettingsPageComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({ page: articleEditorPage }),
          },
        },
        { provide: WorkspaceStateService, useValue: fakeState },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(PageViewComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('expands the editor form when the preview is hidden', async () => {
    const layout = fixture.nativeElement.querySelector('.article-layout') as HTMLElement | null;
    const previewToggle = fixture.nativeElement.querySelector('.preview-toggle input[type="checkbox"]') as HTMLInputElement | null;
    const buttonLabels = Array.from(fixture.nativeElement.querySelectorAll('button') as NodeListOf<HTMLButtonElement>)
      .map((button) => button.textContent?.trim())
      .filter((label): label is string => Boolean(label));

    expect(layout).toBeTruthy();
    expect(fixture.nativeElement.querySelector('.preview-drawer')).toBeNull();
    expect(previewToggle?.checked).toBeFalse();
    expect(buttonLabels).toEqual(['Save']);
    expect(fixture.nativeElement.querySelector('input[formcontrolname="slug"]')).toBeNull();
    expect(fixture.nativeElement.querySelector('select[formcontrolname="status"]')).toBeNull();
    expect(fixture.nativeElement.querySelector('input[formcontrolname="id"]')?.value).toBe('article-1');
    expect(fixture.nativeElement.querySelector('.field-row')).toBeTruthy();
    expect(fixture.nativeElement.querySelector('.form-grid textarea[formcontrolname="contentMarkdown"]')).toBeNull();
    expect(fixture.nativeElement.querySelector('.markdown-panel textarea[formcontrolname="contentMarkdown"]')).toBeTruthy();
    expect(fixture.nativeElement.querySelector('.page-error')?.textContent).toContain('invalid tag id');

    previewToggle?.click();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fixture.nativeElement.querySelector('.preview-drawer')).toBeTruthy();
  });
});

describe('Blank article editor page', () => {
  let fixture: ComponentFixture<PageViewComponent>;

  const fakeState = {
    loading: () => false,
    isAuthenticated: () => true,
    error: () => null,
    selectedSite: () => ({
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
      aiConfig: '{}',
      storageConfig: '{}',
      updatedAt: '2026-05-23T00:00:00.000Z',
    }),
    selectedSiteId: () => 'site-example',
    selectedArticle: () => null,
    sites: () => [],
    authSession: () => null,
    authors: () => [{ id: 'author-1', siteId: 'site-example', name: 'Author', slug: 'author', bio: '' }],
    categories: () => [{ id: 'category-1', siteId: 'site-example', name: 'Category', slug: 'category', description: '' }],
    articles: () => [],
    tags: () => [],
    mediaAssets: () => [],
    builds: () => [],
    dashboardStats: () => [],
    landingSections: () => [],
    reportError: jasmine.createSpy('reportError'),
    login: jasmine.createSpy('login').and.resolveTo(),
    logout: jasmine.createSpy('logout').and.resolveTo(),
    selectSite: jasmine.createSpy('selectSite').and.resolveTo(),
    selectArticle: jasmine.createSpy('selectArticle').and.resolveTo(),
    clearSelectedArticle: jasmine.createSpy('clearSelectedArticle'),
    createSite: jasmine.createSpy('createSite').and.resolveTo(),
    updateSelectedSite: jasmine.createSpy('updateSelectedSite').and.resolveTo(),
    createArticleDraft: jasmine.createSpy('createArticleDraft').and.resolveTo({ id: 'article-1' }),
    saveArticle: jasmine.createSpy('saveArticle').and.resolveTo({ id: 'article-1' }),
    triggerBuild: jasmine.createSpy('triggerBuild').and.resolveTo(),
    toggleLandingSection: jasmine.createSpy('toggleLandingSection').and.resolveTo(),
    moveLandingSection: jasmine.createSpy('moveLandingSection').and.resolveTo(),
    uploadMedia: jasmine.createSpy('uploadMedia').and.resolveTo(),
    toggleFeatured: jasmine.createSpy('toggleFeatured').and.resolveTo(),
    setArticleStatus: jasmine.createSpy('setArticleStatus').and.resolveTo(),
    updateSelectedSiteSettings: jasmine.createSpy('updateSelectedSiteSettings').and.resolveTo(),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [PageViewComponent],
      imports: [CommonModule, ReactiveFormsModule, RouterTestingModule, DashboardPageComponent, SettingsPageComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({ page: articleEditorPage }),
          },
        },
        { provide: WorkspaceStateService, useValue: fakeState },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(PageViewComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('starts blank when no article is selected', () => {
    expect(fixture.nativeElement.querySelector('input[formcontrolname="id"]')?.value).toBe('');
    expect(fixture.nativeElement.querySelector('.preview-drawer')).toBeNull();
  });
});
