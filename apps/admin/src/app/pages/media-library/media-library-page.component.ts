import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-media-library-page',
  templateUrl: './media-library-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './media-library-page.component.css'],
  standalone: true,
  imports: [CommonModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MediaLibraryPageComponent {
  private readonly router = inject(Router);
  readonly state = inject(WorkspaceStateService);
  readonly search = signal('');
  readonly selectedIds = signal<string[]>([]);
  readonly uploading = signal(false);
  readonly storageBytes = computed(() => this.state.mediaAssets().reduce((total, asset) => total + asset.sizeBytes, 0));

  readonly assets = computed(() => {
    const search = this.search().trim().toLowerCase();
    return this.state.mediaAssets().filter((asset) => !search || `${asset.fileName} ${asset.altText} ${asset.mimeType}`.toLowerCase().includes(search));
  });

  onSearch(value: string): void {
    this.search.set(value);
  }

  toggleSelection(assetId: string): void {
    const selected = new Set(this.selectedIds());
    selected.has(assetId) ? selected.delete(assetId) : selected.add(assetId);
    this.selectedIds.set([...selected]);
  }

  toggleAll(): void {
    const visibleIds = this.assets().map((asset) => asset.id);
    const allSelected = visibleIds.length > 0 && visibleIds.every((id) => this.selectedIds().includes(id));
    this.selectedIds.set(allSelected ? [] : visibleIds);
  }

  isSelected(assetId: string): boolean {
    return this.selectedIds().includes(assetId);
  }

  openAsset(assetId: string): void {
    void this.router.navigate(['/content/media', assetId]);
  }

  async uploadFile(event: Event): Promise<void> {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;
    this.uploading.set(true);
    this.state.clearError();
    try {
      await this.state.uploadMediaFile(file, file.name.replace(/\.[^.]+$/, '').replace(/[-_]/g, ' '));
    } catch (error) {
      this.state.reportError(error instanceof Error ? error.message : 'Unable to upload media.');
    } finally {
      this.uploading.set(false);
      input.value = '';
    }
  }

  formatSize(bytes: number): string {
    if (bytes >= 1_048_576) return `${(bytes / 1_048_576).toFixed(1)} MB`;
    return `${Math.max(1, Math.round(bytes / 1024))} KB`;
  }

  usageCount(fileUrl: string): number {
    return this.state.articles().filter((article) => article.coverImageUrl === fileUrl).length;
  }
}
