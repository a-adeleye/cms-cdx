import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, DestroyRef, effect, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { RouterModule } from '@angular/router';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import {
  DEPLOY_PROVIDER_OPTIONS,
  TEMPLATE_OPTIONS,
  defaultDeployConfigTemplate,
  isDeployConfigTemplate,
  shouldReplaceDeployConfigTemplate,
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
  private readonly destroyRef = inject(DestroyRef);
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);

  readonly siteSettingsForm = this.fb.nonNullable.group({
    templateKey: ['default-blog', [Validators.required]],
    themeConfig: ['', [Validators.required]],
    deployProvider: ['none', [Validators.required]],
    deployConfig: ['', [Validators.required]],
    previewDeployProvider: ['none', [Validators.required]],
    previewDeployConfig: ['', [Validators.required]],
    aiConfig: ['', [Validators.required]],
    storageConfig: ['', [Validators.required]],
  });

  readonly templateOptions = TEMPLATE_OPTIONS;
  readonly deployProviderOptions = DEPLOY_PROVIDER_OPTIONS;

  constructor() {
    this.siteSettingsForm.controls.deployProvider.valueChanges.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((provider) => {
      this.applyDeployConfigTemplate(provider, 'deployConfig');
    });
    this.siteSettingsForm.controls.previewDeployProvider.valueChanges
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe((provider) => {
        this.applyDeployConfigTemplate(provider, 'previewDeployConfig');
      });

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
          deployConfig: this.ensureDeployConfig(site.deployProvider || 'none', site.deployConfig),
          previewDeployProvider: site.previewDeployProvider || 'none',
          previewDeployConfig: this.ensureDeployConfig(site.previewDeployProvider || 'none', site.previewDeployConfig),
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

    const {
      templateKey,
      themeConfig,
      deployProvider,
      deployConfig,
      previewDeployProvider,
      previewDeployConfig,
      aiConfig,
      storageConfig,
    } = this.siteSettingsForm.getRawValue();
    try {
      await this.state.updateSelectedSite({
        templateKey,
        themeConfig,
        deployProvider,
        deployConfig,
        previewDeployProvider,
        previewDeployConfig,
        aiConfig,
        storageConfig,
      });
    } catch (error) {
      const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
      this.state.reportError(`Unable to save site settings.${detail}`);
    }
  }

  private applyDeployConfigTemplate(provider: string, controlName: 'deployConfig' | 'previewDeployConfig'): void {
    const control = this.siteSettingsForm.controls[controlName];
    if (!(shouldReplaceDeployConfigTemplate(provider, control.value) || isDeployConfigTemplate(control.value))) {
      return;
    }

    control.setValue(defaultDeployConfigTemplate(provider), { emitEvent: false });
  }

  private ensureDeployConfig(provider: string, value: string): string {
    return shouldReplaceDeployConfigTemplate(provider, value) ? defaultDeployConfigTemplate(provider) : value;
  }
}
