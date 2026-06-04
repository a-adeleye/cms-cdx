import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { TemplatesPageComponent } from './templates-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('TemplatesPageComponent', () => {
  let fixture: ComponentFixture<TemplatesPageComponent>;
  let state: WorkspaceStateService;
  let resolveCreateTemplate: ((value: { id: string; name: string; slug: string; updatedAt: string }) => void) | null;

  beforeEach(async () => {
    resolveCreateTemplate = null;
    state = {
      templates: jasmine.createSpy('templates').and.returnValue([
        { id: 'template-default', name: 'Default Blog', slug: 'default-blog', updatedAt: '2026-05-23T00:00:00.000Z' },
      ]),
      error: jasmine.createSpy('error').and.returnValue(null),
      clearError: jasmine.createSpy('clearError'),
      createTemplate: jasmine.createSpy('createTemplate').and.callFake(
        () =>
          new Promise<{ id: string; name: string; slug: string; updatedAt: string }>((resolve) => {
            resolveCreateTemplate = resolve;
          }),
      ),
      reportError: jasmine.createSpy('reportError'),
    } as unknown as WorkspaceStateService;

    await TestBed.configureTestingModule({
      imports: [CommonModule, ReactiveFormsModule, RouterTestingModule, TemplatesPageComponent],
      providers: [{ provide: WorkspaceStateService, useValue: state }],
    }).compileComponents();

    fixture = TestBed.createComponent(TemplatesPageComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
  });

  it('renders registered templates', () => {
    expect(fixture.nativeElement.textContent).toContain('Registered templates');
    expect(fixture.nativeElement.textContent).toContain('Default Blog');
  });

  it('creates a template from the form', async () => {
    const nameInput = fixture.nativeElement.querySelector('input[formcontrolname="name"]') as HTMLInputElement | null;
    const slugInput = fixture.nativeElement.querySelector('input[formcontrolname="slug"]') as HTMLInputElement | null;
    if (nameInput) {
      nameInput.value = 'Magazine';
      nameInput.dispatchEvent(new Event('input', { bubbles: true }));
    }
    if (slugInput) {
      slugInput.value = 'magazine';
      slugInput.dispatchEvent(new Event('input', { bubbles: true }));
    }

    const submitButton = fixture.nativeElement.querySelector('button[type="submit"]') as HTMLButtonElement | null;
    submitButton?.click();
    fixture.detectChanges();

    expect(submitButton?.disabled).toBeTrue();
    expect(fixture.nativeElement.textContent).toContain('Saving template...');

    resolveCreateTemplate?.({
      id: 'template-new',
      name: 'Magazine',
      slug: 'magazine',
      updatedAt: '2026-05-23T00:00:00.000Z',
    });

    await fixture.whenStable();
    fixture.detectChanges();

    expect(state.createTemplate).toHaveBeenCalledWith({ name: 'Magazine', slug: 'magazine' });
    expect(fixture.nativeElement.textContent).toContain('Template registered successfully.');
  });
});
