import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { SummaryMetric } from '../../features/pages/page-view.types';
import { AuthorRecord } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-authors-page',
  templateUrl: './authors-page.component.html',
  styleUrl: '../../features/pages/page-view.component.css',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AuthorsPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);

  readonly editingId = signal<string>('');
  readonly metrics = computed<SummaryMetric[]>(() => [
    { label: 'Authors', value: String(this.state.authors().length), detail: 'Contributors available for attribution.' },
    { label: 'Editors', value: '1', detail: 'Editorial operators with review access.' },
  ]);

  readonly highlights = ['Contributor attribution', 'Role management', 'Reusable profiles'];
  readonly authors = computed(() => this.state.authors());

  readonly form = this.fb.nonNullable.group({
    id: [''],
    name: ['', [Validators.required, Validators.minLength(2)]],
    bio: ['', [Validators.required, Validators.minLength(12)]],
  });

  startCreate(): void {
    this.state.clearError();
    this.resetForm();
  }

  edit(author: AuthorRecord): void {
    this.state.clearError();
    this.editingId.set(author.id);
    this.form.reset(
      {
        id: author.id,
        name: author.name,
        bio: author.bio,
      },
      { emitEvent: false },
    );
  }

  async save(): Promise<void> {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const { id, name, bio } = this.form.getRawValue();
    try {
      this.state.clearError();
      await this.state.saveAuthor({ id: id.trim() || undefined, name, bio });
      this.resetForm();
    } catch (error) {
      this.reportActionError('Unable to save author.', error);
    }
  }

  async remove(author: AuthorRecord): Promise<void> {
    const confirmed = this.confirmDelete(author);
    if (!confirmed) {
      return;
    }

    try {
      this.state.clearError();
      await this.state.deleteAuthor(author.id);
      if (this.editingId() === author.id) {
        this.resetForm();
      }
    } catch (error) {
      this.reportActionError('Unable to delete author.', error);
    }
  }

  authorDetail(author: AuthorRecord): string {
    return author.bio ? author.bio : 'No bio provided.';
  }

  private resetForm(): void {
    this.editingId.set('');
    this.form.reset(
      {
        id: '',
        name: '',
        bio: '',
      },
      { emitEvent: false },
    );
  }

  private confirmDelete(author: AuthorRecord): boolean {
    if (typeof globalThis.confirm !== 'function') {
      return true;
    }

    return globalThis.confirm(`Delete ${author.name}? This cannot be undone.`);
  }

  private reportActionError(message: string, error: unknown): void {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    this.state.reportError(`${message}${detail}`);
  }
}
