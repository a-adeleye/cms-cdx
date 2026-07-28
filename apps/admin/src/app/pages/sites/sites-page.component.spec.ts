import { CommonModule, DOCUMENT } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { SitesPageComponent } from './sites-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('SitesPageComponent', () => {
  let fixture: ComponentFixture<SitesPageComponent>;
  let fakeState: WorkspaceStateService;

  const selectedSite = {
    id: 'site-example',
    name: 'Example Site',
    slug: 'example',
    domain: 'https://example.test',
    blogPath: '/articles',
    contentContext: 'standalone_blog' as const,
    status: 'active' as const,
    templateKey: 'default-blog',
    themeConfig: '{}',
    deployProvider: '',
    deployConfig: '{}',
    previewDeployProvider: '',
    previewDeployConfig: '{}',
    aiConfig: '{}',
    storageConfig: '{}',
    updatedAt: '2026-05-23T00:00:00.000Z',
  };

  beforeEach(async () => {
    fakeState = {
      selectedSite: jasmine.createSpy('selectedSite').and.returnValue(selectedSite),
      templates: jasmine.createSpy('templates').and.returnValue([
        { id: 'template-default', name: 'Default Blog', slug: 'default-blog', updatedAt: '2026-05-23T00:00:00.000Z' },
        { id: 'template-premium', name: 'Premium SaaS', slug: 'premium-saas', updatedAt: '2026-05-23T00:00:00.000Z' },
      ]),
      sites: jasmine.createSpy('sites').and.returnValue([selectedSite, { ...selectedSite, id: 'site-other', name: 'Other Site' }]),
      selectedSiteId: jasmine.createSpy('selectedSiteId').and.returnValue('site-example'),
      createSite: jasmine.createSpy('createSite').and.resolveTo(selectedSite),
      updateSite: jasmine.createSpy('updateSite').and.resolveTo(selectedSite),
      deleteSite: jasmine.createSpy('deleteSite').and.resolveTo(),
      selectSite: jasmine.createSpy('selectSite').and.resolveTo(),
      reportError: jasmine.createSpy('reportError'),
    } as unknown as WorkspaceStateService;

    await TestBed.configureTestingModule({
      imports: [CommonModule, ReactiveFormsModule, RouterTestingModule, SitesPageComponent],
      providers: [{ provide: WorkspaceStateService, useValue: fakeState }],
    }).compileComponents();

    fixture = TestBed.createComponent(SitesPageComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('opens the create modal from the header action', async () => {
    const createButton = Array.from(fixture.nativeElement.querySelectorAll('button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.trim() === 'Create site',
    );

    createButton?.click();
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fixture.nativeElement.querySelector('dialog')?.open).toBeTrue();
    expect(fixture.nativeElement.querySelector('input[formcontrolname="name"]')?.value).toBe('');
  });

  it('shows slug validation feedback when create site is submitted with an invalid slug', async () => {
    const createButton = Array.from(fixture.nativeElement.querySelectorAll('button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.trim() === 'Create site',
    );

    createButton?.click();
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    const slugInput = fixture.nativeElement.querySelector('input[formcontrolname="slug"]') as HTMLInputElement | null;
    const submitButton = Array.from(fixture.nativeElement.querySelectorAll('dialog button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.trim() === 'Create site',
    );

    if (slugInput) {
      slugInput.value = 'Bad Slug';
      slugInput.dispatchEvent(new Event('input', { bubbles: true }));
      slugInput.dispatchEvent(new Event('blur', { bubbles: true }));
    }

    submitButton?.click();
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fixture.nativeElement.textContent).toContain('Slug must use lowercase letters, numbers, and hyphens only.');
    expect(fixture.nativeElement.textContent).toContain('Fix the highlighted fields to create the site.');
    expect(fakeState.createSite).not.toHaveBeenCalled();
  });

  it('keeps the create form focused on name and domain and derives an omitted slug', async () => {
    fixture.componentInstance.openCreateSiteModal();
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('details.site-advanced')).toBeTruthy();

    fixture.componentInstance.siteCreateForm.controls.name.setValue('New publication');
    fixture.componentInstance.siteCreateForm.controls.domain.setValue('https://new.example.test');

    await fixture.componentInstance.saveSiteFromModal();

    expect(fakeState.createSite).toHaveBeenCalledWith(
      jasmine.objectContaining({
        name: 'New publication',
        domain: 'https://new.example.test',
        slug: 'new-publication',
      }),
    );
  });

  it('opens the edit modal with the selected site values', async () => {
    const editButton = Array.from(fixture.nativeElement.querySelectorAll('button') as NodeListOf<HTMLButtonElement>).find(
      (button) => button.textContent?.trim() === 'Edit site',
    );

    editButton?.click();
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(fixture.nativeElement.querySelector('dialog')?.open).toBeTrue();
    expect(fixture.nativeElement.querySelector('input[formcontrolname="slug"]')?.value).toBe('example');
    expect(fixture.nativeElement.querySelector('select[formcontrolname="status"]')?.value).toBe('active');
  });

  it('deletes a site only after confirmation', async () => {
    spyOn(TestBed.inject(DOCUMENT).defaultView!, 'confirm').and.returnValue(true);

    await fixture.componentInstance.deleteSite(selectedSite);

    expect(fakeState.deleteSite).toHaveBeenCalledWith('site-example');
  });
});
