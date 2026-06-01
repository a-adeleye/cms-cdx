import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { AuthorsPageComponent } from './authors-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('AuthorsPageComponent', () => {
  let fixture: ComponentFixture<AuthorsPageComponent>;

  beforeEach(async () => {
    const state = {
      authors: () => [
        {
          id: 'author-1',
          siteId: 'site-example',
          name: 'Ada Lovelace',
          slug: 'ada-lovelace',
          bio: 'Computing pioneer.',
        },
      ],
      selectedSite: () => ({
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
        aiConfig: '{}',
        storageConfig: '{}',
        updatedAt: '2026-05-23T00:00:00.000Z',
      }),
    } as unknown as WorkspaceStateService;

    await TestBed.configureTestingModule({
      imports: [CommonModule, RouterTestingModule, AuthorsPageComponent],
      providers: [{ provide: WorkspaceStateService, useValue: state }],
    }).compileComponents();

    fixture = TestBed.createComponent(AuthorsPageComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('renders the authors overview and back link', () => {
    expect(fixture.nativeElement.textContent).toContain('Back to settings');
    expect(fixture.nativeElement.textContent).toContain('Ada Lovelace');
    expect(fixture.nativeElement.textContent).toContain('Contributor attribution');
  });
});
