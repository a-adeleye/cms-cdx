import { CommonModule, DOCUMENT } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { externalSiteUrl } from '../../features/pages/external-url';

@Component({
  selector: 'app-deployment-history-page',
  templateUrl: './deployment-history-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './deployment-history-page.component.css'],
  standalone: true,
  imports: [CommonModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DeploymentHistoryPageComponent {
  private readonly document = inject(DOCUMENT);
  readonly state = inject(WorkspaceStateService);
  readonly search = signal('');
  readonly environment = signal<'all' | 'preview' | 'published'>('all');
  readonly selectedSiteUrl = computed(() => externalSiteUrl(this.state.selectedSite().domain));
  readonly builds = computed(() => {
    const search = this.search().trim().toLowerCase();
    return this.state.builds().filter((build) => {
      const matchesEnvironment = this.environment() === 'all' || build.buildType === this.environment();
      return matchesEnvironment && (!search || `${build.id} ${build.deployProvider} ${build.status}`.toLowerCase().includes(search));
    });
  });

  setSearch(value: string): void { this.search.set(value); }
  setEnvironment(value: string): void { this.environment.set(value as 'all' | 'preview' | 'published'); }

  async clearHistory(): Promise<void> {
    if (this.state.builds().length === 0) {
      return;
    }
    const siteName = this.state.selectedSite().name || 'this site';
    const confirmed = this.document.defaultView?.confirm(`Clear all deployment history for ${siteName}? This cannot be undone.`) ?? false;
    if (!confirmed) {
      return;
    }

    try {
      await this.state.clearBuildHistory();
    } catch (error) {
      this.state.reportError(error instanceof Error ? error.message : 'Unable to clear deployment history.');
    }
  }

  shortId(id: string): string { return id.length > 8 ? id.slice(0, 8) : id; }
  deployNumber(index: number): number { return Math.max(1, this.state.builds().length - index); }
  duration(startedAt: string | null, finishedAt: string | null): string {
    if (!startedAt || !finishedAt) return '—';
    const seconds = Math.max(0, Math.round((Date.parse(finishedAt) - Date.parse(startedAt)) / 1000));
    return seconds >= 60 ? `${Math.floor(seconds / 60)}m ${seconds % 60}s` : `${seconds}s`;
  }
}
