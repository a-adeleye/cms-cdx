import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, effect, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import {
  DEFAULT_TEMPLATE_SLUG,
  DEPLOY_PROVIDER_OPTIONS,
  templateSelectOptions,
} from '../../features/pages/site-config-options';

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
  readonly feedback = createPageActionFeedback();

  readonly siteSettingsForm = this.fb.nonNullable.group({
    templateKey: ['default-blog', [Validators.required]],
    themeConfig: ['', [Validators.required]],
    deployProvider: ['none', [Validators.required]],
    previewDeployProvider: ['none', [Validators.required]],
  });

  readonly templateOptions = computed(() =>
    templateSelectOptions(this.state.templates(), this.state.selectedSite().templateKey),
  );
  readonly deployProviderOptions = DEPLOY_PROVIDER_OPTIONS;

  constructor() {
    effect(() => {
      const site = this.state.selectedSite();
      if (!site) {
        return;
      }

      this.siteSettingsForm.reset(
        {
          templateKey: site.templateKey || this.templateOptions()[0]?.value || DEFAULT_TEMPLATE_SLUG,
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
      this.feedback.error('Fix the highlighted fields to save the site settings.');
      return;
    }

    const {
      templateKey,
      themeConfig,
      deployProvider,
      previewDeployProvider,
    } = this.siteSettingsForm.getRawValue();
    try {
      this.state.clearError();
      this.feedback.loading('Saving site settings...');
      await this.state.updateSelectedSite({
        templateKey,
        themeConfig,
        deployProvider,
        previewDeployProvider,
      });
      this.feedback.success('Site settings saved successfully.');
    } catch (error) {
      this.feedback.error(this.buildErrorMessage('Unable to save site settings.', error));
    }
  }

  private buildErrorMessage(message: string, error: unknown): string {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    return `${message}${detail}`;
  }
}
