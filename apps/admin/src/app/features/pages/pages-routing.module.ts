import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { WORKSPACE_PAGES } from './pages.data';
import { PageViewComponent } from './page-view.component';

const pageRoutes: Routes = WORKSPACE_PAGES.map((page) => ({
  path: page.path,
  component: PageViewComponent,
  data: { page },
}));

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'login',
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
