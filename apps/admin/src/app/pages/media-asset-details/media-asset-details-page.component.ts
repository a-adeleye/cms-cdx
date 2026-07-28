import { CommonModule, DOCUMENT, NgOptimizedImage } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, effect, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { map } from 'rxjs';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-media-asset-details-page',
  templateUrl: './media-asset-details-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './media-asset-details-page.component.css'],
  standalone: true,
  imports: [CommonModule, NgOptimizedImage, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MediaAssetDetailsPageComponent {
  private readonly document = inject(DOCUMENT);
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  readonly state = inject(WorkspaceStateService);
  readonly feedback = createPageActionFeedback();
  readonly assetId = toSignal(this.route.paramMap.pipe(map((params) => params.get('assetId') ?? '')), { initialValue: '' });
  readonly asset = computed(() => this.state.mediaAssets().find((asset) => asset.id === this.assetId()) ?? null);
  readonly usageArticles = computed(() => {
    const assetUrl = this.asset()?.fileUrl;
    return assetUrl ? this.state.articles().filter((article) => article.coverImageUrl === assetUrl) : [];
  });
  readonly altTextForm = this.fb.nonNullable.group({
    altText: ['', [Validators.maxLength(500)]],
  });

  constructor() {
    effect(() => {
      this.altTextForm.controls.altText.setValue(this.asset()?.altText ?? '', { emitEvent: false });
    });
  }

  async saveAltText(): Promise<void> {
    const item = this.asset();
    if (!item || this.altTextForm.invalid) {
      this.altTextForm.markAllAsTouched();
      return;
    }
    try {
      this.feedback.loading('Saving alt text...');
      await this.state.updateMediaAltText(item.id, this.altTextForm.controls.altText.value);
      this.feedback.success('Alt text saved.');
    } catch (error) {
      this.feedback.error(error instanceof Error ? error.message : 'Unable to save alt text.');
    }
  }

  async replaceFile(event: Event): Promise<void> {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    const item = this.asset();
    if (!file || !item || this.altTextForm.invalid) {
      return;
    }
    try {
      this.feedback.loading('Replacing asset...');
      await this.state.replaceMediaFile(item.id, file, this.altTextForm.controls.altText.value);
      this.feedback.success('Asset replaced.');
    } catch (error) {
      this.feedback.error(error instanceof Error ? error.message : 'Unable to replace asset.');
    } finally {
      input.value = '';
    }
  }

  async deleteAsset(): Promise<void> {
    const item = this.asset();
    if (!item) {
      return;
    }
    const confirmed = this.document.defaultView?.confirm(`Delete ${item.fileName}? This cannot be undone.`) ?? false;
    if (!confirmed) {
      return;
    }
    try {
      this.feedback.loading('Deleting asset...');
      await this.state.deleteMediaAsset(item.id);
      this.feedback.success('Asset deleted.');
      await this.router.navigate(['/content/media']);
    } catch (error) {
      this.feedback.error(error instanceof Error ? error.message : 'Unable to delete asset.');
    }
  }

  formatSize(bytes: number): string {
    return bytes >= 1_048_576 ? `${(bytes / 1_048_576).toFixed(1)} MB` : `${Math.max(1, Math.round(bytes / 1024))} KB`;
  }
}
