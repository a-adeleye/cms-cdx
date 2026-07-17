import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { TemplatesPageComponent } from './templates-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('TemplatesPageComponent', () => {
  let fixture: ComponentFixture<TemplatesPageComponent>;
  let state: jasmine.SpyObj<WorkspaceStateService>;

  beforeEach(async () => {
    state = jasmine.createSpyObj<WorkspaceStateService>('WorkspaceStateService', ['selectedSite', 'templates', 'error', 'updateSelectedSite']);
    state.selectedSite.and.returnValue({ templateKey: 'default-blog' } as never);
    state.templates.and.returnValue([{ id: 'template-default', name: 'Default Blog', slug: 'default-blog', updatedAt: '2026-05-23T00:00:00Z', previewUrl: '/api/v1/template-previews/default-blog' }]);
    state.error.and.returnValue(null);
    state.updateSelectedSite.and.resolveTo();
    await TestBed.configureTestingModule({ imports: [RouterTestingModule, TemplatesPageComponent], providers: [{ provide: WorkspaceStateService, useValue: state }] }).compileComponents();
    fixture = TestBed.createComponent(TemplatesPageComponent);
    fixture.detectChanges();
  });

  it('renders the production-renderer preview instead of a mockup', () => {
    const frame = fixture.nativeElement.querySelector('iframe') as HTMLIFrameElement;
    expect(fixture.nativeElement.textContent).toContain('Default Blog');
    expect(frame.getAttribute('src')).toBe('/api/v1/template-previews/default-blog');
  });

  it('applies an implemented template', async () => {
    await fixture.componentInstance.selectTemplate('premium-saas');
    expect(state.updateSelectedSite).toHaveBeenCalledWith({ templateKey: 'premium-saas' });
  });
});
