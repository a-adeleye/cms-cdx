import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-templates-page',
  templateUrl: './templates-page.component.html',
  styleUrl: '../../features/pages/page-view.component.css',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TemplatesPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);
  readonly feedback = createPageActionFeedback();

  readonly templateForm = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    slug: ['', [Validators.required, Validators.pattern(/^[a-z0-9-]+$/)]],
  });

  readonly templates = computed(() => this.state.templates());

  async saveTemplate(): Promise<void> {
    if (this.templateForm.invalid) {
      this.templateForm.markAllAsTouched();
      this.feedback.error('Fix the highlighted fields to save the template.');
      return;
    }

    try {
      this.state.clearError();
      this.feedback.loading('Saving template...');
      await this.state.createTemplate(this.templateForm.getRawValue());
      this.feedback.success('Template registered successfully.');
      this.templateForm.reset({ name: '', slug: '' });
    } catch (error) {
      this.feedback.error(this.buildErrorMessage('Unable to create template.', error));
    }
  }

  private buildErrorMessage(message: string, error: unknown): string {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    return `${message}${detail}`;
  }
}
