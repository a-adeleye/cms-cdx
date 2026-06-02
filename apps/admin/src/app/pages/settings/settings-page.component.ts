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
      links: [
        { label: 'Sites', path: '/settings/sites' },
        { label: 'Site settings', path: '/settings/site-settings' },
        { label: 'Landing page editor', path: '/landing-page-editor' },
      ],
    },
    {
      title: 'Content structure',
      links: [
        { label: 'Authors', path: '/settings/authors' },
        { label: 'Categories', path: '/settings/categories' },
        { label: 'Tags', path: '/settings/tags' },
      ],
    },
    {
      title: 'Publishing',
      links: [
        { label: 'Publishing', path: '/publishing' },
        { label: 'Media library', path: '/media-library' },
        { label: 'AI assistant', path: '/ai-assistant' },
        { label: 'Builds', path: '/builds' },
        { label: 'Deployment settings', path: '/deployment-settings' },
      ],
    },
  ];
}
