import { ChangeDetectionStrategy, Component, ElementRef, computed, forwardRef, input, signal, viewChild } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

type BlockStyle = 'paragraph' | 'heading-2' | 'heading-3';
type InlineStyle = 'bold' | 'italic' | 'strikethrough' | 'code' | 'link' | 'image';
type LineStyle = 'quote' | 'bulleted-list' | 'numbered-list';

@Component({
  selector: 'app-rich-text-editor',
  templateUrl: './rich-text-editor.component.html',
  styleUrl: './rich-text-editor.component.css',
  providers: [{ provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => RichTextEditorComponent), multi: true }],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RichTextEditorComponent implements ControlValueAccessor {
  readonly placeholder = input('Start writing your article...');
  readonly invalid = input(false);
  readonly describedBy = input('');
  private readonly contentInput = viewChild.required<ElementRef<HTMLTextAreaElement>>('contentInput');
  readonly value = signal('');
  readonly disabled = signal(false);
  readonly wordCount = computed(() => this.value().trim() ? this.value().trim().split(/\s+/).length : 0);

  private onChange: (value: string) => void = () => undefined;
  private onTouched: () => void = () => undefined;
  private selectionStart = 0;
  private selectionEnd = 0;

  writeValue(value: string | null): void {
    this.value.set(value ?? '');
  }

  registerOnChange(onChange: (value: string) => void): void {
    this.onChange = onChange;
  }

  registerOnTouched(onTouched: () => void): void {
    this.onTouched = onTouched;
  }

  setDisabledState(isDisabled: boolean): void {
    this.disabled.set(isDisabled);
  }

  updateValue(event: Event): void {
    const input = event.target as HTMLTextAreaElement;
    this.commit(input.value);
    this.captureTextSelection(event);
  }

  markTouched(): void {
    this.onTouched();
  }

  applyBlockStyle(event: Event): void {
    const select = event.target as HTMLSelectElement;
    const style = select.value as BlockStyle;
    this.applyLineStyle(style);
    select.value = 'paragraph';
  }

  applyInlineStyle(style: InlineStyle): void {
    const wrappers: Record<InlineStyle, readonly [string, string, string]> = {
      bold: ['**', '**', 'bold text'],
      italic: ['*', '*', 'italic text'],
      strikethrough: ['~~', '~~', 'struck text'],
      code: ['`', '`', 'code'],
      link: ['[', '](https://)', 'link text'],
      image: ['![', '](https://)', 'image description'],
    };
    const [prefix, suffix, placeholder] = wrappers[style];
    this.withSelection((content, start, end) => {
      const selected = content.slice(start, end) || placeholder;
      const replacement = `${prefix}${selected}${suffix}`;
      return { content: content.slice(0, start) + replacement + content.slice(end), start: start + prefix.length, end: start + prefix.length + selected.length };
    });
  }

  applyLineStyle(style: BlockStyle | LineStyle): void {
    this.withSelection((content, start, end) => {
      const lineStart = content.lastIndexOf('\n', Math.max(0, start - 1)) + 1;
      const lineEndIndex = content.indexOf('\n', end);
      const lineEnd = lineEndIndex === -1 ? content.length : lineEndIndex;
      const lines = content.slice(lineStart, lineEnd).split('\n');
      const replacement = lines.map((line, index) => formatLine(line, style, index)).join('\n');
      return { content: content.slice(0, lineStart) + replacement + content.slice(lineEnd), start: lineStart, end: lineStart + replacement.length };
    });
  }

  handleKeyboardShortcut(event: KeyboardEvent): void {
    if (!(event.ctrlKey || event.metaKey) || event.altKey) {
      return;
    }
    if (event.key.toLowerCase() === 'b') {
      event.preventDefault();
      this.applyInlineStyle('bold');
    } else if (event.key.toLowerCase() === 'i') {
      event.preventDefault();
      this.applyInlineStyle('italic');
    }
  }

  preserveTextSelection(event: MouseEvent): void {
    if (event.target instanceof HTMLElement && event.target.closest('button')) {
      const textarea = this.contentInput().nativeElement;
      this.selectionStart = textarea.selectionStart;
      this.selectionEnd = textarea.selectionEnd;
    }
  }

  captureTextSelection(event: Event): void {
    const textarea = event.target as HTMLTextAreaElement;
    this.selectionStart = textarea.selectionStart;
    this.selectionEnd = textarea.selectionEnd;
  }

  private withSelection(transform: (content: string, start: number, end: number) => SelectionUpdate): void {
    if (this.disabled()) {
      return;
    }
    const textarea = this.contentInput().nativeElement;
    const update = transform(this.value(), this.selectionStart, this.selectionEnd);
    textarea.value = update.content;
    textarea.focus();
    textarea.setSelectionRange(update.start, update.end);
    this.selectionStart = update.start;
    this.selectionEnd = update.end;
    this.commit(update.content);
  }

  private commit(value: string): void {
    this.value.set(value);
    this.onChange(value);
  }
}

interface SelectionUpdate {
  content: string;
  start: number;
  end: number;
}

function formatLine(line: string, style: BlockStyle | LineStyle, index: number): string {
  const plainLine = line.replace(/^(#{1,6}\s+|>\s?|[-*+]\s+|\d+\.\s+)/, '');
  switch (style) {
    case 'heading-2': return `## ${plainLine}`;
    case 'heading-3': return `### ${plainLine}`;
    case 'quote': return `> ${plainLine}`;
    case 'bulleted-list': return `- ${plainLine}`;
    case 'numbered-list': return `${index + 1}. ${plainLine}`;
    default: return plainLine;
  }
}
