import { CommonModule, NgOptimizedImage } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { map } from 'rxjs';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-media-asset-details-page',
  templateUrl: './media-asset-details-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './media-asset-details-page.component.css'],
  standalone: true,
  imports: [CommonModule, NgOptimizedImage, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MediaAssetDetailsPageComponent {
  private readonly route = inject(ActivatedRoute);
  readonly state = inject(WorkspaceStateService);
  readonly assetId = toSignal(this.route.paramMap.pipe(map((params) => params.get('assetId') ?? '')), { initialValue: '' });
  readonly asset = computed(() => this.state.mediaAssets().find((asset) => asset.id === this.assetId()) ?? null);
  readonly usageArticles = computed(() => {
    const assetUrl = this.asset()?.fileUrl;
    return assetUrl ? this.state.articles().filter((article) => article.coverImageUrl === assetUrl) : [];
  });

  formatSize(bytes: number): string {
    return bytes >= 1_048_576 ? `${(bytes / 1_048_576).toFixed(1)} MB` : `${Math.max(1, Math.round(bytes / 1024))} KB`;
  }
}
