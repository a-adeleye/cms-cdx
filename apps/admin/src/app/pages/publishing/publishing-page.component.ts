import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { BuildRecord, BuildType } from '../../features/pages/pages.model';

interface BuildActionFeedback {
  kind: 'error' | 'progress' | 'success';
  message: string;
}

interface EnvironmentBuildState {
  build: BuildRecord | null;
  isHealthy: boolean;
  label: string;
}

@Component({
  selector: 'app-publishing-page',
  templateUrl: './publishing-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './publishing-page.component.css'],
  standalone: true,
  imports: [CommonModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PublishingPageComponent {
  readonly state = inject(WorkspaceStateService);
  readonly actionFeedback = signal<BuildActionFeedback | null>(null);
  readonly isBuildPending = signal(false);
  readonly selectedArticleIds = signal<string[]>([]);

  readonly selectedArticleIdSet = computed(() => new Set(this.selectedArticleIds()));
  readonly selectedCount = computed(() => this.selectedArticleIds().length);
  readonly deploymentWarnings = computed(() => this.state.selectedSite().deploymentWarnings ?? []);
  readonly articles = computed(() => this.state.articles());
  readonly recentBuilds = computed(() => this.state.builds().slice(0, 8));
  readonly latestPreviewBuild = computed(() => this.latestEnvironmentBuild('preview'));
  readonly latestProductionBuild = computed(() => this.latestEnvironmentBuild('published'));
  readonly latestSuccessfulProductionDeployment = computed(
    () => this.state.builds().find((build) => build.buildType === 'published' && this.isSuccessfulDeployment(build)) ?? null,
  );
  readonly previewEnvironmentState = computed(() => this.environmentState('preview'));
  readonly productionEnvironmentState = computed(() => this.environmentState('published'));
  readonly selectedArticles = computed(() => this.articles().filter((article) => this.selectedArticleIdSet().has(article.id)));
  readonly isAllSelected = computed(() => {
    const articles = this.articles();
    return articles.length > 0 && this.selectedCount() === articles.length;
  });

  async triggerPreviewBuild(): Promise<void> {
    if (this.isBuildPending()) {
      return;
    }

    this.isBuildPending.set(true);
    this.actionFeedback.set({ kind: 'progress', message: 'Starting preview build...' });
    try {
      const build = await this.state.triggerBuild('preview', this.selectedArticleIds());
      this.actionFeedback.set(this.buildActionFeedback('Preview build', build));
    } catch (error) {
      this.reportActionError('Unable to start preview build.', error);
    } finally {
      this.isBuildPending.set(false);
    }
  }

  async triggerPublishedBuild(): Promise<void> {
    if (this.isBuildPending()) {
      return;
    }

    this.isBuildPending.set(true);
    this.actionFeedback.set({ kind: 'progress', message: 'Starting published build...' });
    try {
      for (const article of this.selectedArticles()) {
        if (article.status !== 'published') {
          await this.state.setArticleStatus(article.id, 'published');
        }
      }

      const build = await this.state.triggerBuild('published', this.selectedArticleIds());
      this.actionFeedback.set(this.buildActionFeedback('Published build', build));
    } catch (error) {
      this.reportActionError('Unable to start published build.', error);
    } finally {
      this.isBuildPending.set(false);
    }
  }

  toggleArticleSelection(articleId: string): void {
    const selectedIds = new Set(this.selectedArticleIds());
    if (selectedIds.has(articleId)) {
      selectedIds.delete(articleId);
    } else {
      selectedIds.add(articleId);
    }
    this.selectedArticleIds.set(Array.from(selectedIds));
  }

  toggleAllArticles(): void {
    if (this.isAllSelected()) {
      this.selectedArticleIds.set([]);
      return;
    }

    this.selectedArticleIds.set(this.articles().map((article) => article.id));
  }

  isArticleSelected(articleId: string): boolean {
    return this.selectedArticleIdSet().has(articleId);
  }

  duration(build: BuildRecord): string {
    if (!build.startedAt || !build.finishedAt) {
      return build.status === 'running' ? 'In progress' : '—';
    }

    const startedAt = Date.parse(build.startedAt);
    const finishedAt = Date.parse(build.finishedAt);
    if (!Number.isFinite(startedAt) || !Number.isFinite(finishedAt) || finishedAt < startedAt) {
      return 'Unavailable';
    }

    const seconds = Math.round((finishedAt - startedAt) / 1000);
    return seconds >= 60 ? `${Math.floor(seconds / 60)}m ${seconds % 60}s` : `${seconds}s`;
  }

  private reportActionError(message: string, error: unknown): void {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    const errorMessage = `${message}${detail}`;
    this.actionFeedback.set({ kind: 'error', message: errorMessage });
    this.state.reportError(errorMessage);
  }

  private latestEnvironmentBuild(buildType: BuildType): BuildRecord | null {
    return this.state.builds().find((build) => build.buildType === buildType) ?? null;
  }

  private environmentState(buildType: BuildType): EnvironmentBuildState {
    const build = this.latestEnvironmentBuild(buildType);
    return {
      build,
      isHealthy: Boolean(build && this.isSuccessfulBuildDeployment(build)),
      label: this.buildStatusLabel(build),
    };
  }

  private buildStatusLabel(build: BuildRecord | null): string {
    if (!build) {
      return 'Not deployed';
    }

    switch (build.status) {
      case 'success':
        return this.isSuccessfulBuildDeployment(build) ? 'Succeeded' : 'Build succeeded';
      case 'failed':
        return 'Failed';
      case 'running':
        return 'Running';
      case 'queued':
        return 'Queued';
    }
  }

  private isSuccessfulDeployment(build: BuildRecord): boolean {
    return this.isSuccessfulBuildDeployment(build) && Boolean(build.deployUrl.trim());
  }

  private isSuccessfulBuildDeployment(build: BuildRecord): boolean {
    return build.status === 'success' && build.deployStatus.trim().toLowerCase() === 'deployed';
  }

  private buildActionFeedback(prefix: string, build: BuildRecord): BuildActionFeedback {
    if (build.status === 'failed') {
      return { kind: 'error', message: `${prefix} failed.` };
    }
    if (build.status === 'running') {
      return { kind: 'progress', message: `${prefix} is running.` };
    }
    if (build.status === 'queued') {
      return { kind: 'progress', message: `${prefix} is queued.` };
    }

    const provider = build.deployProvider.trim();
    const deployUrl = build.deployUrl.trim();
    if (build.deployStatus.trim().toLowerCase() === 'deployed') {
      const destination = provider ? ` to ${provider}` : '';
      const link = deployUrl ? ` Open ${deployUrl}` : '';
      return { kind: 'success', message: `${prefix} deployed${destination}.${link}` };
    }

    if (build.deployStatus.trim()) {
      return { kind: 'error', message: `${prefix} completed, but deployment status is ${build.deployStatus.trim()}.` };
    }

    return { kind: 'success', message: `${prefix} completed successfully.` };
  }
}
