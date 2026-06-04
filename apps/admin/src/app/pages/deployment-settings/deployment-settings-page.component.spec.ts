import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { DeploymentSettingsPageComponent } from './deployment-settings-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('DeploymentSettingsPageComponent', () => {
  let fixture: ComponentFixture<DeploymentSettingsPageComponent>;
  let state: WorkspaceStateService;
  let resolveUpdateSelectedSite: (() => void) | null;

  beforeEach(async () => {
    resolveUpdateSelectedSite = null;
    const selectedSite = {
      id: 'site-example',
      name: 'Example Site',
      slug: 'example',
      domain: 'https://example.test',
      blogPath: '/articles',
      status: 'active' as const,
      templateKey: 'default-blog',
      themeConfig: '{"tone":"professional"}',
      deployProvider: 'firebase',
      deployConfig: '',
      previewDeployProvider: 'cloudflare',
      previewDeployConfig: '',
      aiConfig: '',
      storageConfig: '',
      updatedAt: '2026-05-23T00:00:00.000Z',
    };

    state = {
      selectedSite: jasmine.createSpy('selectedSite').and.returnValue(selectedSite),
      error: jasmine.createSpy('error').and.returnValue(null),
      clearError: jasmine.createSpy('clearError'),
      updateSelectedSite: jasmine.createSpy('updateSelectedSite').and.callFake(
        () =>
          new Promise<void>((resolve) => {
            resolveUpdateSelectedSite = resolve;
          }),
      ),
      reportError: jasmine.createSpy('reportError'),
    } as unknown as WorkspaceStateService;

    await TestBed.configureTestingModule({
      imports: [CommonModule, ReactiveFormsModule, RouterTestingModule, DeploymentSettingsPageComponent],
      providers: [{ provide: WorkspaceStateService, useValue: state }],
    }).compileComponents();

    fixture = TestBed.createComponent(DeploymentSettingsPageComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('loads default templates and saves deployment settings', async () => {
    expect(fixture.nativeElement.textContent).toContain('Deployment settings');
    expect(fixture.nativeElement.textContent).toContain('firebase');

    const submitButton = fixture.nativeElement.querySelector('button[type="submit"]') as HTMLButtonElement | null;
    submitButton?.click();
    fixture.detectChanges();

    expect(submitButton?.disabled).toBeTrue();
    expect(fixture.nativeElement.textContent).toContain('Saving deployment settings...');

    resolveUpdateSelectedSite?.();

    await fixture.whenStable();
    fixture.detectChanges();

    expect(state.updateSelectedSite).toHaveBeenCalledWith({
      deployConfig: JSON.stringify(
        {
          provider: 'firebase',
          projectId: '',
          siteId: '',
          serviceAccountSecretRef: '',
        },
        null,
        2,
      ),
      previewDeployConfig: JSON.stringify(
        {
          provider: 'cloudflare',
          accountId: '',
          projectName: '',
          apiTokenSecretRef: '',
        },
        null,
        2,
      ),
      aiConfig: JSON.stringify(
        {
          provider: '',
          model: '',
          tone: '',
          brand_context: '',
        },
        null,
        2,
      ),
      storageConfig: JSON.stringify(
        {
          provider: '',
          bucket: '',
          region: '',
          prefix: '',
          public_url: '',
        },
        null,
        2,
      ),
    });
    expect(fixture.nativeElement.textContent).toContain('Deployment settings saved successfully.');
  });
});
