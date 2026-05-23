import { TestBed } from '@angular/core/testing';
import { AdminApiService } from './admin-api.service';
import { AuthTokenService } from './auth-token.service';
import { WorkspaceStateService } from './workspace-state.service';

describe('WorkspaceStateService', () => {
  let service: WorkspaceStateService;
  let tokenStore: jasmine.SpyObj<AuthTokenService>;
  let api: jasmine.SpyObj<AdminApiService>;

  beforeEach(async () => {
    tokenStore = jasmine.createSpyObj<AuthTokenService>('AuthTokenService', ['getToken', 'setToken', 'clear']);
    tokenStore.getToken.and.returnValue(null);

    api = jasmine.createSpyObj<AdminApiService>('AdminApiService', [
      'login',
      'logout',
      'me',
      'loadWorkspace',
      'createSite',
      'updateSite',
      'upsertArticle',
      'updateLandingSection',
      'reorderLandingSections',
      'createBuild',
      'createMediaAsset',
    ]);

    api.loadWorkspace.and.callFake(async (siteId?: string) => ({
      user: {
        id: 'user-1',
        email: 'admin@example.com',
        fullName: 'Admin User',
        role: 'admin',
      },
      selectedSiteId: siteId ?? 'site-example',
      selectedArticleId: siteId === 'site-new' ? null : 'article-1',
      sites: [
        {
          id: siteId === 'site-new' ? 'site-new' : 'site-example',
          name: siteId === 'site-new' ? 'New Site' : 'Example Site',
          slug: siteId === 'site-new' ? 'new-site' : 'example',
          domain: siteId === 'site-new' ? 'https://new.example' : 'https://example.test',
          blogPath: '/articles',
          status: 'active',
          templateKey: 'default-blog',
          themeConfig: '{}',
          deployProvider: '',
          deployConfig: '{}',
          aiConfig: '{}',
          storageConfig: '{}',
          updatedAt: '2026-05-23T00:00:00.000Z',
        },
      ],
      landingSections: [],
      articles: [],
      authors: [],
      categories: [],
      tags: [],
      mediaAssets: [],
      builds: [],
    }));

    api.login.and.resolveTo({
      token: 'jwt-token',
      user: {
        id: 'user-1',
        email: 'admin@example.com',
        fullName: 'Admin User',
        role: 'admin',
      },
    });

    api.me.and.resolveTo({
      user: {
        id: 'user-1',
        email: 'admin@example.com',
        fullName: 'Admin User',
        role: 'admin',
      },
    });

    api.createSite.and.resolveTo({
      id: 'site-new',
      name: 'New Site',
      slug: 'new-site',
      domain: 'https://new.example',
      blogPath: '/articles',
      status: 'active',
      templateKey: 'default-blog',
      themeConfig: '{}',
      deployProvider: '',
      deployConfig: '{}',
      aiConfig: '{}',
      storageConfig: '{}',
      updatedAt: '2026-05-23T00:00:00.000Z',
    });

    api.upsertArticle.and.resolveTo({
      id: 'article-1',
      siteId: 'site-example',
      authorId: 'author-1',
      categoryId: 'category-1',
      title: 'Saved article',
      slug: 'saved-article',
      excerpt: 'A short test article.',
      contentMarkdown: '# Saved article\n\nBody copy for the test article.',
      coverImageUrl: '',
      status: 'draft',
      isFeatured: false,
      publishedAt: null,
      seoTitle: 'Saved article',
      seoDescription: 'A test article for the CMS builder.',
      canonicalUrl: 'https://example.test/articles/saved-article/',
      generatedByAi: false,
      humanReviewed: false,
      aiPrompt: '',
      aiModel: '',
      tagIds: [],
      updatedAt: '2026-05-23T00:00:00.000Z',
    });

    await TestBed.configureTestingModule({
      providers: [
        WorkspaceStateService,
        { provide: AuthTokenService, useValue: tokenStore },
        { provide: AdminApiService, useValue: api },
      ],
    }).compileComponents();

    service = TestBed.inject(WorkspaceStateService);
    await Promise.resolve();
  });

  it('logs in with JWT and hydrates the workspace from the API', async () => {
    await service.login('admin@example.com', 'admin123');

    expect(tokenStore.setToken).toHaveBeenCalledWith('jwt-token');
    expect(service.isAuthenticated()).toBeTrue();
    expect(service.selectedSite().id).toBe('site-example');
    expect(service.authSession()?.email).toBe('admin@example.com');
  });

  it('creates a site and refreshes the selected site from the API', async () => {
    await service.login('admin@example.com', 'admin123');

    await service.createSite({
      name: 'New Site',
      slug: 'new-site',
      domain: 'https://new.example',
      blogPath: '/articles',
      templateKey: 'default-blog',
    });

    expect(api.createSite).toHaveBeenCalled();
    expect(service.selectedSite().id).toBe('site-new');
    expect(service.sites().some((site) => site.slug === 'new-site')).toBeTrue();
  });
});
