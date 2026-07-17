import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { of } from 'rxjs';
import { DashboardPageComponent } from '../../pages/dashboard/dashboard-page.component';
import { PageViewComponent } from './page-view.component';
import { WorkspaceStateService } from './workspace-state.service';

@Component({
  selector: 'app-route-stub',
  standalone: false,
  template: '<p>Route stub active</p>',
})
class RouteStubComponent {}

const loginPage = {
  path: 'login',
  navLabel: 'Login',
  kind: 'login' as const,
  eyebrow: 'Access',
  title: 'Login',
  primaryAction: { label: 'Open dashboard', path: '/dashboard' },
};

const dashboardPage = {
  path: 'dashboard',
  navLabel: 'Dashboard',
  kind: 'dashboard' as const,
  eyebrow: 'Overview',
  title: 'Dashboard',
  primaryAction: { label: 'Review articles', path: '/content/articles' },
  secondaryAction: { label: 'Manage sites', path: '/configuration/sites' },
};

describe('PageViewComponent', () => {
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
    uploadMediaFile: jasmine.createSpy('uploadMediaFile').and.resolveTo(),
  };

  function baseImports(routeConfig: any) {
    return [
      CommonModule,
      ReactiveFormsModule,
      RouterTestingModule.withRoutes(routeConfig),
      DashboardPageComponent,
    ];
  }

  it('renders only the sign-in form on the login page', async () => {
    await TestBed.configureTestingModule({
      declarations: [PageViewComponent],
      imports: baseImports([
        { path: 'login', component: RouteStubComponent },
      ]),
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

    const fixture = TestBed.createComponent(PageViewComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fixture.nativeElement.querySelector('.hero')).toBeNull();
    expect(fixture.nativeElement.querySelector('.two-column')).toBeNull();
    expect(fixture.nativeElement.querySelectorAll('.panel').length).toBe(1);
    expect(fixture.nativeElement.querySelector('form')).toBeTruthy();
  });

  it('hides the dashboard action buttons in the hero', async () => {
    await TestBed.configureTestingModule({
      declarations: [PageViewComponent],
      imports: baseImports([
        { path: 'dashboard', component: RouteStubComponent },
      ]),
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({ page: dashboardPage }),
          },
        },
        { provide: WorkspaceStateService, useValue: fakeState },
      ],
    }).compileComponents();

    const fixture = TestBed.createComponent(PageViewComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fixture.nativeElement.textContent).toContain('Dashboard');
    expect(fixture.nativeElement.textContent).not.toContain('Review articles');
    expect(fixture.nativeElement.textContent).not.toContain('Manage sites');
  });

});
