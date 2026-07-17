import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, input, output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ArticleRecord, BuildRecord, BuildType, SiteRecord } from '../../features/pages/pages.model';
import { SummaryMetric } from '../../features/pages/page-view.types';

@Component({
  selector: 'app-dashboard-page',
  templateUrl: './dashboard-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './dashboard-page.component.css'],
  standalone: true,
  imports: [CommonModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DashboardPageComponent {
  readonly dashboardStats = input.required<SummaryMetric[]>();
  readonly selectedSite = input<SiteRecord | null>(null);
  readonly recentArticles = input.required<ArticleRecord[]>();
  readonly recentBuilds = input.required<BuildRecord[]>();
  readonly articleSelected = output<string>();
  readonly latestProductionBuild = computed(() => this.recentBuilds().find((build) => build.buildType === 'published') ?? null);
  readonly latestPreviewBuild = computed(() => this.recentBuilds().find((build) => build.buildType === 'preview') ?? null);
  readonly latestLiveProductionDeployment = computed(
    () => this.recentBuilds().find((build) => build.buildType === 'published' && this.isLiveDeployment(build)) ?? null,
  );
  readonly productionDeploymentState = computed(() => this.deploymentState('published'));
  readonly previewDeploymentState = computed(() => this.deploymentState('preview'));
  readonly recentBuildActivity = computed(() =>
    this.recentBuilds()
      .slice(0, 3)
      .map((build) => ({
        build,
        title: this.activityTitle(build),
        statusLabel: this.buildStatusLabel(build),
      })),
  );

  selectArticle(articleId: string): void {
    this.articleSelected.emit(articleId);
  }

  private deploymentState(buildType: BuildType): { build: BuildRecord | null; label: string; isHealthy: boolean } {
    const build = this.recentBuilds().find((candidate) => candidate.buildType === buildType) ?? null;
    return {
      build,
      label: this.buildStatusLabel(build),
      isHealthy: Boolean(build && this.isSuccessfulDeployment(build)),
    };
  }

  private activityTitle(build: BuildRecord): string {
    const environment = build.buildType === 'published' ? 'Deployment' : 'Preview build';
    if (build.status === 'success' && !this.isSuccessfulDeployment(build)) {
      return `${environment} succeeded; deployment not recorded`;
    }
    return `${environment} ${this.buildStatusLabel(build).toLowerCase()}`;
  }

  private buildStatusLabel(build: BuildRecord | null): string {
    if (!build) {
      return 'Not deployed';
    }

    switch (build.status) {
      case 'success':
        return this.isSuccessfulDeployment(build) ? 'Succeeded' : 'Build succeeded';
      case 'failed':
        return 'Failed';
      case 'running':
        return 'Running';
      case 'queued':
        return 'Queued';
    }
  }

  private isSuccessfulDeployment(build: BuildRecord): boolean {
    return build.status === 'success' && build.deployStatus.trim().toLowerCase() === 'deployed';
  }

  private isLiveDeployment(build: BuildRecord): boolean {
    return this.isSuccessfulDeployment(build) && Boolean(build.deployUrl.trim());
  }
}
