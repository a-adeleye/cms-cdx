import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SummaryMetric } from '../../features/pages/page-view.types';
import { AuthorRecord } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-authors-page',
  templateUrl: './authors-page.component.html',
  styleUrl: '../../features/pages/page-view.component.css',
  standalone: true,
  imports: [CommonModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AuthorsPageComponent {
  readonly state = inject(WorkspaceStateService);

  readonly metrics = computed<SummaryMetric[]>(() => [
    { label: 'Authors', value: String(this.state.authors().length), detail: 'Contributors available for attribution.' },
    { label: 'Editors', value: '1', detail: 'Editorial operators with review access.' },
  ]);

  readonly highlights = ['Contributor attribution', 'Role management', 'Reusable profiles'];
  readonly authors = computed(() => this.state.authors());

  authorDetail(author: AuthorRecord): string {
    return author.bio ? author.bio : 'No bio provided.';
  }
}
