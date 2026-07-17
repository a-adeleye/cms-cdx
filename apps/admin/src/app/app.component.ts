import { ChangeDetectionStrategy, Component, ElementRef, ViewChild, computed, effect, inject } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { filter, map, startWith } from 'rxjs';
import { WorkspaceStateService } from './features/pages/workspace-state.service';

type NavigationItem = {
  label: string;
  path: string;
  iconPath: string;
  exact?: boolean;
};

type NavigationSection = {
  label: string;
  items: NavigationItem[];
};

const NAVIGATION_SECTIONS: NavigationSection[] = [
  {
    label: 'Content',
    items: [
      { label: 'Articles', path: '/content/articles', iconPath: 'M6 2h9l3 3v17H6z M14 2v5h4 M9 12h6 M9 16h6' },
      { label: 'Media Library', path: '/content/media', iconPath: 'M3 5h18v14H3z M7 10a2 2 0 1 0 0-4 2 2 0 0 0 0 4z M3 17l5-5 4 4 3-3 6 6', exact: true },
    ],
  },
  {
    label: 'Publishing',
    items: [
      { label: 'Publishing', path: '/publishing', iconPath: 'M4 13l5-2 4-7 7-2-2 7-7 4-2 5-1-6z M12 12l-3 3', exact: true },
      { label: 'Deployment History', path: '/publishing/history', iconPath: 'M12 8v5l3 2 M21 12a9 9 0 1 1-3-6.7 M18 2v5h-5' },
    ],
  },
  {
    label: 'Configuration',
    items: [
      { label: 'Sites', path: '/configuration/sites', iconPath: 'M12 21a9 9 0 1 0 0-18 9 9 0 0 0 0 18z M3 12h18 M12 3c3 3 3 15 0 18 M12 3c-3 3-3 15 0 18', exact: true },
      { label: 'Site settings', path: '/configuration/site-settings', iconPath: 'M12 15.5a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7z M19.4 15a1.7 1.7 0 0 0 .3 1.9l.1.1-2.8 2.8-.1-.1a1.7 1.7 0 0 0-1.9-.3 1.7 1.7 0 0 0-1 1.6V21h-4v-.1a1.7 1.7 0 0 0-1-1.6 1.7 1.7 0 0 0-1.9.3l-.1.1L4.2 17l.1-.1a1.7 1.7 0 0 0 .3-1.9A1.7 1.7 0 0 0 3 14H3v-4h.1a1.7 1.7 0 0 0 1.6-1 1.7 1.7 0 0 0-.3-1.9L4.2 7 7 4.2l.1.1a1.7 1.7 0 0 0 1.9.3A1.7 1.7 0 0 0 10 3V3h4v.1a1.7 1.7 0 0 0 1 1.6 1.7 1.7 0 0 0 1.9-.3l.1-.1L19.8 7l-.1.1a1.7 1.7 0 0 0-.3 1.9A1.7 1.7 0 0 0 21 10h.1v4H21a1.7 1.7 0 0 0-1.6 1z', exact: true },
      { label: 'Templates', path: '/configuration/templates', iconPath: 'M3 4h18v16H3z M3 9h18 M8 9v11' },
      { label: 'Taxonomy', path: '/configuration/taxonomy/categories', iconPath: 'M3 6h7l2 2h9v11H3z' },
      { label: 'Users', path: '/configuration/users', iconPath: 'M12 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8z M4 21c0-4 3.6-7 8-7s8 3 8 7' , exact: true },
      { label: 'AI', path: '/configuration/ai', iconPath: 'M12 3v3 M12 18v3 M3 12h3 M18 12h3 M5.6 5.6l2.1 2.1 M16.3 16.3l2.1 2.1 M18.4 5.6l-2.1 2.1 M7.7 16.3l-2.1 2.1 M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8z', exact: true },
      { label: 'Deployment', path: '/configuration/deployment', iconPath: 'M7 18H5a4 4 0 0 1-.5-8A7 7 0 0 1 18 8.5 4.5 4.5 0 0 1 18.5 18H17 M9 15l3-3 3 3 M12 12v9', exact: true },
    ],
  },
];

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AppComponent {
  @ViewChild('mainContent') private mainContent?: ElementRef<HTMLElement>;

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
  readonly navigationSections = NAVIGATION_SECTIONS;
  readonly dashboardIconPath = 'M3 11l9-8 9 8v10h-6v-6H9v6H3z';

  constructor() {
    effect(() => {
      const url = this.currentUrl();
      if (!url.startsWith('/login')) {
        setTimeout(() => this.mainContent?.nativeElement.focus());
      }
    });
  }

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
