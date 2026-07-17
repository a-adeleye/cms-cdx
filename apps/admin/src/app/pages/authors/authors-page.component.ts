import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { AuthorRecord } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { siteHostname } from '../../features/pages/external-url';

@Component({
  selector: 'app-authors-page',
  templateUrl: './authors-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './authors-page.component.css'],
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AuthorsPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);

  readonly editingId = signal<string>('');
  readonly search = signal('');
  readonly authors = computed(() => {
    const search = this.search().trim().toLowerCase();
    return this.state.authors().filter((author) => !search || `${author.name} ${author.slug} ${author.bio}`.toLowerCase().includes(search));
  });
  readonly selectedSiteHost = computed(() => siteHostname(this.state.selectedSite().domain));

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
    if (!this.confirmDelete(author)) {
      return;
    }

    try {
      this.state.clearError();
      await this.state.deleteAuthor(author.id);
      this.resetForm();
    } catch (error) {
      this.reportActionError('Unable to delete author.', error);
    }
  }

  articleCount(authorId: string): number {
    return this.state.articles().filter((article) => article.authorId === authorId).length;
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
    return typeof globalThis.confirm !== 'function' || globalThis.confirm(`Delete ${author.name}? This cannot be undone.`);
  }

  private reportActionError(message: string, error: unknown): void {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    this.state.reportError(`${message}${detail}`);
  }
}
