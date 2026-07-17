import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-templates-page', templateUrl: './templates-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './templates-page.component.css'],
  standalone: true, imports: [CommonModule, RouterModule], changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TemplatesPageComponent {
  readonly state = inject(WorkspaceStateService);
  readonly feedback = createPageActionFeedback();
  private readonly sanitizer = inject(DomSanitizer);
  readonly search = signal('');
  readonly templates = computed(() => {
    const query = this.search().trim().toLowerCase();
    return this.state.templates().filter((template) => !query || `${template.name} ${template.slug}`.toLowerCase().includes(query));
  });

  async selectTemplate(templateKey: string): Promise<void> {
    try {
      this.feedback.loading('Applying template...');
      await this.state.updateSelectedSite({ templateKey });
      this.feedback.success('Template applied successfully.');
    } catch (error) {
      this.feedback.error(error instanceof Error ? error.message : 'Unable to apply template.');
    }
  }

  previewUrl(slug: string): SafeResourceUrl {
    return this.sanitizer.bypassSecurityTrustResourceUrl(`/api/v1/template-previews/${encodeURIComponent(slug)}`);
  }
}
