import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, effect, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { DEPLOY_PROVIDER_OPTIONS, TEMPLATE_OPTIONS } from '../../features/pages/site-config-options';

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
    templateKey: ['default-blog', [Validators.required]],
    themeConfig: ['', [Validators.required]],
    deployProvider: ['none', [Validators.required]],
    previewDeployProvider: ['none', [Validators.required]],
  });

  readonly templateOptions = TEMPLATE_OPTIONS;
  readonly deployProviderOptions = DEPLOY_PROVIDER_OPTIONS;

  constructor() {
    effect(() => {
      const site = this.state.selectedSite();
      if (!site) {
        return;
      }

      this.siteSettingsForm.reset(
        {
          templateKey: site.templateKey || 'default-blog',
          themeConfig: site.themeConfig,
          deployProvider: site.deployProvider || 'none',
          previewDeployProvider: site.previewDeployProvider || 'none',
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

    const {
      templateKey,
      themeConfig,
      deployProvider,
      previewDeployProvider,
    } = this.siteSettingsForm.getRawValue();
    try {
      await this.state.updateSelectedSite({
        templateKey,
        themeConfig,
        deployProvider,
        previewDeployProvider,
      });
    } catch (error) {
      const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
      this.state.reportError(`Unable to save site settings.${detail}`);
    }
  }
}
