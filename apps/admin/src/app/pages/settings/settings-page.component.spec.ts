import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { SettingsPageComponent } from './settings-page.component';

describe('SettingsPageComponent', () => {
  let fixture: ComponentFixture<SettingsPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CommonModule, RouterTestingModule, SettingsPageComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(SettingsPageComponent);
    fixture.componentRef.setInput('authSession', {
      email: 'admin@example.com',
      fullName: 'Admin User',
      role: 'admin',
    });
    fixture.detectChanges();
  });

  it('renders shortcut groups and account details', () => {
    expect(fixture.nativeElement.textContent).toContain('Shortcut links');
    expect(fixture.nativeElement.textContent).toContain('Site setup');
    expect(fixture.nativeElement.querySelectorAll('.settings-link').length).toBeGreaterThan(0);
    expect(fixture.nativeElement.textContent).not.toContain('Maintain contributor profiles and ownership');
    expect(fixture.nativeElement.textContent).not.toContain('Organize content with stable taxonomy groups');
    expect(fixture.nativeElement.textContent).not.toContain('Open');
  });
});
