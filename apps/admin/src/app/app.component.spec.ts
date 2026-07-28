import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { AppComponent } from './app.component';
import { WorkspaceStateService } from './features/pages/workspace-state.service';

@Component({
  selector: 'app-login-route-stub',
  standalone: false,
  template: '<p>Login route is active</p>',
})
class LoginRouteStubComponent {}

@Component({
  selector: 'app-dashboard-route-stub',
  standalone: false,
  template: '<p>Dashboard route is active</p>',
})
class DashboardRouteStubComponent {}

@Component({
  selector: 'app-configuration-route-stub',
  standalone: false,
  template: '<p>Configuration route is active</p>',
})
class ConfigurationRouteStubComponent {}

@Component({
  selector: 'app-sites-route-stub',
  standalone: false,
  template: '<p>Sites route is active</p>',
})
class SitesRouteStubComponent {}

@Component({
  selector: 'app-articles-route-stub',
  standalone: false,
  template: '<p>Articles route is active</p>',
})
class ArticlesRouteStubComponent {}

@Component({
  selector: 'app-publishing-route-stub',
  standalone: false,
  template: '<p>Publishing route is active</p>',
})
class PublishingRouteStubComponent {}

describe('AppComponent', () => {
  let fixture: ComponentFixture<AppComponent>;
  let location: Location;
  let router: Router;
  let selectedSite: { id: string; name: string; domain: string; blogPath: string; templateKey: string; deployProvider: string; previewDeployProvider: string; previewDeployConfig: string; logoUrl?: string };
  const fakeState = {
    isAuthenticated: () => true,
    selectedSiteId: () => 'site-example',
    sites: () => [
      {
        id: 'site-example',
        name: 'Example Site',
        domain: 'https://example.test',
        previewDeployProvider: 'none',
        previewDeployConfig: '{}',
      },
    ],
    selectedSite: () => selectedSite,
    authSession: () => ({
      email: 'admin@example.com',
      fullName: 'Admin User',
      role: 'admin' as const,
    }),
    selectSite: jasmine.createSpy('selectSite').and.resolveTo(),
    logout: jasmine.createSpy('logout').and.resolveTo(),
    reportError: jasmine.createSpy('reportError'),
  };

  beforeEach(async () => {
    selectedSite = {
      id: 'site-example',
      name: 'Example Site',
      domain: 'https://example.test',
      blogPath: '/articles',
      templateKey: 'default-blog',
      deployProvider: 'netlify',
      previewDeployProvider: 'none',
      previewDeployConfig: '{}',
    };
    await TestBed.configureTestingModule({
      declarations: [
        AppComponent,
        LoginRouteStubComponent,
        DashboardRouteStubComponent,
        ConfigurationRouteStubComponent,
        SitesRouteStubComponent,
        ArticlesRouteStubComponent,
        PublishingRouteStubComponent,
      ],
      imports: [
        RouterTestingModule.withRoutes([
          { path: '', pathMatch: 'full', redirectTo: 'login' },
          { path: 'login', component: LoginRouteStubComponent },
          { path: 'dashboard', component: DashboardRouteStubComponent },
          { path: 'publishing', component: PublishingRouteStubComponent },
          {
            path: 'configuration',
            component: ConfigurationRouteStubComponent,
            children: [{ path: 'sites', component: SitesRouteStubComponent }],
          },
          {
            path: 'content',
            component: ArticlesRouteStubComponent,
            children: [
              {
                path: 'articles',
                component: ArticlesRouteStubComponent,
                children: [
                  { path: 'new', component: ArticlesRouteStubComponent },
                  { path: ':articleId/edit', component: ArticlesRouteStubComponent },
                ],
              },
              { path: 'media', component: ArticlesRouteStubComponent },
            ],
          },
        ]),
      ],
      providers: [{ provide: WorkspaceStateService, useValue: fakeState }],
    }).compileComponents();

    fixture = TestBed.createComponent(AppComponent);
    location = TestBed.inject(Location);
    router = TestBed.inject(Router);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('hides the workspace shell on the login route', async () => {
    await router.navigate(['/login']);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fixture.nativeElement.querySelector('.sidebar')).toBeNull();
    expect(fixture.nativeElement.textContent).toContain('Login route is active');
  });

  it('routes sidebar navigation to dashboard instead of using hash placeholders', async () => {
    await router.navigate(['/dashboard']);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    const siteSwitcher = fixture.nativeElement.querySelector('#site-switcher') as HTMLSelectElement | null;
    expect(siteSwitcher).toBeTruthy();
    expect(siteSwitcher?.value).toBe('site-example');

    siteSwitcher!.value = 'site-example';
    siteSwitcher!.dispatchEvent(new Event('change'));
    await fixture.whenStable();

    const links = Array.from(fixture.nativeElement.querySelectorAll('nav a')) as HTMLAnchorElement[];
    expect(links.map((link) => link.textContent?.trim())).toEqual([
      'Dashboard',
      'Articles',
      'Media Library',
      'Publishing',
      'Deployment History',
      'Sites',
      'Site settings',
      'Templates',
      'Taxonomy',
      'Users',
      'AI',
      'Deployment',
    ]);

    const dashboardLink = links.find((link) => link.textContent?.includes('Dashboard'));
    const sitesLink = links.find((link) => link.textContent?.trim() === 'Sites');

    expect(dashboardLink).toBeTruthy();
    expect(dashboardLink?.getAttribute('href')).toContain('/dashboard');
    expect(dashboardLink?.getAttribute('href')).not.toContain('#dashboard');
    expect(sitesLink?.getAttribute('href')).toContain('/configuration/sites');
    expect(sitesLink?.getAttribute('href')).not.toContain('#configuration');

    dashboardLink?.click();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(location.path()).toBe('/dashboard');
    expect(fixture.nativeElement.textContent).toContain('Dashboard route is active');
    expect(fakeState.selectSite).toHaveBeenCalledWith('site-example');
  });

  it('keeps configuration active on nested configuration routes', async () => {
    await router.navigate(['/configuration/sites']);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    const sitesLink = Array.from(
      fixture.nativeElement.querySelectorAll('nav a') as NodeListOf<HTMLAnchorElement>,
    ).find((link) => link.textContent?.trim() === 'Sites');

    expect(location.path()).toBe('/configuration/sites');
    expect(sitesLink?.classList.contains('is-active')).toBeTrue();
    expect(fixture.nativeElement.textContent).toContain('Configuration route is active');
  });

  it('keeps articles active on nested article routes', async () => {
    await router.navigate(['/content/articles/article-1/edit']);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    const articlesLink = Array.from(
      fixture.nativeElement.querySelectorAll('nav a') as NodeListOf<HTMLAnchorElement>,
    ).find((link) => link.textContent?.includes('Articles'));

    expect(location.path()).toBe('/content/articles/article-1/edit');
    expect(articlesLink?.classList.contains('is-active')).toBeTrue();
    expect(fixture.nativeElement.textContent).toContain('Articles route is active');
  });

  it('keeps publishing active on the publishing route', async () => {
    await router.navigate(['/publishing']);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    const publishingLink = Array.from(
      fixture.nativeElement.querySelectorAll('nav a') as NodeListOf<HTMLAnchorElement>,
    ).find((link) => link.textContent?.trim() === 'Publishing');

    expect(location.path()).toBe('/publishing');
    expect(publishingLink?.classList.contains('is-active')).toBeTrue();
    expect(fixture.nativeElement.textContent).toContain('Publishing route is active');
  });

  it('shows the sign out action in the sidebar and returns to login when clicked', async () => {
    await router.navigate(['/dashboard']);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    const signOutButton = fixture.nativeElement.querySelector('.sidebar-logout') as HTMLButtonElement | null;
    expect(signOutButton).toBeTruthy();
    expect(signOutButton?.textContent).toContain('Sign out');

    signOutButton?.click();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fakeState.logout).toHaveBeenCalled();
    expect(location.path()).toBe('/login');
  });

  it('renders the selected site logo in the sidebar brand mark', async () => {
    selectedSite = { ...selectedSite, logoUrl: 'https://cdn.example/logo.png' };
    await router.navigate(['/dashboard']);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    const logo = fixture.nativeElement.querySelector('.sidebar-brand__mark img') as HTMLImageElement | null;

    expect(logo?.src).toBe('https://cdn.example/logo.png');
    expect(logo?.alt).toBe('');
  });
});
