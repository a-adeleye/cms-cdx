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
    expect(cards[2].classList).toContain('deployment-config-card--latest');
  });
});
