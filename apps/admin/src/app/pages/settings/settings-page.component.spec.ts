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
    expect(fixture.nativeElement.textContent).toContain('Admin User');
    expect(fixture.nativeElement.querySelectorAll('.settings-link').length).toBeGreaterThan(0);
  });
});
