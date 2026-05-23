import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { PageViewComponent } from './page-view.component';
import { PagesRoutingModule } from './pages-routing.module';
import { DashboardPageComponent } from '../../pages/dashboard/dashboard-page.component';
import { TaxonomyPageComponent } from '../../pages/taxonomy/taxonomy-page.component';
import { SettingsPageComponent } from '../../pages/settings/settings-page.component';

@NgModule({
  declarations: [PageViewComponent],
  imports: [CommonModule, ReactiveFormsModule, RouterModule, PagesRoutingModule, DashboardPageComponent, SettingsPageComponent, TaxonomyPageComponent],
})
export class PagesModule {}
