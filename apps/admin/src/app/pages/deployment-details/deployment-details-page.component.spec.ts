import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { of } from 'rxjs';
import { BuildRecord } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { DeploymentDetailsPageComponent } from './deployment-details-page.component';

describe('DeploymentDetailsPageComponent', () => {
  const selectedSite = {
    id: 'site-example',
    name: 'Example Site',
    slug: 'example',
    domain: 'https://example.test',
    blogPath: '/articles',
    status: 'active' as const,
    templateKey: 'default-blog',
    themeConfig: '{}',
    deployProvider: 'netlify',
    deployConfig: '{}',
    previewDeployProvider: 'cloudflare',
    previewDeployConfig: '{}',
    aiConfig: '{}',
    storageConfig: '{}',
    updatedAt: '2026-05-23T00:00:00.000Z',
  };

  async function render(builds: BuildRecord[], buildId: string): Promise<ComponentFixture<DeploymentDetailsPageComponent>> {
    await TestBed.configureTestingModule({
      imports: [CommonModule, RouterTestingModule, DeploymentDetailsPageComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: { paramMap: of(convertToParamMap({ buildId })) },
        },
        {
          provide: WorkspaceStateService,
          useValue: {
            selectedSite: jasmine.createSpy('selectedSite').and.returnValue(selectedSite),
            builds: jasmine.createSpy('builds').and.returnValue(builds),
          },
        },
      ],
    }).compileComponents();

    const fixture = TestBed.createComponent(DeploymentDetailsPageComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
    return fixture;
  }

  afterEach(() => TestBed.resetTestingModule());

  it('renders only recorded failure events, output, and timing', async () => {
    const fixture = await render(
      [
        {
          id: 'failed-build',
          siteId: 'site-example',
          status: 'failed',
          buildType: 'published',
          logs: 'Provider rejected the deployment.',
          outputPath: 'dist/published/site',
          deployProvider: 'netlify',
          deployStatus: 'failed',
          deployUrl: '',
          startedAt: '2026-05-23T10:00:00.000Z',
          finishedAt: '2026-05-23T10:01:05.000Z',
        },
      ],
      'failed-build',
    );

    expect(fixture.nativeElement.textContent).toContain('Failed');
    expect(fixture.nativeElement.textContent).toContain('Deployment failed');
    expect(fixture.nativeElement.textContent).toContain('Provider rejected the deployment.');
    expect(fixture.nativeElement.textContent).toContain('dist/published/site');
    expect(fixture.nativeElement.textContent).toContain('1m 5s');
    expect(fixture.nativeElement.textContent).toContain('Build log');
    expect(fixture.nativeElement.textContent).not.toContain('Application built successfully');
    expect(fixture.nativeElement.textContent).not.toContain('/assets/app.js');
    expect(fixture.nativeElement.textContent).not.toContain('View live site');
  });

  it('renders a running build as incomplete with no fabricated artifacts', async () => {
    const fixture = await render(
      [
        {
          id: 'running-build',
          siteId: 'site-example',
          status: 'running',
          buildType: 'preview',
          logs: '',
          outputPath: '',
          deployProvider: 'cloudflare',
          deployStatus: 'running',
          deployUrl: '',
          startedAt: '2026-05-23T10:00:00.000Z',
          finishedAt: null,
        },
      ],
      'running-build',
    );

    expect(fixture.nativeElement.textContent).toContain('Running');
    expect(fixture.nativeElement.textContent).toContain('In progress');
    expect(fixture.nativeElement.textContent).toContain('Not generated');
    expect(fixture.nativeElement.textContent).not.toContain('Deployed');
    expect(fixture.nativeElement.textContent).not.toContain('Live');
  });

  it('renders a successful deployment URL and recorded output as live', async () => {
    const fixture = await render(
      [
        {
          id: 'successful-build',
          siteId: 'site-example',
          status: 'success',
          buildType: 'published',
          logs: 'Production deployment completed.',
          outputPath: 'dist/published/site',
          deployProvider: 'netlify',
          deployStatus: 'deployed',
          deployUrl: 'https://example.test',
          startedAt: '2026-05-23T10:00:00.000Z',
          finishedAt: '2026-05-23T10:00:30.000Z',
        },
      ],
      'successful-build',
    );

    expect(fixture.nativeElement.textContent).toContain('Successfully deployed to the production environment.');
    expect(fixture.nativeElement.textContent).toContain('Open deployment');
    expect(fixture.nativeElement.textContent).toContain('dist/published/site');
    expect(fixture.nativeElement.textContent).toContain('30s');
  });

  it('renders the not-found state for an unknown build id', async () => {
    const fixture = await render([], 'missing-build');

    expect(fixture.nativeElement.textContent).toContain('Deployment not found');
  });
});
