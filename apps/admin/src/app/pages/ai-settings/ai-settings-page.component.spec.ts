import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { AiSettingsPageComponent } from './ai-settings-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('AiSettingsPageComponent', () => {
  let fixture: ComponentFixture<AiSettingsPageComponent>;
  let state: jasmine.SpyObj<WorkspaceStateService>;

  beforeEach(async () => {
    state = jasmine.createSpyObj<WorkspaceStateService>('WorkspaceStateService', ['selectedSite', 'error', 'updateSelectedSite']);
    state.selectedSite.and.returnValue({ id: 'site-1', name: 'Example Site', aiConfig: '{}' } as never);
    state.error.and.returnValue(null);
    state.updateSelectedSite.and.resolveTo();
    await TestBed.configureTestingModule({ imports: [RouterTestingModule, AiSettingsPageComponent], providers: [{ provide: WorkspaceStateService, useValue: state }] }).compileComponents();
    fixture = TestBed.createComponent(AiSettingsPageComponent);
    fixture.detectChanges();
  });

  it('offers provider selection and saves a secret reference instead of an API key', async () => {
    fixture.componentInstance.setProvider('openai');
    fixture.componentInstance.form.patchValue({ model: 'gpt-4.1-mini', apiKeySecretRef: 'OPENAI_API_KEY' });
    await fixture.componentInstance.save();

    expect(fixture.nativeElement.textContent).toContain('OpenAI');
    expect(state.updateSelectedSite).toHaveBeenCalledWith({ aiConfig: '{"provider":"openai","model":"gpt-4.1-mini","apiKeySecretRef":"OPENAI_API_KEY"}' });
  });
});
