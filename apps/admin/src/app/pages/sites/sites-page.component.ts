import { CommonModule } from '@angular/common';
import { DOCUMENT } from '@angular/common';
import { ChangeDetectionStrategy, Component, ElementRef, ViewChild, computed, effect, inject, signal } from '@angular/core';
import { FormBuilder, Validators, ReactiveFormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { TEMPLATE_OPTIONS } from '../../features/pages/site-config-options';
import { SiteRecord } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-sites-page',
  templateUrl: './sites-page.component.html',
  styleUrl: '../../features/pages/page-view.component.css',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SitesPageComponent {
  @ViewChild('siteDialog') private siteDialog?: ElementRef<HTMLDialogElement>;

  private readonly document = inject(DOCUMENT);
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);
  private siteDialogReturnFocus: HTMLElement | null = null;
  readonly templateOptions = TEMPLATE_OPTIONS;

  readonly siteDialogMode = signal<'create' | 'edit' | null>(null);
  readonly editingSiteId = signal<string | null>(null);
  readonly siteCreateForm = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    slug: ['', [Validators.required, Validators.pattern(/^[a-z0-9-]+$/)]],
    domain: ['', [Validators.required]],
    blogPath: ['/articles', [Validators.required]],
    templateKey: ['default-blog', [Validators.required]],
  });

  readonly siteEditForm = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    slug: ['', [Validators.required, Validators.pattern(/^[a-z0-9-]+$/)]],
    domain: ['', [Validators.required]],
    blogPath: ['/articles', [Validators.required]],
    templateKey: ['default-blog', [Validators.required]],
    status: ['active' as 'active' | 'inactive', [Validators.required]],
  });

  readonly sites = computed(() => this.state.sites());

  constructor() {
    effect(() => {
      const site = this.state.selectedSite();
      if (!site) {
        return;
      }

      this.siteEditForm.reset(
        {
          name: site.name,
          slug: site.slug,
          domain: site.domain,
          blogPath: site.blogPath,
          templateKey: site.templateKey,
          status: site.status,
        },
        { emitEvent: false },
      );
    });
  }

  openCreateSiteModal(): void {
    this.siteDialogMode.set('create');
    this.editingSiteId.set(null);
    this.resetCreateSiteForm();
    this.openSiteDialog();
  }

  openEditSiteModal(siteId: string): void {
    const site = this.state.sites().find((entry) => entry.id === siteId);
    if (!site) {
      this.state.reportError('Unable to open site editor.');
      return;
    }

    this.siteDialogMode.set('edit');
    this.editingSiteId.set(siteId);
    this.resetEditSiteForm(site);
    this.openSiteDialog();
  }

  async selectSite(siteId: string): Promise<void> {
    try {
      await this.state.selectSite(siteId);
    } catch {
      this.state.reportError('Unable to switch sites.');
    }
  }

  async saveSiteFromModal(): Promise<void> {
    const mode = this.siteDialogMode();
    if (!mode) {
      return;
    }

    if (mode === 'create' && this.siteCreateForm.invalid) {
      this.siteCreateForm.markAllAsTouched();
      return;
    }

    if (mode === 'edit' && this.siteEditForm.invalid) {
      this.siteEditForm.markAllAsTouched();
      return;
    }

    try {
      if (mode === 'create') {
        const site = await this.state.createSite(this.siteCreateForm.getRawValue());
        this.closeSiteDialog();
        this.resetCreateSiteForm();
        await this.state.selectSite(site.id);
        return;
      }

      const siteId = this.editingSiteId();
      if (!siteId) {
        this.state.reportError('Unable to open site editor.');
        return;
      }

      const { status, ...site } = this.siteEditForm.getRawValue();
      await this.state.updateSite(siteId, {
        ...site,
        status,
      });
      this.closeSiteDialog();
    } catch (error) {
      const message = mode === 'create' ? 'Unable to create site.' : 'Unable to save site.';
      this.state.reportError(error instanceof Error && error.message ? `${message} ${error.message}` : message);
    }
  }

  closeSiteDialog(): void {
    const dialog = this.siteDialog?.nativeElement;
    if (dialog?.open) {
      dialog.close();
      return;
    }

    this.onSiteDialogClosed();
  }

  onSiteDialogClosed(): void {
    this.siteDialogMode.set(null);
    this.editingSiteId.set(null);
    this.restoreSiteDialogFocus();
  }

  onSiteDialogCancel(event: Event): void {
    event.preventDefault();
    this.closeSiteDialog();
  }

  private openSiteDialog(): void {
    const dialog = this.siteDialog?.nativeElement;
    if (!dialog || dialog.open) {
      return;
    }

    this.captureSiteDialogFocus();
    queueMicrotask(() => {
      if (!dialog.open) {
        dialog.showModal();
      }
    });
  }

  private captureSiteDialogFocus(): void {
    const activeElement = this.document.activeElement;
    this.siteDialogReturnFocus = activeElement instanceof HTMLElement ? activeElement : null;
  }

  private restoreSiteDialogFocus(): void {
    this.siteDialogReturnFocus?.focus();
    this.siteDialogReturnFocus = null;
  }

  private resetCreateSiteForm(): void {
    this.siteCreateForm.reset({
      name: '',
      slug: '',
      domain: '',
      blogPath: '/articles',
      templateKey: 'default-blog',
    });
  }

  private resetEditSiteForm(site: SiteRecord): void {
    this.siteEditForm.reset(
      {
        name: site.name,
        slug: site.slug,
        domain: site.domain,
        blogPath: site.blogPath,
        templateKey: site.templateKey,
        status: site.status,
      },
      { emitEvent: false },
    );
  }
}
