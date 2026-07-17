import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { PublishingPageComponent } from './publishing-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('PublishingPageComponent', () => {
  let fixture: ComponentFixture<PublishingPageComponent>;
  let state: jasmine.SpyObj<WorkspaceStateService>;
  const article = { id: 'article-1', title: 'Real article', status: 'draft', updatedAt: '2026-07-17T00:00:00Z', coverImageUrl: '' };

  beforeEach(async () => {
    state = jasmine.createSpyObj<WorkspaceStateService>('WorkspaceStateService', ['selectedSite', 'articles', 'builds', 'error', 'triggerBuild', 'setArticleStatus', 'reportError']);
    state.selectedSite.and.returnValue({ name: 'Example', deployProvider: 'cloudflare_pages', previewDeployProvider: 'git_repository', deploymentWarnings: [] } as never);
    state.articles.and.returnValue([article] as never);
    state.builds.and.returnValue([]);
    state.error.and.returnValue(null);
    state.setArticleStatus.and.resolveTo();
    state.triggerBuild.and.resolveTo({ id: 'build-1', status: 'success', buildType: 'preview', deployStatus: 'deployed', deployProvider: 'git_repository', deployUrl: 'https://preview.test', logs: '', outputPath: '', siteId: '', startedAt: null, finishedAt: null } as never);
    await TestBed.configureTestingModule({ imports: [RouterTestingModule, PublishingPageComponent], providers: [{ provide: WorkspaceStateService, useValue: state }] }).compileComponents();
    fixture = TestBed.createComponent(PublishingPageComponent);
    fixture.detectChanges();
  });

  it('renders actual article and empty build history', () => {
    expect(fixture.nativeElement.textContent).toContain('Real article');
    expect(fixture.nativeElement.textContent).toContain('No build has been attempted');
  });

  it('passes selected IDs to preview and production builds', async () => {
    fixture.componentInstance.toggleArticleSelection('article-1');
    await fixture.componentInstance.triggerPreviewBuild();
    expect(state.triggerBuild).toHaveBeenCalledWith('preview', ['article-1']);
    await fixture.componentInstance.triggerPublishedBuild();
    expect(state.setArticleStatus).toHaveBeenCalledWith('article-1', 'published');
    expect(state.triggerBuild).toHaveBeenCalledWith('published', ['article-1']);
  });
});
