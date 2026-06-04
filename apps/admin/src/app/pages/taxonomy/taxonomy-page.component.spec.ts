import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { TaxonomyPageComponent } from './taxonomy-page.component';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

describe('TaxonomyPageComponent', () => {
  const category = {
    id: 'category-1',
    siteId: 'site-example',
    name: 'Category',
    slug: 'category',
    description: 'Category description',
  };

  const tag = {
    id: 'tag-1',
    siteId: 'site-example',
    name: 'Tag',
    slug: 'tag',
  };

  let fixture: ComponentFixture<TaxonomyPageComponent>;

  const createState = () => ({
    loading: () => false,
    isAuthenticated: () => true,
    error: () => null,
    selectedSite: () => ({
      id: 'site-example',
      name: 'Example Site',
      slug: 'example',
      domain: 'https://example.test',
      blogPath: '/articles',
      status: 'active' as const,
      templateKey: 'default-blog',
      themeConfig: '{}',
      deployProvider: 'netlify',
      deployConfig: '{}',
      previewDeployProvider: 'cloudflare',
      previewDeployConfig: '{}',
      aiConfig: '{}',
      storageConfig: '{}',
      updatedAt: '2026-05-23T00:00:00.000Z',
    }),
    selectedSiteId: () => 'site-example',
    selectedArticle: () => null,
    sites: () => [],
    authSession: () => null,
    authors: () => [],
    categories: () => [category],
    articles: () => [],
    tags: () => [tag],
    mediaAssets: () => [],
    builds: () => [],
    dashboardStats: () => [],
    landingSections: () => [],
    reportError: jasmine.createSpy('reportError'),
    clearError: jasmine.createSpy('clearError'),
    login: jasmine.createSpy('login').and.resolveTo(),
    logout: jasmine.createSpy('logout').and.resolveTo(),
    selectSite: jasmine.createSpy('selectSite').and.resolveTo(),
    selectArticle: jasmine.createSpy('selectArticle').and.resolveTo(),
    clearSelectedArticle: jasmine.createSpy('clearSelectedArticle'),
    createSite: jasmine.createSpy('createSite').and.resolveTo(),
    updateSelectedSite: jasmine.createSpy('updateSelectedSite').and.resolveTo(),
    createArticleDraft: jasmine.createSpy('createArticleDraft').and.resolveTo({ id: 'article-1' }),
    saveArticle: jasmine.createSpy('saveArticle').and.resolveTo({ id: 'article-1' }),
    saveCategory: jasmine.createSpy('saveCategory').and.resolveTo(category),
    deleteCategory: jasmine.createSpy('deleteCategory').and.resolveTo(),
    saveTag: jasmine.createSpy('saveTag').and.resolveTo(tag),
    deleteTag: jasmine.createSpy('deleteTag').and.resolveTo(),
    triggerBuild: jasmine.createSpy('triggerBuild').and.resolveTo(),
    toggleLandingSection: jasmine.createSpy('toggleLandingSection').and.resolveTo(),
    moveLandingSection: jasmine.createSpy('moveLandingSection').and.resolveTo(),
    uploadMedia: jasmine.createSpy('uploadMedia').and.resolveTo(),
    toggleFeatured: jasmine.createSpy('toggleFeatured').and.resolveTo(),
    setArticleStatus: jasmine.createSpy('setArticleStatus').and.resolveTo(),
    updateSelectedSiteSettings: jasmine.createSpy('updateSelectedSiteSettings').and.resolveTo(),
  });

  describe('categories', () => {
    let state: ReturnType<typeof createState>;

    beforeEach(async () => {
      state = createState();

      await TestBed.configureTestingModule({
        imports: [CommonModule, ReactiveFormsModule, RouterTestingModule, TaxonomyPageComponent],
        providers: [{ provide: WorkspaceStateService, useValue: state }],
      }).compileComponents();

      fixture = TestBed.createComponent(TaxonomyPageComponent);
      fixture.componentRef.setInput('kind', 'categories');
      fixture.detectChanges();
      await fixture.whenStable();
      fixture.detectChanges();
    });

    it('renders the category description field and uses category CRUD actions', async () => {
      expect(fixture.nativeElement.textContent).toContain('Categories');
      expect(fixture.nativeElement.textContent).toContain('New category');
      expect(fixture.nativeElement.textContent).toContain('Description');
      expect(fixture.nativeElement.querySelector('textarea[formcontrolname="description"]')).toBeTruthy();
      expect(fixture.nativeElement.querySelectorAll('.table-row').length).toBe(1);

      fixture.componentInstance.form.controls.name.setValue('Editorial');
      fixture.componentInstance.form.controls.description.setValue('Editorial content');
      await fixture.componentInstance.save();
      fixture.detectChanges();

      expect(state.clearError).toHaveBeenCalled();
      expect(state.saveCategory).toHaveBeenCalledWith({
        id: undefined,
        name: 'Editorial',
        description: 'Editorial content',
      });
      expect(fixture.nativeElement.textContent).toContain('Categories saved successfully.');

      fixture.componentInstance.edit(category);
      fixture.detectChanges();
      expect(fixture.nativeElement.querySelector('input[formcontrolname="id"]')?.value).toBe('category-1');

      spyOn(window, 'confirm').and.returnValue(true);
      await fixture.componentInstance.remove(category);
      fixture.detectChanges();

      expect(state.deleteCategory).toHaveBeenCalledWith('category-1');
      expect(fixture.nativeElement.textContent).toContain('Categories deleted successfully.');
    });
  });

  describe('tags', () => {
    let state: ReturnType<typeof createState>;

    beforeEach(async () => {
      state = createState();

      await TestBed.configureTestingModule({
        imports: [CommonModule, ReactiveFormsModule, RouterTestingModule, TaxonomyPageComponent],
        providers: [{ provide: WorkspaceStateService, useValue: state }],
      }).compileComponents();

      fixture = TestBed.createComponent(TaxonomyPageComponent);
      fixture.componentRef.setInput('kind', 'tags');
      fixture.detectChanges();
      await fixture.whenStable();
      fixture.detectChanges();
    });

    it('hides the description field and uses tag CRUD actions', async () => {
      expect(fixture.nativeElement.textContent).toContain('Tags');
      expect(fixture.nativeElement.textContent).not.toContain('New tag');
      expect(fixture.nativeElement.querySelector('textarea[formcontrolname="description"]')).toBeNull();
      expect(fixture.nativeElement.querySelectorAll('.table-row').length).toBe(1);

      fixture.componentInstance.form.controls.name.setValue('Launch');
      await fixture.componentInstance.save();
      fixture.detectChanges();

      expect(state.clearError).toHaveBeenCalled();
      expect(state.saveTag).toHaveBeenCalledWith({
        id: undefined,
        name: 'Launch',
      });
      expect(fixture.nativeElement.textContent).toContain('Tags saved successfully.');

      spyOn(window, 'confirm').and.returnValue(true);
      await fixture.componentInstance.remove(tag);
      fixture.detectChanges();

      expect(state.deleteTag).toHaveBeenCalledWith('tag-1');
      expect(fixture.nativeElement.textContent).toContain('Tags deleted successfully.');
    });
  });
});
