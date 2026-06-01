import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, effect, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-site-settings-page',
  templateUrl: './site-settings-page.component.html',
  styleUrl: '../../features/pages/page-view.component.css',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SiteSettingsPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);

  readonly siteSettingsForm = this.fb.nonNullable.group({
    themeConfig: ['', [Validators.required]],
    deployProvider: ['none', [Validators.required]],
    deployConfig: ['', [Validators.required]],
    aiConfig: ['', [Validators.required]],
    storageConfig: ['', [Validators.required]],
  });

  constructor() {
    effect(() => {
      const site = this.state.selectedSite();
      if (!site) {
        return;
      }

      this.siteSettingsForm.reset(
        {
          themeConfig: site.themeConfig,
          deployProvider: site.deployProvider,
          deployConfig: site.deployConfig,
          aiConfig: site.aiConfig,
          storageConfig: site.storageConfig,
        },
        { emitEvent: false },
      );
    });
  }

  async saveSiteSettings(): Promise<void> {
    if (this.siteSettingsForm.invalid) {
      this.siteSettingsForm.markAllAsTouched();
      return;
    }

    const { themeConfig, deployProvider, deployConfig, aiConfig, storageConfig } = this.siteSettingsForm.getRawValue();
    try {
      await this.state.updateSelectedSite({
        themeConfig,
        deployProvider,
        deployConfig,
        aiConfig,
        storageConfig,
      });
    } catch (error) {
      const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
      this.state.reportError(`Unable to save site settings.${detail}`);
    }
  }
}
