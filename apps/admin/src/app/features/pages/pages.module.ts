import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { PageViewComponent } from './page-view.component';
import { PagesRoutingModule } from './pages-routing.module';

@NgModule({
  declarations: [PageViewComponent],
  imports: [CommonModule, ReactiveFormsModule, PagesRoutingModule],
})
export class PagesModule {}
