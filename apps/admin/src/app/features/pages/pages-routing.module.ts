import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { WORKSPACE_PAGES } from './pages.data';
import { PageViewComponent } from './page-view.component';
import { SitesPageComponent } from '../../pages/sites/sites-page.component';
import { SiteSettingsPageComponent } from '../../pages/site-settings/site-settings-page.component';
import { AuthorsPageComponent } from '../../pages/authors/authors-page.component';
import { TaxonomyPageComponent } from '../../pages/taxonomy/taxonomy-page.component';

const pageRoutes: Routes = WORKSPACE_PAGES.filter((page) => page.path !== 'settings' && page.path !== 'sites').map((page) => ({
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
