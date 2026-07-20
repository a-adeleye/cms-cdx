import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
import { CategoryRecord } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-taxonomy-page',
  templateUrl: './taxonomy-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './taxonomy-page.component.css'],
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TaxonomyPageComponent {
  readonly state = inject(WorkspaceStateService);
  readonly feedback = createPageActionFeedback();

  private readonly fb = inject(FormBuilder);

  readonly form = this.fb.nonNullable.group({
    id: [''],
    name: ['', [Validators.required, Validators.minLength(2)]],
    description: [''],
  });

  readonly editingId = signal<string>('');
  readonly search = signal('');
  readonly records = computed(() => {
    const search = this.search().trim().toLowerCase();
    return this.state.categories().filter((record) => !search || `${record.name} ${record.slug}`.toLowerCase().includes(search));
  });
  readonly pageTitle = 'Categories';
  readonly singularLabel = 'category';
  readonly showDescription = true;

  articleCount(record: CategoryRecord): number {
    return this.state.articles().filter((article) => article.categoryId === record.id).length;
  }

  recordDescription(record: CategoryRecord): string {
    return record.description ? record.description : 'Organize related content and topics.';
  }

  startCreate(): void {
    this.state.clearError();
    this.feedback.clear();
    this.resetForm();
  }

  edit(record: CategoryRecord): void {
    this.state.clearError();
    this.feedback.clear();
    this.editingId.set(record.id);
    this.form.reset(
      {
        id: record.id,
        name: record.name,
        description: record.description,
      },
      { emitEvent: false },
    );
  }

  async save(): Promise<void> {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      this.feedback.error(`Fix the highlighted fields to save the ${this.singularLabel}.`);
      return;
    }

    const { id, name, description } = this.form.getRawValue();
    try {
      this.state.clearError();
      this.feedback.loading(`Saving ${this.singularLabel}...`);
      await this.state.saveCategory({ id: id.trim() || undefined, name, description });
      this.feedback.success(`${this.pageTitle} saved successfully.`);
      this.resetForm();
    } catch (error) {
      this.feedback.error(this.buildErrorMessage(`Unable to save ${this.singularLabel}.`, error));
    }
  }

  async remove(record: CategoryRecord): Promise<void> {
    if (!this.confirmDelete(record)) {
      return;
    }

    try {
      this.state.clearError();
      this.feedback.loading(`Deleting ${this.singularLabel}...`);
      await this.state.deleteCategory(record.id);
      this.feedback.success(`${this.pageTitle} deleted successfully.`);
      this.resetForm();
    } catch (error) {
      this.feedback.error(this.buildErrorMessage(`Unable to delete ${this.singularLabel}.`, error));
    }
  }

  private resetForm(): void {
    this.editingId.set('');
    this.form.reset(
      {
        id: '',
        name: '',
        description: '',
      },
      { emitEvent: false },
    );
  }

  private buildErrorMessage(message: string, error: unknown): string {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    return `${message}${detail}`;
  }

  private confirmDelete(record: CategoryRecord): boolean {
    return typeof globalThis.confirm !== 'function' || globalThis.confirm(`Delete ${record.name}? This cannot be undone.`);
  }
}
