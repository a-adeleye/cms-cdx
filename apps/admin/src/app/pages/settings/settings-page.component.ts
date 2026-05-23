import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { RouterModule } from '@angular/router';
import { AuthSession } from '../../features/pages/pages.model';
import { SettingsGroup } from '../../features/pages/page-view.types';

@Component({
  selector: 'app-settings-page',
  templateUrl: './settings-page.component.html',
  styleUrl: '../../features/pages/page-view.component.css',
  standalone: true,
  imports: [CommonModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SettingsPageComponent {
  readonly authSession = input<AuthSession | null>(null);
  readonly settingsGroups: SettingsGroup[] = [
    {
      title: 'Site setup',
      description: 'Configure the current site and switch between workspaces.',
      links: [
        { label: 'Sites', description: 'Create, select, and edit the active site.', path: '/sites' },
        { label: 'Site settings', description: 'Update domain, storage, AI, and deployment config.', path: '/site-settings' },
        { label: 'Landing page editor', description: 'Manage the home page sections for the selected site.', path: '/landing-page-editor' },
      ],
    },
    {
      title: 'Content structure',
      description: 'Keep editorial content organized and consistent.',
      links: [
        { label: 'Authors', description: 'Maintain contributor profiles and ownership.', path: '/authors' },
        { label: 'Categories', description: 'Organize content with stable taxonomy groups.', path: '/categories' },
        { label: 'Tags', description: 'Reuse topic labels across campaigns and pages.', path: '/tags' },
      ],
    },
    {
      title: 'Publishing',
      description: 'Handle assets, AI drafting, and deployment output.',
      links: [
        { label: 'Media library', description: 'Upload and manage reusable site assets.', path: '/media-library' },
        { label: 'AI assistant', description: 'Generate draft ideas without publishing automatically.', path: '/ai-assistant' },
        { label: 'Builds', description: 'Review preview and published build history.', path: '/builds' },
        { label: 'Deployment settings', description: 'Control deploy targets and secret references.', path: '/deployment-settings' },
      ],
    },
  ];
}
