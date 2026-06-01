import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, input, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { CategoryRecord, TagRecord } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

type TaxonomyKind = 'categories' | 'tags';
type TaxonomyRecord = CategoryRecord | TagRecord;

@Component({
  selector: 'app-taxonomy-page',
  templateUrl: './taxonomy-page.component.html',
  styleUrl: '../../features/pages/page-view.component.css',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TaxonomyPageComponent {
  readonly kind = input.required<TaxonomyKind>();
  readonly state = inject(WorkspaceStateService);

  private readonly fb = inject(FormBuilder);

  readonly form = this.fb.nonNullable.group({
    id: [''],
    name: ['', [Validators.required, Validators.minLength(2)]],
    description: [''],
  });

  readonly editingId = signal<string>('');
  readonly records = computed(() => (this.kind() === 'categories' ? this.state.categories() : this.state.tags()));
  readonly isCategoryPage = computed(() => this.kind() === 'categories');
  readonly pageTitle = computed(() => (this.isCategoryPage() ? 'Categories' : 'Tags'));
  readonly singularLabel = computed(() => (this.isCategoryPage() ? 'category' : 'tag'));
  readonly introText = computed(() =>
    this.isCategoryPage()
      ? 'Use categories to group articles into stable editorial sections.'
      : 'Use tags to label articles with reusable topics and campaign markers.',
  );
  readonly showDescription = computed(() => this.isCategoryPage());

  startCreate(): void {
    this.state.clearError();
    this.resetForm();
  }

  edit(record: TaxonomyRecord): void {
    this.state.clearError();
    this.editingId.set(record.id);
    this.form.reset(
      {
        id: record.id,
        name: record.name,
        description: 'description' in record ? record.description : '',
      },
      { emitEvent: false },
    );
  }

  async save(): Promise<void> {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const { id, name, description } = this.form.getRawValue();
    try {
      this.state.clearError();
      if (this.isCategoryPage()) {
        await this.state.saveCategory({ id: id.trim() || undefined, name, description });
      } else {
        await this.state.saveTag({ id: id.trim() || undefined, name });
      }
      this.resetForm();
    } catch (error) {
      this.reportActionError(`Unable to save ${this.singularLabel()}.`, error);
    }
  }

  async remove(record: TaxonomyRecord): Promise<void> {
    const confirmed = this.confirmDelete(record);
    if (!confirmed) {
      return;
    }

    try {
      this.state.clearError();
      if (this.isCategoryPage()) {
        await this.state.deleteCategory(record.id);
      } else {
        await this.state.deleteTag(record.id);
      }
      if (this.editingId() === record.id) {
        this.resetForm();
      }
    } catch (error) {
      this.reportActionError(`Unable to delete ${this.singularLabel()}.`, error);
    }
  }

  recordDetail(record: TaxonomyRecord): string {
    if (this.isCategoryPage()) {
      return 'description' in record && record.description ? record.description : 'No description provided.';
    }

    return 'Reusable topic label for article filtering.';
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

  private confirmDelete(record: TaxonomyRecord): boolean {
    if (typeof globalThis.confirm !== 'function') {
      return true;
    }

    return globalThis.confirm(`Delete ${record.name}? This cannot be undone.`);
  }

  private reportActionError(message: string, error: unknown): void {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    this.state.reportError(`${message}${detail}`);
  }
}
