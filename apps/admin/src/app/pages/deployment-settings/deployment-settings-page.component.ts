import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, effect, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
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
  readonly feedback = createPageActionFeedback();

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
      this.feedback.error('Fix the highlighted fields to save the deployment settings.');
      return;
    }

    const { deployConfig, previewDeployConfig, aiConfig, storageConfig } = this.deploymentSettingsForm.getRawValue();
    try {
      this.state.clearError();
      this.feedback.loading('Saving deployment settings...');
      await this.state.updateSelectedSite({
        deployConfig,
        previewDeployConfig,
        aiConfig,
        storageConfig,
      });
      this.feedback.success('Deployment settings saved successfully.');
    } catch (error) {
      this.feedback.error(this.buildErrorMessage('Unable to save deployment settings.', error));
    }
  }

  private ensureTemplate(value: string, template: string): string {
    if (!value.trim() || isJsonTemplate(value)) {
      return template;
    }

    return value;
  }

  private buildErrorMessage(message: string, error: unknown): string {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    return `${message}${detail}`;
  }
}
