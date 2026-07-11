import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { AuthorsPageComponent } from './authors-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('AuthorsPageComponent', () => {
  let fixture: ComponentFixture<AuthorsPageComponent>;
  let state: {
    authors: () => Array<{ id: string; siteId: string; name: string; slug: string; bio: string }>;
    selectedSite: () => unknown;
    error: () => null;
    clearError: jasmine.Spy;
    reportError: jasmine.Spy;
    saveAuthor: jasmine.Spy;
    deleteAuthor: jasmine.Spy;
  };

  beforeEach(async () => {
    const author = {
      id: 'author-1',
      siteId: 'site-example',
      name: 'Ada Lovelace',
      slug: 'ada-lovelace',
      bio: 'Computing pioneer.',
    };

    state = {
      authors: () => [author],
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
        previewDeployProvider: 'cloudflare',
        previewDeployConfig: '{}',
        aiConfig: '{}',
        storageConfig: '{}',
        updatedAt: '2026-05-23T00:00:00.000Z',
      }),
      error: () => null,
      clearError: jasmine.createSpy('clearError'),
      reportError: jasmine.createSpy('reportError'),
      saveAuthor: jasmine.createSpy('saveAuthor').and.resolveTo(author),
      deleteAuthor: jasmine.createSpy('deleteAuthor').and.resolveTo(),
    };

    await TestBed.configureTestingModule({
      imports: [CommonModule, ReactiveFormsModule, RouterTestingModule, AuthorsPageComponent],
      providers: [{ provide: WorkspaceStateService, useValue: state }],
    }).compileComponents();

    fixture = TestBed.createComponent(AuthorsPageComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('creates, edits, and deletes authors', async () => {
    expect(fixture.nativeElement.textContent).toContain('Back to settings');
    expect(fixture.nativeElement.textContent).toContain('New author');
    expect(fixture.nativeElement.textContent).toContain('Ada Lovelace');
    expect(fixture.nativeElement.querySelectorAll('.table-row').length).toBeGreaterThan(1);

    const createButton = Array.from(fixture.nativeElement.querySelectorAll('button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.trim() === 'New author',
    );
    createButton?.click();
    fixture.detectChanges();

    fixture.componentInstance.form.controls.name.setValue('Grace Hopper');
    fixture.componentInstance.form.controls.bio.setValue('Pioneer in computer programming.');
    await fixture.componentInstance.save();

    expect(state.clearError).toHaveBeenCalled();
    expect(state.saveAuthor).toHaveBeenCalledWith({
      id: undefined,
      name: 'Grace Hopper',
      bio: 'Pioneer in computer programming.',
    });

    fixture.componentInstance.edit({
      id: 'author-1',
      siteId: 'site-example',
      name: 'Ada Lovelace',
      slug: 'ada-lovelace',
      bio: 'Computing pioneer.',
    });
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('input[formcontrolname="id"]')?.value).toBe('author-1');

    spyOn(window, 'confirm').and.returnValue(true);
    await fixture.componentInstance.remove({
      id: 'author-1',
      siteId: 'site-example',
      name: 'Ada Lovelace',
      slug: 'ada-lovelace',
      bio: 'Computing pioneer.',
    });

    expect(state.deleteAuthor).toHaveBeenCalledWith('author-1');
  });
});
