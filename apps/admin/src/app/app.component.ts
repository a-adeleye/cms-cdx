import { ChangeDetectionStrategy, Component, computed, inject } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { filter, map, startWith } from 'rxjs';
import { WORKSPACE_PAGES } from './features/pages/pages.data';
import { WorkspaceStateService } from './features/pages/workspace-state.service';

const PRIMARY_NAV_PATHS = ['dashboard', 'articles', 'publishing', 'settings'] as const;

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AppComponent {
  private readonly router = inject(Router);
  readonly state = inject(WorkspaceStateService);
  readonly currentUrl = toSignal(
    this.router.events.pipe(
      filter((event): event is NavigationEnd => event instanceof NavigationEnd),
      map((event) => event.urlAfterRedirects),
      startWith('/login'),
    ),
    { initialValue: '/login' },
  );

  readonly showWorkspaceShell = computed(() => !this.currentUrl().startsWith('/login'));
  readonly navItems = computed(() =>
    PRIMARY_NAV_PATHS.map((path) => WORKSPACE_PAGES.find((item) => item.path === path)).filter(
      (item): item is (typeof WORKSPACE_PAGES)[number] => item !== undefined,
    ),
  );

  async onSiteChange(siteId: string): Promise<void> {
    try {
      await this.state.selectSite(siteId);
    } catch {
      this.state.reportError('Unable to switch site.');
    }
  }

  async signOut(): Promise<void> {
    try {
      await this.state.logout();
      void this.router.navigate(['/login']);
    } catch {
      this.state.reportError('Unable to sign out cleanly.');
    }
  }
}
