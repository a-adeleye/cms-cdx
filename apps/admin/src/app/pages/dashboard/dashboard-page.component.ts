import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, input, output } from '@angular/core';
import { ArticleRecord, SiteRecord } from '../../features/pages/pages.model';
import { SummaryMetric } from '../../features/pages/page-view.types';

@Component({
  selector: 'app-dashboard-page',
  templateUrl: './dashboard-page.component.html',
  styleUrl: '../../features/pages/page-view.component.css',
  standalone: true,
  imports: [CommonModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DashboardPageComponent {
  readonly dashboardStats = input.required<SummaryMetric[]>();
  readonly selectedSite = input<SiteRecord | null>(null);
  readonly recentArticles = input.required<ArticleRecord[]>();
  readonly articleSelected = output<string>();

  selectArticle(articleId: string): void {
    this.articleSelected.emit(articleId);
  }
}
