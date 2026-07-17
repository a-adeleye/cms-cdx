import { ComponentFixture, TestBed } from '@angular/core/testing';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { Component } from '@angular/core';
import { RichTextEditorComponent } from './rich-text-editor.component';

@Component({
  imports: [ReactiveFormsModule, RichTextEditorComponent],
  template: '<app-rich-text-editor [formControl]="content"></app-rich-text-editor>',
})
class RichTextEditorHostComponent {
  readonly content = new FormControl('A selected phrase');
}

describe('RichTextEditorComponent', () => {
  let fixture: ComponentFixture<RichTextEditorHostComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({ imports: [RichTextEditorHostComponent] }).compileComponents();
    fixture = TestBed.createComponent(RichTextEditorHostComponent);
    fixture.detectChanges();
  });

  it('wraps selected text in Markdown when bold is chosen', () => {
    const textarea = fixture.nativeElement.querySelector('textarea') as HTMLTextAreaElement;
    textarea.setSelectionRange(2, 10);
    clickToolbarButton(fixture.nativeElement.querySelector('[aria-label="Bold"]') as HTMLButtonElement);

    expect(fixture.componentInstance.content.value).toBe('A **selected** phrase');
  });

  it('turns selected lines into a block quote', () => {
    const textarea = fixture.nativeElement.querySelector('textarea') as HTMLTextAreaElement;
    textarea.value = 'First line\nSecond line';
    textarea.dispatchEvent(new Event('input'));
    textarea.setSelectionRange(0, textarea.value.length);
    clickToolbarButton(fixture.nativeElement.querySelector('[aria-label="Block quote"]') as HTMLButtonElement);

    expect(fixture.componentInstance.content.value).toBe('> First line\n> Second line');
  });

  it('captures the textarea selection before a toolbar button receives focus', () => {
    const boldButton = fixture.nativeElement.querySelector('[aria-label="Bold"]') as HTMLButtonElement;
    const event = new MouseEvent('mousedown', { bubbles: true, cancelable: true });

    boldButton.dispatchEvent(event);

    expect(event.defaultPrevented).toBeFalse();
  });
});

function clickToolbarButton(button: HTMLButtonElement): void {
  button.dispatchEvent(new MouseEvent('mousedown', { bubbles: true, cancelable: true }));
  button.click();
}
