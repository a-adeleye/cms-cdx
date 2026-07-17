import { ChangeDetectionStrategy, Component, effect, inject } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { map } from 'rxjs';
import { WORKSPACE_PAGES } from './pages.data';
import { WorkspacePageConfig } from './pages.model';
import { WorkspaceStateService } from './workspace-state.service';

@Component({
  selector: 'app-page-view',
  templateUrl: './page-view.component.html',
  styleUrl: './page-view.component.css',
  standalone: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
  host: {
    class: 'page-view',
  },
})
export class PageViewComponent {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);

  readonly page = toSignal(
    this.route.data.pipe(map((data) => (data['page'] as WorkspacePageConfig | undefined) ?? WORKSPACE_PAGES[1])),
    { initialValue: WORKSPACE_PAGES[1] },
  );
  readonly loginForm = this.fb.nonNullable.group({
    email: ['admin@example.com', [Validators.required, Validators.email]],
    password: ['admin123', [Validators.required, Validators.minLength(6)]],
  });

  constructor() {
    effect(() => {
      if (this.state.loading()) {
        return;
      }

      const page = this.page();
      if (!this.state.isAuthenticated() && page.kind !== 'login') {
        void this.router.navigate(['/login']);
        return;
      }

      if (this.state.isAuthenticated() && page.kind === 'login') {
        void this.router.navigate(['/dashboard']);
      }
    });

  }

  private reportActionError(message: string, error: unknown): void {
    const detail = error instanceof Error && error.message ? ` ${error.message}` : '';
    this.state.reportError(`${message}${detail}`);
  }

  async signIn(): Promise<void> {
    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      return;
    }

    const { email, password } = this.loginForm.getRawValue();
    try {
      await this.state.login(email, password);
      void this.router.navigate(['/dashboard']);
    } catch (error) {
      this.reportActionError('Unable to sign in.', error);
    }
  }

  openArticle(articleId: string): void {
    this.state.clearError();
    void this.router.navigate(['/content/articles', articleId, 'edit']);
  }

  async toggleSection(sectionId: string): Promise<void> {
    try {
      await this.state.toggleLandingSection(sectionId);
    } catch (error) {
      this.reportActionError('Unable to update landing section.', error);
    }
  }

  async moveSection(sectionId: string, direction: 'up' | 'down'): Promise<void> {
    try {
      await this.state.moveLandingSection(sectionId, direction);
    } catch (error) {
      this.reportActionError('Unable to reorder landing sections.', error);
    }
  }

}
