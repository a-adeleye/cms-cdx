import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, effect, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import {
  defaultAiConfigTemplate,
  defaultDeployConfigTemplate,
  defaultStorageConfigTemplate,
  isJsonTemplate,
} from '../../features/pages/site-config-options';

@Component({
  selector: 'app-deployment-settings-page',
  templateUrl: './deployment-settings-page.component.html',
  styleUrl: '../../features/pages/page-view.component.css',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DeploymentSettingsPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);

  readonly deploymentSettingsForm = this.fb.nonNullable.group({
    deployConfig: ['', [Validators.required]],
    previewDeployConfig: ['', [Validators.required]],
    aiConfig: ['', [Validators.required]],
    storageConfig: ['', [Validators.required]],
  });

  constructor() {
    effect(() => {
      const site = this.state.selectedSite();
      if (!site) {
        return;
      }

      this.deploymentSettingsForm.reset(
        {
          deployConfig: this.ensureTemplate(site.deployConfig ?? '', defaultDeployConfigTemplate(site.deployProvider || 'none')),
          previewDeployConfig: this.ensureTemplate(
            site.previewDeployConfig ?? '',
            defaultDeployConfigTemplate(site.previewDeployProvider || 'none'),
          ),
          aiConfig: this.ensureTemplate(site.aiConfig ?? '', defaultAiConfigTemplate()),
          storageConfig: this.ensureTemplate(site.storageConfig ?? '', defaultStorageConfigTemplate()),
        },
        { emitEvent: false },
      );
    });
  }

  async saveDeploymentSettings(): Promise<void> {
    if (this.deploymentSettingsForm.invalid) {
      this.deploymentSettingsForm.markAllAsTouched();
      return;
    }

    const { deployConfig, previewDeployConfig, aiConfig, storageConfig } = this.deploymentSettingsForm.getRawValue();
    try {
      await this.state.updateSelectedSite({
        deployConfig,
        previewDeployConfig,
        aiConfig,
        storageConfig,
      });
    } catch (error) {
      const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
      this.state.reportError(`Unable to save deployment settings.${detail}`);
    }
  }

  private ensureTemplate(value: string, template: string): string {
    if (!value.trim() || isJsonTemplate(value)) {
      return template;
    }

    return value;
  }
}
