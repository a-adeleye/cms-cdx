import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, effect, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { DEPLOY_PROVIDER_OPTIONS, defaultDeployConfigTemplate } from '../../features/pages/site-config-options';
import { BuildRecord } from '../../features/pages/pages.model';

type Channel = 'production' | 'preview';

@Component({
  selector: 'app-deployment-settings-page',
  templateUrl: './deployment-settings-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './deployment-settings-page.component.css'],
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DeploymentSettingsPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);
  readonly feedback = createPageActionFeedback();
  readonly deployProviderOptions = DEPLOY_PROVIDER_OPTIONS;
  readonly latestBuild = computed(() => this.state.builds().at(0) ?? null);
  readonly latestSuccessfulPreviewDeployment = computed(() => this.state.builds().find((b) => b.buildType === 'preview' && this.isSuccessfulDeployment(b)) ?? null);
  readonly latestSuccessfulProductionDeployment = computed(() => this.state.builds().find((b) => b.buildType === 'published' && this.isSuccessfulDeployment(b)) ?? null);
  readonly latestBuildStatus = computed(() => this.buildStatusLabel(this.latestBuild()));
  readonly deploymentSettingsForm = this.fb.nonNullable.group({
    deployProvider: ['none', [Validators.required]], previewDeployProvider: ['none', [Validators.required]],
    deployConfig: ['{}', [Validators.required]], previewDeployConfig: ['{}', [Validators.required]], aiConfig: ['{}'], storageConfig: ['{}'],
    productionProjectName: [''], productionBranch: ['main'], productionFirebaseSiteId: [''], productionSecretRef: [''],
    productionRepositoryUrl: [''], productionContentPath: ['public/blog'], productionTokenSecretRef: [''], productionPublicUrl: [''],
    previewProjectName: [''], previewBranch: ['preview'], previewFirebaseSiteId: [''], previewSecretRef: [''],
    previewRepositoryUrl: [''], previewContentPath: ['public/blog'], previewTokenSecretRef: [''], previewPublicUrl: [''],
  });

  constructor() {
    effect(() => {
      const site = this.state.selectedSite();
      this.deploymentSettingsForm.reset({
        deployProvider: site.deployProvider || 'none', previewDeployProvider: site.previewDeployProvider || 'none',
        deployConfig: site.deployConfig || '{}', previewDeployConfig: site.previewDeployConfig || '{}',
        aiConfig: site.aiConfig || '{}', storageConfig: site.storageConfig || '{}',
        ...this.channelValues('production', site.deployConfig), ...this.channelValues('preview', site.previewDeployConfig),
      }, { emitEvent: false });
    });
  }

  provider(channel: Channel): string {
    return channel === 'production' ? this.deploymentSettingsForm.controls.deployProvider.value : this.deploymentSettingsForm.controls.previewDeployProvider.value;
  }

  resetDeploymentConfig(channel: Channel): void {
    const provider = this.provider(channel);
    const config = defaultDeployConfigTemplate(provider);
    const control = channel === 'production' ? this.deploymentSettingsForm.controls.deployConfig : this.deploymentSettingsForm.controls.previewDeployConfig;
    control.setValue(config);
    this.deploymentSettingsForm.patchValue(this.channelValues(channel, config));
  }

  async saveDeploymentSettings(): Promise<void> {
    try {
      const productionConfig = this.serializeChannel('production');
      const previewConfig = this.serializeChannel('preview');
      this.deploymentSettingsForm.controls.deployConfig.setValue(productionConfig);
      this.deploymentSettingsForm.controls.previewDeployConfig.setValue(previewConfig);
      this.feedback.loading('Saving deployment settings...');
      await this.state.updateSelectedSite({
        deployProvider: this.provider('production'), deployConfig: productionConfig,
        previewDeployProvider: this.provider('preview'), previewDeployConfig: previewConfig,
      });
      this.feedback.success('Deployment settings saved successfully.');
    } catch (error) {
      this.feedback.error(error instanceof Error ? error.message : 'Unable to save deployment settings.');
    }
  }

  private serializeChannel(channel: Channel): string {
    const p = channel;
    const controls = this.deploymentSettingsForm.controls;
    const provider = this.provider(channel);
    let config: Record<string, string> = {};
    if (provider === 'cloudflare_pages') config = { projectName: controls[`${p}ProjectName`].value.trim(), productionBranch: controls[`${p}Branch`].value.trim() };
    if (provider === 'firebase') config = { siteId: controls[`${p}FirebaseSiteId`].value.trim(), serviceAccountSecretRef: controls[`${p}SecretRef`].value.trim(), publicUrl: controls[`${p}PublicUrl`].value.trim() };
    if (provider === 'git_repository') config = { repositoryUrl: controls[`${p}RepositoryUrl`].value.trim(), branch: controls[`${p}Branch`].value.trim(), contentPath: controls[`${p}ContentPath`].value.trim(), tokenSecretRef: controls[`${p}TokenSecretRef`].value.trim(), publicUrl: controls[`${p}PublicUrl`].value.trim() };
    for (const [key, value] of Object.entries(config)) {
      if (!value && key !== 'publicUrl') throw new Error(`${channel === 'production' ? 'Production' : 'Preview'} ${key} is required.`);
    }
    return JSON.stringify(config, null, 2);
  }

  private channelValues(channel: Channel, raw: string): Record<string, string> {
    let value: Record<string, string> = {};
    try { value = JSON.parse(raw || '{}') as Record<string, string>; } catch { value = {}; }
    return {
      [`${channel}ProjectName`]: value['projectName'] || '', [`${channel}Branch`]: value['productionBranch'] || value['branch'] || (channel === 'production' ? 'main' : 'preview'),
      [`${channel}FirebaseSiteId`]: value['siteId'] || value['projectId'] || '', [`${channel}SecretRef`]: value['serviceAccountSecretRef'] || '',
      [`${channel}RepositoryUrl`]: value['repositoryUrl'] || '', [`${channel}ContentPath`]: value['contentPath'] || 'public/blog',
      [`${channel}TokenSecretRef`]: value['tokenSecretRef'] || '', [`${channel}PublicUrl`]: value['publicUrl'] || '',
    };
  }

  private isSuccessfulDeployment(build: BuildRecord): boolean { return build.status === 'success' && build.deployStatus === 'deployed' && Boolean(build.deployUrl); }
  private buildStatusLabel(build: BuildRecord | null): string { return build ? (build.status === 'success' ? 'Succeeded' : build.status[0].toUpperCase() + build.status.slice(1)) : 'No deployments yet'; }
}
