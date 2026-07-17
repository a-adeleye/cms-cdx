import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { map } from 'rxjs';
import { BuildRecord, BuildType } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

interface DeploymentTimelineEvent {
  detail: string;
  id: string;
  time: string | null;
  title: string;
}

interface DeploymentTargetState {
  label: string;
  isLive: boolean;
}

@Component({
  selector: 'app-deployment-details-page',
  templateUrl: './deployment-details-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './deployment-details-page.component.css'],
  standalone: true,
  imports: [CommonModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DeploymentDetailsPageComponent {
  private readonly route = inject(ActivatedRoute);
  readonly state = inject(WorkspaceStateService);
  readonly buildId = toSignal(this.route.paramMap.pipe(map((params) => params.get('buildId') ?? '')), { initialValue: '' });
  readonly build = computed(() => this.state.builds().find((build) => build.id === this.buildId()) ?? null);
  readonly buildDuration = computed(() => this.duration(this.build()));
  readonly statusDescription = computed(() => this.describeStatus(this.build()));
  readonly timeline = computed(() => this.buildTimeline(this.build()));
  readonly productionTargetState = computed(() => this.deploymentTargetState('published'));
  readonly previewTargetState = computed(() => this.deploymentTargetState('preview'));
  readonly isCurrentBuildLive = computed(() => {
    const build = this.build();
    if (!build || !this.isSuccessfulDeployment(build)) {
      return false;
    }
    const latestLiveBuild = this.latestLiveDeployment(build.buildType);
    return latestLiveBuild?.id === build.id;
  });

  private duration(build: BuildRecord | null): string {
    if (!build?.startedAt) {
      return build?.status === 'queued' ? 'Not started' : 'Unavailable';
    }
    if (!build.finishedAt) {
      return build.status === 'running' ? 'In progress' : 'Not recorded';
    }

    const startedAt = Date.parse(build.startedAt);
    const finishedAt = Date.parse(build.finishedAt);
    if (!Number.isFinite(startedAt) || !Number.isFinite(finishedAt) || finishedAt < startedAt) {
      return 'Unavailable';
    }

    const seconds = Math.round((finishedAt - startedAt) / 1000);
    return seconds >= 60 ? `${Math.floor(seconds / 60)}m ${seconds % 60}s` : `${seconds}s`;
  }

  private describeStatus(build: BuildRecord | null): string {
    if (!build) {
      return '';
    }
    const environment = build.buildType === 'published' ? 'production' : 'preview';
    switch (build.status) {
      case 'success':
        return this.isSuccessfulDeployment(build)
          ? `Successfully deployed to the ${environment} environment.`
          : `Build succeeded, but no completed ${environment} deployment URL was recorded.`;
      case 'failed':
        return `The ${environment} deployment failed.`;
      case 'running':
        return `The ${environment} build is in progress.`;
      case 'queued':
        return `The ${environment} build is queued and has not started.`;
    }
  }

  private buildTimeline(build: BuildRecord | null): DeploymentTimelineEvent[] {
    if (!build) {
      return [];
    }

    const events: DeploymentTimelineEvent[] = [];
    if (!build.startedAt) {
      events.push({ id: 'queued', time: null, title: 'Build queued', detail: build.logs || 'No start time was recorded.' });
      return events;
    }

    events.push({ id: 'started', time: build.startedAt, title: 'Build started', detail: `${build.buildType === 'published' ? 'Production' : 'Preview'} build started.` });
    if (build.status === 'running') {
      events.push({ id: 'running', time: null, title: 'Build running', detail: build.logs || 'No completion was recorded.' });
      return events;
    }
    if (build.status === 'queued') {
      events.push({ id: 'queued', time: null, title: 'Build queued', detail: build.logs || 'Waiting to run.' });
      return events;
    }

    const deployed = this.isSuccessfulDeployment(build);
    events.push({
      id: 'completed',
      time: build.finishedAt,
      title:
        build.status === 'failed'
          ? build.deployStatus.trim()
            ? 'Deployment failed'
            : 'Build failed'
          : deployed
            ? 'Deployment succeeded'
            : 'Build succeeded',
      detail: build.logs || (build.status === 'failed' ? 'No failure detail was recorded.' : 'No completion detail was recorded.'),
    });
    return events;
  }

  private deploymentTargetState(buildType: BuildType): DeploymentTargetState {
    const latestBuild = this.state.builds().find((build) => build.buildType === buildType) ?? null;
    const liveBuild = this.latestLiveDeployment(buildType);
    if (liveBuild) {
      return {
        isLive: true,
        label: latestBuild?.status === 'failed' ? 'Live; latest build failed' : 'Live',
      };
    }

    return {
      isLive: false,
      label: this.statusLabel(latestBuild),
    };
  }

  private latestLiveDeployment(buildType: BuildType): BuildRecord | null {
    return this.state.builds().find((build) => build.buildType === buildType && this.isSuccessfulDeployment(build)) ?? null;
  }

  private isSuccessfulDeployment(build: BuildRecord): boolean {
    return build.status === 'success' && build.deployStatus.trim().toLowerCase() === 'deployed' && Boolean(build.deployUrl.trim());
  }

  private statusLabel(build: BuildRecord | null): string {
    if (!build) {
      return 'Not deployed';
    }
    switch (build.status) {
      case 'success':
        return 'Build succeeded';
      case 'failed':
        return 'Failed';
      case 'running':
        return 'Running';
      case 'queued':
        return 'Queued';
    }
  }
}
