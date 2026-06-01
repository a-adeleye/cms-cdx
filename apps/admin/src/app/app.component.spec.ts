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
  selector: 'app-settings-route-stub',
  standalone: false,
  template: '<p>Settings route is active</p>',
})
class SettingsRouteStubComponent {}

@Component({
  selector: 'app-sites-route-stub',
  standalone: false,
  template: '<p>Sites route is active</p>',
})
class SitesRouteStubComponent {}

describe('AppComponent', () => {
  let fixture: ComponentFixture<AppComponent>;
  let location: Location;
  let router: Router;
  const fakeState = {
    isAuthenticated: () => true,
    selectedSiteId: () => 'site-example',
    sites: () => [
      {
        id: 'site-example',
        name: 'Example Site',
        domain: 'https://example.test',
      },
    ],
    selectedSite: () => ({
      id: 'site-example',
      name: 'Example Site',
      domain: 'https://example.test',
      blogPath: '/articles',
      templateKey: 'default-blog',
    }),
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
    await TestBed.configureTestingModule({
      declarations: [AppComponent, LoginRouteStubComponent, DashboardRouteStubComponent, SettingsRouteStubComponent],
      imports: [
        RouterTestingModule.withRoutes([
          { path: '', pathMatch: 'full', redirectTo: 'login' },
          { path: 'login', component: LoginRouteStubComponent },
          { path: 'dashboard', component: DashboardRouteStubComponent },
        {
          path: 'settings',
          component: SettingsRouteStubComponent,
          children: [{ path: 'sites', component: SitesRouteStubComponent }],
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
    expect(links.map((link) => link.textContent?.trim())).toEqual(['Dashboard', 'Articles', 'Settings']);

    const dashboardLink = links.find((link) => link.textContent?.includes('Dashboard'));
    const settingsLink = links.find((link) => link.textContent?.includes('Settings'));

    expect(dashboardLink).toBeTruthy();
    expect(dashboardLink?.getAttribute('href')).toContain('/dashboard');
    expect(dashboardLink?.getAttribute('href')).not.toContain('#dashboard');
    expect(settingsLink?.getAttribute('href')).toContain('/settings');
    expect(settingsLink?.getAttribute('href')).not.toContain('#settings');

    dashboardLink?.click();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(location.path()).toBe('/dashboard');
    expect(fixture.nativeElement.textContent).toContain('Dashboard route is active');
    expect(fakeState.selectSite).toHaveBeenCalledWith('site-example');
  });

  it('keeps settings active on nested settings routes', async () => {
    await router.navigate(['/settings/sites']);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    const settingsLink = Array.from(
      fixture.nativeElement.querySelectorAll('nav a') as NodeListOf<HTMLAnchorElement>,
    ).find((link) => link.textContent?.includes('Settings'));

    expect(location.path()).toBe('/settings/sites');
    expect(settingsLink?.classList.contains('is-active')).toBeTrue();
    expect(fixture.nativeElement.textContent).toContain('Settings route is active');
  });

  it('shows the sign out action in the sidebar and returns to login when clicked', async () => {
    await router.navigate(['/dashboard']);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    const signOutButton = fixture.nativeElement.querySelector('.sidebar-logout') as HTMLButtonElement | null;
    expect(signOutButton).toBeTruthy();
    expect(signOutButton?.textContent?.trim()).toBe('Sign out');

    signOutButton?.click();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fakeState.logout).toHaveBeenCalled();
    expect(location.path()).toBe('/login');
  });
});
