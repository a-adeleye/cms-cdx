import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, effect, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { DEFAULT_TEMPLATE_SLUG, templateSelectOptions } from '../../features/pages/site-config-options';
import { externalSiteUrl } from '../../features/pages/external-url';
import type { SiteContentContext } from '../../features/pages/pages.model';

@Component({
  selector: 'app-site-settings-page',
  templateUrl: './site-settings-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './site-settings-page.component.css'],
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SiteSettingsPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);
  readonly feedback = createPageActionFeedback();
  readonly selectedSiteUrl = computed(() => externalSiteUrl(this.state.selectedSite().domain));
  readonly templateOptions = computed(() => templateSelectOptions(this.state.templates(), this.state.selectedSite().templateKey));
  readonly form = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    domain: ['', [Validators.required]],
    blogPath: ['/blog', [Validators.required, Validators.pattern(/^\/[a-z0-9][a-z0-9/_-]*$/)]],
    description: ['', [Validators.maxLength(180)]],
    contentContext: ['standalone_blog' as SiteContentContext, [Validators.required]],
    templateKey: [DEFAULT_TEMPLATE_SLUG, [Validators.required]],
    accentColor: ['#2563eb', [Validators.required, Validators.pattern(/^#[0-9a-fA-F]{6}$/)]],
  });

  constructor() {
    effect(() => {
      const site = this.state.selectedSite();
      let theme: Record<string, unknown> = {};
      try { theme = JSON.parse(site.themeConfig || '{}') as Record<string, unknown>; } catch { theme = {}; }
      this.form.reset({
        name: site.name,
        domain: site.domain,
        blogPath: site.blogPath,
        description: site.description || '',
        contentContext: site.contentContext || 'standalone_blog',
        templateKey: site.templateKey || DEFAULT_TEMPLATE_SLUG,
        accentColor: typeof theme['accent'] === 'string' ? theme['accent'] : '#2563eb',
      }, { emitEvent: false });
    });
  }

  async save(): Promise<void> {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      this.feedback.error('Fix the highlighted fields before saving.');
      return;
    }
    const value = this.form.getRawValue();
    try {
      this.feedback.loading('Saving site configuration...');
      await this.state.updateSelectedSite({
        name: value.name.trim(), domain: value.domain.trim(), blogPath: value.blogPath.trim(),
        description: value.description.trim(), contentContext: value.contentContext, templateKey: value.templateKey,
        themeConfig: JSON.stringify({ accent: value.accentColor }),
      });
      this.feedback.success('Site configuration saved.');
    } catch (error) {
      this.feedback.error(error instanceof Error ? error.message : 'Unable to save site configuration.');
    }
  }

  async uploadBrandAsset(kind: 'logo' | 'favicon', event: Event): Promise<void> {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;
    try {
      this.feedback.loading(`Uploading ${kind}...`);
      const asset = await this.state.uploadMediaFile(file, `${this.form.controls.name.value} ${kind}`);
      await this.state.updateSelectedSite(kind === 'logo' ? { logoMediaId: asset.id } : { faviconMediaId: asset.id });
      this.feedback.success(`${kind === 'logo' ? 'Logo' : 'Favicon'} uploaded and assigned.`);
    } catch (error) {
      this.feedback.error(error instanceof Error ? error.message : `Unable to upload ${kind}.`);
    } finally {
      input.value = '';
    }
  }
}
