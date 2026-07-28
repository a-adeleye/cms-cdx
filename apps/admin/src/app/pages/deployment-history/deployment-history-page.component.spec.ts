import { DOCUMENT } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { DeploymentHistoryPageComponent } from './deployment-history-page.component';

describe('DeploymentHistoryPageComponent', () => {
  let fixture: ComponentFixture<DeploymentHistoryPageComponent>;
  const previewBuild = {
    id: 'preview-build-123', siteId: 'site-1', status: 'success', buildType: 'preview', logs: '', outputPath: 'dist',
    deployProvider: 'cloudflare', deployStatus: 'deployed', deployUrl: 'https://preview.example.test',
    startedAt: '2026-07-14T12:00:00.000Z', finishedAt: '2026-07-14T12:01:05.000Z',
  };
  const productionBuild = { ...previewBuild, id: 'production-build-1', buildType: 'published', deployProvider: 'netlify' };
  const state = {
    selectedSite: () => ({ name: 'Example Site', domain: 'https://anonime.io' }),
    builds: () => [previewBuild, productionBuild],
    clearBuildHistory: jasmine.createSpy('clearBuildHistory').and.resolveTo(),
    reportError: jasmine.createSpy('reportError'),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RouterTestingModule, DeploymentHistoryPageComponent],
      providers: [{ provide: WorkspaceStateService, useValue: state }],
    }).compileComponents();
    fixture = TestBed.createComponent(DeploymentHistoryPageComponent);
    fixture.detectChanges();
  });

  it('renders deployments and filters by environment and search', () => {
    expect(fixture.nativeElement.querySelectorAll('tbody tr').length).toBe(2);
    fixture.componentInstance.setEnvironment('preview');
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelectorAll('tbody tr').length).toBe(1);
    expect(fixture.nativeElement.textContent).toContain('preview-');
    fixture.componentInstance.setSearch('missing');
    fixture.detectChanges();
    expect(fixture.nativeElement.textContent).toContain('No deployments match');
  });

  it('normalizes the site link and formats deployment duration', () => {
    expect(fixture.componentInstance.duration(previewBuild.startedAt, previewBuild.finishedAt)).toBe('1m 5s');
  });

  it('clears the current site deployment history after confirmation', async () => {
    spyOn(TestBed.inject(DOCUMENT).defaultView!, 'confirm').and.returnValue(true);

    await fixture.componentInstance.clearHistory();

    expect(state.clearBuildHistory).toHaveBeenCalled();
  });
});
