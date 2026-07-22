import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { DeploymentSettingsPageComponent } from './deployment-settings-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('DeploymentSettingsPageComponent', () => {
  let fixture: ComponentFixture<DeploymentSettingsPageComponent>;
  let state: jasmine.SpyObj<WorkspaceStateService>;

  beforeEach(async () => {
    state = jasmine.createSpyObj<WorkspaceStateService>('WorkspaceStateService', ['selectedSite', 'builds', 'error', 'updateSelectedSite']);
    state.selectedSite.and.returnValue({ id: 'site-1', name: 'Example', deployProvider: 'none', deployConfig: '{}', previewDeployProvider: 'none', previewDeployConfig: '{}', aiConfig: '{}', storageConfig: '{}', deploymentWarnings: [] } as never);
    state.builds.and.returnValue([]);
    state.error.and.returnValue(null);
    state.updateSelectedSite.and.resolveTo();
    await TestBed.configureTestingModule({ imports: [RouterTestingModule, DeploymentSettingsPageComponent], providers: [{ provide: WorkspaceStateService, useValue: state }] }).compileComponents();
    fixture = TestBed.createComponent(DeploymentSettingsPageComponent);
    fixture.detectChanges();
  });

  it('serializes typed repository settings for each channel', async () => {
    fixture.componentInstance.deploymentSettingsForm.patchValue({
      deployProvider: 'git_repository', productionRepositoryUrl: 'https://github.com/example/site.git',
      productionBranch: 'main', productionContentPath: 'public/blog', productionTokenSecretRef: 'GITHUB_TOKEN',
    });
    await fixture.componentInstance.saveDeploymentSettings();
    expect(state.updateSelectedSite).toHaveBeenCalledWith(jasmine.objectContaining({
      deployProvider: 'git_repository',
      deployConfig: JSON.stringify({ repositoryUrl: 'https://github.com/example/site.git', branch: 'main', contentPath: 'public/blog', tokenSecretRef: 'GITHUB_TOKEN', publicUrl: '' }, null, 2),
    }));
  });

  it('shows only supported deployment modes', () => {
    const text = fixture.nativeElement.textContent;
    expect(text).toContain('Cloudflare Pages');
    expect(text).toContain('Firebase');
    expect(text).toContain('Repository branch');
    expect(text).not.toContain('Netlify');
  });

  it('stacks grouped fields only inside the preview card', () => {
    const cards = fixture.nativeElement.querySelectorAll('.deployment-config-card');
    expect(cards[0].classList).toContain('deployment-config-card--production');
    expect(cards[0].classList).not.toContain('deployment-config-card--preview');
    expect(cards[1].classList).toContain('deployment-config-card--preview');
    expect(cards[2].classList).toContain('deployment-config-card--storage');
    expect(cards[3].classList).toContain('deployment-config-card--latest');
  });

  it('serializes production storage fields and sends them with deployment settings', async () => {
    fixture.componentInstance.deploymentSettingsForm.patchValue({
      storageBucket: 'my-site-prod-media', storageRegion: 'us-west-2', storagePublicUrl: 'https://cdn.example.com',
      storageAccessKeySecretRef: 'PROD_S3_ACCESS_KEY', storageSecretKeySecretRef: 'PROD_S3_SECRET_KEY',
    });
    await fixture.componentInstance.saveDeploymentSettings();
    expect(state.updateSelectedSite).toHaveBeenCalledWith(jasmine.objectContaining({
      storageConfig: JSON.stringify({
        bucket: 'my-site-prod-media', region: 'us-west-2', endpoint: '', publicUrl: 'https://cdn.example.com',
        accessKeySecretRef: 'PROD_S3_ACCESS_KEY', secretKeySecretRef: 'PROD_S3_SECRET_KEY',
      }, null, 2),
    }));
  });

  it('leaves storage config empty when no fields are filled', async () => {
    await fixture.componentInstance.saveDeploymentSettings();
    expect(state.updateSelectedSite).toHaveBeenCalledWith(jasmine.objectContaining({ storageConfig: '{}' }));
  });

  it('rejects an incomplete production storage configuration', async () => {
    fixture.componentInstance.deploymentSettingsForm.patchValue({ storageBucket: 'my-site-prod-media' });
    await fixture.componentInstance.saveDeploymentSettings();
    fixture.detectChanges();
    expect(state.updateSelectedSite).not.toHaveBeenCalled();
    expect(fixture.nativeElement.textContent).toContain('Production storage publicUrl is required.');
  });
});
