import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { SiteSettingsPageComponent } from './site-settings-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('SiteSettingsPageComponent', () => {
  let fixture: ComponentFixture<SiteSettingsPageComponent>;
  let state: jasmine.SpyObj<WorkspaceStateService>;

  beforeEach(async () => {
    state = jasmine.createSpyObj<WorkspaceStateService>('WorkspaceStateService', ['selectedSite', 'templates', 'error', 'updateSelectedSite', 'uploadMediaFile']);
    state.selectedSite.and.returnValue({ id: 'site-1', name: 'Example Site', domain: 'https://example.test', blogPath: '/blog', templateKey: 'default-blog', themeConfig: '{"accent":"#2563eb"}' } as never);
    state.templates.and.returnValue([{ id: 'default', name: 'Default Blog', slug: 'default-blog', updatedAt: '', previewUrl: '/api/v1/template-previews/default-blog' }]);
    state.error.and.returnValue(null);
    state.updateSelectedSite.and.resolveTo();
    state.uploadMediaFile.and.resolveTo({ id: 'asset-1', fileUrl: 'https://cdn.example/logo.png' } as never);
    await TestBed.configureTestingModule({ imports: [RouterTestingModule, SiteSettingsPageComponent], providers: [{ provide: WorkspaceStateService, useValue: state }] }).compileComponents();
    fixture = TestBed.createComponent(SiteSettingsPageComponent);
    fixture.detectChanges();
  });

  it('saves the fields consumed by the renderer', async () => {
    await fixture.componentInstance.save();
    expect(state.updateSelectedSite).toHaveBeenCalledWith(jasmine.objectContaining({ name: 'Example Site', domain: 'https://example.test', blogPath: '/blog', templateKey: 'default-blog', themeConfig: '{"accent":"#2563eb"}' }));
  });

  it('uploads and assigns a logo media asset', async () => {
    const file = new File(['png'], 'logo.png', { type: 'image/png' });
    await fixture.componentInstance.uploadBrandAsset('logo', { target: { files: [file], value: 'logo.png' } } as unknown as Event);
    expect(state.uploadMediaFile).toHaveBeenCalledWith(file, 'Example Site logo');
    expect(state.updateSelectedSite).toHaveBeenCalledWith({ logoMediaId: 'asset-1' });
  });

  it('saves the selected AI writing context with the site', async () => {
    fixture.componentInstance.form.controls.contentContext.setValue('application_blog');
    await fixture.componentInstance.save();
    expect(state.updateSelectedSite).toHaveBeenCalledWith(jasmine.objectContaining({ contentContext: 'application_blog' }));
  });

  it('does not inherit Anonime\'s master prompt when the selected site has none saved', () => {
    state.selectedSite.and.returnValue({
      id: 'site-2', name: 'New Anonime Site', domain: 'https://new.example.test', blogPath: '/blog', templateKey: 'anonime', themeConfig: '{}', aiConfig: '{}',
    } as never);
    const newSiteFixture = TestBed.createComponent(SiteSettingsPageComponent);

    newSiteFixture.detectChanges();

    expect(newSiteFixture.componentInstance.form.controls.masterPrompt.value).toBe('');
  });

  it('persists the AI master prompt alongside the existing AI configuration', async () => {
    state.selectedSite.and.returnValue({
      id: 'site-1', name: 'Example Site', domain: 'https://example.test', blogPath: '/blog', templateKey: 'anonime', themeConfig: '{"accent":"#2563eb"}',
      aiConfig: '{"provider":"google","model":"gemini-2.5-flash","apiKeySecretRef":"GEMINI_API_KEY"}',
    } as never);
    fixture.componentInstance.form.controls.masterPrompt.setValue('Write accurate Anonime privacy guides.');

    await fixture.componentInstance.save();

    expect(state.updateSelectedSite).toHaveBeenCalledWith(jasmine.objectContaining({
      aiConfig: '{"provider":"google","model":"gemini-2.5-flash","apiKeySecretRef":"GEMINI_API_KEY","masterPrompt":"Write accurate Anonime privacy guides."}',
    }));
  });
});
