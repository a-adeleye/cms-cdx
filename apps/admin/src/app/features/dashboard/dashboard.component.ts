import { ChangeDetectionStrategy, Component } from '@angular/core';

type DashboardSection = {
  id: string;
  title: string;
  description: string;
};

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css',
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
  host: {
    class: 'dashboard',
  },
})
export class DashboardComponent {
  readonly sections: DashboardSection[] = [
    { id: 'login', title: 'Login', description: 'JWT-authenticated entry point.' },
    { id: 'dashboard', title: 'Dashboard', description: 'Build and content overview.' },
    { id: 'sites', title: 'Sites', description: 'Multi-site administration.' },
    { id: 'site-settings', title: 'Site Settings', description: 'Domain, theme, storage, and deployment configuration.' },
    { id: 'landing-page-editor', title: 'Landing Page Editor', description: 'Manage sections, reorder content, and preview JSON.' },
    { id: 'articles', title: 'Articles', description: 'Draft, review, publish, and archive workflows.' },
    { id: 'article-editor', title: 'Article Editor', description: 'Markdown, SEO, cover image, tags, and AI panel.' },
    { id: 'authors', title: 'Authors', description: 'Author profile management.' },
    { id: 'categories', title: 'Categories', description: 'Article taxonomy management.' },
    { id: 'tags', title: 'Tags', description: 'Reusable tagging system.' },
    { id: 'media-library', title: 'Media Library', description: 'S3-compatible uploads and asset management.' },
    { id: 'ai-assistant', title: 'AI Assistant', description: 'Draft generation and review-only content.' },
    { id: 'builds', title: 'Builds', description: 'Preview and published build logs.' },
    { id: 'deployment-settings', title: 'Deployment Settings', description: 'Provider-specific deployment configuration.' },
  ];
}
