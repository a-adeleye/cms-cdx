import { CommonModule } from '@angular/common';
import { DOCUMENT } from '@angular/common';
import { ChangeDetectionStrategy, Component, ElementRef, ViewChild, computed, effect, inject, signal } from '@angular/core';
import { FormBuilder, Validators, ReactiveFormsModule } from '@angular/forms';
import {
  DEFAULT_TEMPLATE_SLUG,
  templateSelectOptions,
} from '../../features/pages/site-config-options';
import { SiteRecord } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

@Component({
  selector: 'app-sites-page',
  templateUrl: './sites-page.component.html',
  styleUrl: '../../features/pages/page-view.component.css',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SitesPageComponent {
  @ViewChild('siteDialog') private siteDialog?: ElementRef<HTMLDialogElement>;

  private readonly document = inject(DOCUMENT);
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);
  private siteDialogReturnFocus: HTMLElement | null = null;

  readonly siteDialogMode = signal<'create' | 'edit' | null>(null);
  readonly editingSiteId = signal<string | null>(null);
  readonly siteCreateForm = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    slug: ['', [Validators.pattern(/^[a-z0-9-]+$/)]],
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
  readonly templateOptions = computed(() =>
    templateSelectOptions(
      this.state.templates(),
      this.siteDialogMode() === 'edit'
        ? this.state.sites().find((site) => site.id === this.editingSiteId())?.templateKey ?? ''
        : this.siteCreateForm.controls.templateKey.value,
    ),
  );

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

  async deleteSite(site: SiteRecord): Promise<void> {
    if (this.sites().length <= 1) {
      this.state.reportError('Create another site before deleting the last site.');
      return;
    }
    const confirmed = this.document.defaultView?.confirm(`Delete ${site.name}? This permanently removes its CMS content.`) ?? false;
    if (!confirmed) {
      return;
    }

    try {
      await this.state.deleteSite(site.id);
    } catch (error) {
      const message = error instanceof Error && error.message ? `Unable to delete site. ${error.message}` : 'Unable to delete site.';
      this.state.reportError(message);
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
        const draft = this.siteCreateForm.getRawValue();
        const site = await this.state.createSite({
          ...draft,
          slug: draft.slug.trim() || this.createSlugFromName(draft.name),
        });
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
      templateKey: this.templateOptions()[0]?.value || DEFAULT_TEMPLATE_SLUG,
    });
  }

  private createSlugFromName(name: string): string {
    return name
      .trim()
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
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
