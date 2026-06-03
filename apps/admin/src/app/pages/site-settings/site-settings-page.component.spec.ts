import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { SiteSettingsPageComponent } from './site-settings-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { defaultDeployConfigTemplate } from '../../features/pages/site-config-options';

describe('SiteSettingsPageComponent', () => {
  let fixture: ComponentFixture<SiteSettingsPageComponent>;
  let state: WorkspaceStateService;

  beforeEach(async () => {
    const selectedSite = {
      id: 'site-example',
      name: 'Example Site',
      slug: 'example',
      domain: 'https://example.test',
      blogPath: '/articles',
      status: 'active' as const,
      templateKey: 'default-blog',
      themeConfig: '{"tone":"professional"}',
      deployProvider: 'netlify',
      deployConfig: '',
      previewDeployProvider: 'cloudflare',
      previewDeployConfig: '',
      aiConfig: '{}',
      storageConfig: '{}',
      updatedAt: '2026-05-23T00:00:00.000Z',
    };

    state = {
      selectedSite: jasmine.createSpy('selectedSite').and.returnValue(selectedSite),
      error: jasmine.createSpy('error').and.returnValue(null),
      updateSelectedSite: jasmine.createSpy('updateSelectedSite').and.resolveTo(),
      reportError: jasmine.createSpy('reportError'),
    } as unknown as WorkspaceStateService;

    await TestBed.configureTestingModule({
      imports: [CommonModule, ReactiveFormsModule, RouterTestingModule, SiteSettingsPageComponent],
      providers: [{ provide: WorkspaceStateService, useValue: state }],
    }).compileComponents();

    fixture = TestBed.createComponent(SiteSettingsPageComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('renders the current site and saves updates', async () => {
    expect(fixture.nativeElement.textContent).toContain('Back to settings');
    expect(fixture.nativeElement.textContent).toContain('Example Site');

    const submitButton = fixture.nativeElement.querySelector('button[type="submit"]') as HTMLButtonElement | null;
    submitButton?.click();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(state.updateSelectedSite).toHaveBeenCalledWith({
      templateKey: 'default-blog',
      themeConfig: '{"tone":"professional"}',
      deployProvider: 'netlify',
      deployConfig: defaultDeployConfigTemplate('netlify'),
      previewDeployProvider: 'cloudflare',
      previewDeployConfig: defaultDeployConfigTemplate('cloudflare'),
      aiConfig: '{}',
      storageConfig: '{}',
    });
  });
});
