import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { PageViewComponent } from './page-view.component';
import { PagesRoutingModule } from './pages-routing.module';
import { DashboardPageComponent } from '../../pages/dashboard/dashboard-page.component';
import { SettingsPageComponent } from '../../pages/settings/settings-page.component';
import { ArticlesPageComponent } from '../../pages/articles/articles-page.component';
import { ArticleEditorPageComponent } from '../../pages/article-editor/article-editor-page.component';
import { ArticleDetailsPageComponent } from '../../pages/article-details/article-details-page.component';
import { DeploymentSettingsPageComponent } from '../../pages/deployment-settings/deployment-settings-page.component';
import { PublishingPageComponent } from '../../pages/publishing/publishing-page.component';

@NgModule({
  declarations: [PageViewComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    PagesRoutingModule,
    DashboardPageComponent,
    SettingsPageComponent,
    ArticlesPageComponent,
    ArticleEditorPageComponent,
    ArticleDetailsPageComponent,
    DeploymentSettingsPageComponent,
    PublishingPageComponent,
  ],
})
export class PagesModule {}
