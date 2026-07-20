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
import { MediaLibraryPageComponent } from '../../pages/media-library/media-library-page.component';
import { MediaAssetDetailsPageComponent } from '../../pages/media-asset-details/media-asset-details-page.component';
import { DeploymentHistoryPageComponent } from '../../pages/deployment-history/deployment-history-page.component';
import { DeploymentDetailsPageComponent } from '../../pages/deployment-details/deployment-details-page.component';
import { AiSettingsPageComponent } from '../../pages/ai-settings/ai-settings-page.component';

const loginPage = WORKSPACE_PAGES.find((page) => page.path === 'login')!;
const dashboardPage = WORKSPACE_PAGES.find((page) => page.path === 'dashboard')!;
const landingPageEditorPage = WORKSPACE_PAGES.find((page) => page.path === 'landing-page-editor')!;
const articlesPage = WORKSPACE_PAGES.find((page) => page.path === 'articles')!;
const mediaPage = WORKSPACE_PAGES.find((page) => page.path === 'media-library')!;
const publishingPage = WORKSPACE_PAGES.find((page) => page.path === 'publishing')!;

const articleRoutes: Routes = [
  {
    path: '',
    component: ArticlesPageComponent,
  },
  {
    path: 'new',
    component: ArticleEditorPageComponent,
    data: { editorMode: 'create' },
  },
  {
    path: ':articleId/edit',
    component: ArticleEditorPageComponent,
    data: { editorMode: 'edit' },
  },
  {
    path: 'editor',
    pathMatch: 'full',
    redirectTo: 'new',
  },
  {
    path: ':articleId',
    component: ArticleDetailsPageComponent,
  },
];

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'login',
  },
  {
    path: 'login',
    component: PageViewComponent,
    data: { page: loginPage },
  },
  {
    path: 'dashboard',
    component: PageViewComponent,
    data: { page: dashboardPage },
  },
  {
    path: 'content',
    children: [
      {
        path: 'articles',
        component: PageViewComponent,
        data: { page: articlesPage },
        children: articleRoutes,
      },
      {
        path: 'media',
        component: PageViewComponent,
        data: { page: mediaPage },
        children: [
          { path: '', component: MediaLibraryPageComponent },
          { path: ':assetId', component: MediaAssetDetailsPageComponent },
        ],
      },
    ],
  },
  {
    path: 'publishing',
    component: PageViewComponent,
    data: { page: publishingPage },
    children: [
      {
        path: '',
        component: PublishingPageComponent,
      },
      { path: 'history', component: DeploymentHistoryPageComponent },
      { path: 'history/:buildId', component: DeploymentDetailsPageComponent },
    ],
  },
  {
    path: 'configuration',
    children: [
      {
        path: '',
        pathMatch: 'full',
        redirectTo: 'sites',
      },
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
        path: 'templates/landing-page',
        component: PageViewComponent,
        data: { page: landingPageEditorPage },
      },
      {
        path: 'taxonomy/categories',
        component: TaxonomyPageComponent,
      },
      {
        path: 'users',
        component: AuthorsPageComponent,
      },
      {
        path: 'deployment',
        component: DeploymentSettingsPageComponent,
      },
      {
        path: 'ai',
        component: AiSettingsPageComponent,
      },
    ],
  },
  {
    path: 'articles',
    component: PageViewComponent,
    data: { page: articlesPage },
    children: articleRoutes,
  },
  {
    path: 'article-editor',
    redirectTo: 'content/articles/new',
    pathMatch: 'full',
  },
  {
    path: 'settings',
    children: [
      { path: '', pathMatch: 'full', redirectTo: '/configuration/sites' },
      { path: 'sites', pathMatch: 'full', redirectTo: '/configuration/sites' },
      { path: 'site-settings', pathMatch: 'full', redirectTo: '/configuration/site-settings' },
      { path: 'templates', pathMatch: 'full', redirectTo: '/configuration/templates' },
      { path: 'deployment-settings', pathMatch: 'full', redirectTo: '/configuration/deployment' },
      { path: 'ai', pathMatch: 'full', redirectTo: '/configuration/ai' },
      { path: 'authors', pathMatch: 'full', redirectTo: '/configuration/users' },
      { path: 'categories', pathMatch: 'full', redirectTo: '/configuration/taxonomy/categories' },
    ],
  },
  { path: 'sites', pathMatch: 'full', redirectTo: 'configuration/sites' },
  { path: 'site-settings', pathMatch: 'full', redirectTo: 'configuration/site-settings' },
  { path: 'deployment-settings', pathMatch: 'full', redirectTo: 'configuration/deployment' },
  { path: 'ai-settings', pathMatch: 'full', redirectTo: 'configuration/ai' },
  { path: 'templates', pathMatch: 'full', redirectTo: 'configuration/templates' },
  { path: 'authors', pathMatch: 'full', redirectTo: 'configuration/users' },
  { path: 'categories', pathMatch: 'full', redirectTo: 'configuration/taxonomy/categories' },
  { path: 'media-library', pathMatch: 'full', redirectTo: 'content/media' },
  { path: 'ai-assistant', pathMatch: 'full', redirectTo: 'content/articles/new' },
  { path: 'builds', pathMatch: 'full', redirectTo: 'publishing' },
  { path: 'landing-page-editor', pathMatch: 'full', redirectTo: 'configuration/templates/landing-page' },
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
