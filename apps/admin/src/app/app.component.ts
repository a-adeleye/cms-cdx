import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AppComponent {
  readonly navItems = [
    { id: 'login', label: 'Login' },
    { id: 'dashboard', label: 'Dashboard' },
    { id: 'sites', label: 'Sites' },
    { id: 'site-settings', label: 'Site Settings' },
    { id: 'landing-page-editor', label: 'Landing Page Editor' },
    { id: 'articles', label: 'Articles' },
    { id: 'article-editor', label: 'Article Editor' },
    { id: 'authors', label: 'Authors' },
    { id: 'categories', label: 'Categories' },
    { id: 'tags', label: 'Tags' },
    { id: 'media-library', label: 'Media Library' },
    { id: 'ai-assistant', label: 'AI Assistant' },
    { id: 'builds', label: 'Builds' },
    { id: 'deployment-settings', label: 'Deployment Settings' },
  ];
}
