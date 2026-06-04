import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { WORKSPACE_PAGES } from './pages.data';
import { PageViewComponent } from './page-view.component';
import { SitesPageComponent } from '../../pages/sites/sites-page.component';
import { SiteSettingsPageComponent } from '../../pages/site-settings/site-settings-page.component';
import { AuthorsPageComponent } from '../../pages/authors/authors-page.component';
import { TaxonomyPageComponent } from '../../pages/taxonomy/taxonomy-page.component';
import { ArticlesPageComponent } from '../../pages/articles/articles-page.component';
import { ArticleEditorPageComponent } from '../../pages/article-editor/article-editor-page.component';
import { ArticleDetailsPageComponent } from '../../pages/article-details/article-details-page.component';
import { PublishingPageComponent } from '../../pages/publishing/publishing-page.component';
import { DeploymentSettingsPageComponent } from '../../pages/deployment-settings/deployment-settings-page.component';
import { TemplatesPageComponent } from '../../pages/templates/templates-page.component';

const pageRoutes: Routes = WORKSPACE_PAGES.filter(
  (page) =>
    page.path !== 'settings' &&
    page.path !== 'sites' &&
    page.path !== 'articles' &&
    page.path !== 'publishing' &&
    page.path !== 'deployment-settings',
).map((page) => ({
  path: page.path,
  component: PageViewComponent,
  data: { page },
}));

const settingsPage = WORKSPACE_PAGES.find((page) => page.path === 'settings')!;

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'login',
  },
  {
    path: 'settings',
    component: PageViewComponent,
    data: { page: settingsPage },
    children: [
      {
        path: 'sites',
        component: SitesPageComponent,
      },
      {
        path: 'site-settings',
        component: SiteSettingsPageComponent,
      },
      {
        path: 'templates',
        component: TemplatesPageComponent,
      },
      {
        path: 'deployment-settings',
        component: DeploymentSettingsPageComponent,
      },
      {
        path: 'authors',
        component: AuthorsPageComponent,
      },
      {
        path: 'categories',
        component: TaxonomyPageComponent,
        data: { kind: 'categories' },
      },
      {
        path: 'tags',
        component: TaxonomyPageComponent,
        data: { kind: 'tags' },
      },
    ],
  },
  {
    path: 'articles',
    component: PageViewComponent,
    data: { page: WORKSPACE_PAGES.find((page) => page.path === 'articles')! },
    children: [
      {
        path: '',
        component: ArticlesPageComponent,
      },
      {
        path: 'editor',
        component: ArticleEditorPageComponent,
      },
      {
        path: ':articleId',
        component: ArticleDetailsPageComponent,
      },
    ],
  },
  {
    path: 'publishing',
    component: PageViewComponent,
    data: { page: WORKSPACE_PAGES.find((page) => page.path === 'publishing')! },
    children: [
      {
        path: '',
        component: PublishingPageComponent,
      },
    ],
  },
  {
    path: 'sites',
    redirectTo: 'settings/sites',
    pathMatch: 'full',
  },
  {
    path: 'site-settings',
    redirectTo: 'settings/site-settings',
    pathMatch: 'full',
  },
  {
    path: 'deployment-settings',
    redirectTo: 'settings/deployment-settings',
    pathMatch: 'full',
  },
  {
    path: 'templates',
    redirectTo: 'settings/templates',
    pathMatch: 'full',
  },
  {
    path: 'authors',
    redirectTo: 'settings/authors',
    pathMatch: 'full',
  },
  {
    path: 'categories',
    redirectTo: 'settings/categories',
    pathMatch: 'full',
  },
  {
    path: 'tags',
    redirectTo: 'settings/tags',
    pathMatch: 'full',
  },
  {
    path: 'article-editor',
    redirectTo: 'articles/editor',
    pathMatch: 'full',
  },
  ...pageRoutes,
  {
    path: '**',
    redirectTo: 'login',
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class PagesRoutingModule {}
