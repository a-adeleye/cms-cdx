import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import {
  ArticleRecord,
  AuthSession,
  AuthorRecord,
  BuildRecord,
  CategoryRecord,
  LandingSectionRecord,
  MediaAssetRecord,
  SiteRecord,
  TemplateRecord,
  TagRecord,
} from './pages.model';
import { AuthTokenService } from './auth-token.service';

interface ApiUser {
  id: string;
  email: string;
  fullName: string;
  role: 'admin' | 'editor';
}

interface LoginResponse {
  token: string;
  user: ApiUser;
}

interface WorkspaceResponse {
  user: ApiUser;
  selectedSiteId: string;
  selectedArticleId: string | null;
  sites: SiteRecord[];
  templates: TemplateRecord[];
  landingSections: LandingSectionRecord[];
  articles: ArticleRecord[];
  authors: Array<{ id: string; siteId: string; name: string; slug: string; bio: string }>;
  categories: Array<{ id: string; siteId: string; name: string; slug: string; description: string }>;
  tags: Array<{ id: string; siteId: string; name: string; slug: string }>;
  mediaAssets: MediaAssetRecord[];
  builds: BuildRecord[];
}

interface ItemsResponse<T> {
  items: T[];
}

interface ArticleUpsertPayload extends Partial<ArticleRecord> {
  id?: string;
  title: string;
  slug: string;
  excerpt: string;
  contentMarkdown: string;
  coverImageUrl: string;
  seoTitle: string;
  seoDescription: string;
  canonicalUrl: string;
  authorId: string;
  categoryId: string;
  tagIds: string[];
  isFeatured: boolean;
  status: ArticleRecord['status'];
}

interface SiteUpsertPayload {
  name: string;
  slug: string;
  domain: string;
  blogPath: string;
  status: SiteRecord['status'];
  templateKey: string;
  themeConfig: string;
  deployProvider: string;
  deployConfig: string;
  previewDeployProvider: string;
  previewDeployConfig: string;
  aiConfig: string;
  storageConfig: string;
}

interface TemplateUpsertPayload {
  name: string;
  slug: string;
}

interface LandingOrderPayload {
  sectionIds: string[];
}

interface LandingUpdatePayload {
  isEnabled?: boolean;
  displayOrder?: number;
}

interface BuildCreatePayload {
  buildType: 'preview' | 'published';
  articleIds?: string[];
}

interface MediaCreatePayload {
  fileName: string;
  fileUrl: string;
  mimeType: string;
  sizeBytes: number;
  storageProvider: string;
  storageKey: string;
  altText: string;
}

interface CategoryUpsertPayload {
  name: string;
  description: string;
}

interface AuthorUpsertPayload {
  name: string;
  bio: string;
}

interface TagUpsertPayload {
  name: string;
}

@Injectable({
  providedIn: 'root',
})
export class AdminApiService {
  private readonly http = inject(HttpClient);
  private readonly tokenStore = inject(AuthTokenService);
  private readonly baseUrl = '/api/v1';

  async login(email: string, password: string): Promise<LoginResponse> {
    return firstValueFrom(this.http.post<LoginResponse>(`${this.baseUrl}/auth/login`, { email, password }));
  }

  async logout(): Promise<void> {
    await firstValueFrom(this.http.post(`${this.baseUrl}/auth/logout`, {}, { headers: this.headers() }));
  }

  async me(): Promise<{ user: ApiUser }> {
    return firstValueFrom(this.http.get<{ user: ApiUser }>(`${this.baseUrl}/auth/me`, { headers: this.headers() }));
  }

  async loadWorkspace(siteId?: string): Promise<WorkspaceResponse> {
    const url = siteId ? `${this.baseUrl}/workspace?siteId=${encodeURIComponent(siteId)}` : `${this.baseUrl}/workspace`;
    return firstValueFrom(this.http.get<WorkspaceResponse>(url, { headers: this.headers() }));
  }

  async listTemplates(): Promise<ItemsResponse<TemplateRecord>> {
    return firstValueFrom(this.http.get<ItemsResponse<TemplateRecord>>(`${this.baseUrl}/templates`, { headers: this.headers() }));
  }

  async createTemplate(payload: TemplateUpsertPayload): Promise<TemplateRecord> {
    return firstValueFrom(this.http.post<TemplateRecord>(`${this.baseUrl}/templates`, payload, { headers: this.headers() }));
  }

  async createSite(payload: SiteUpsertPayload): Promise<SiteRecord> {
    return firstValueFrom(this.http.post<SiteRecord>(`${this.baseUrl}/sites`, payload, { headers: this.headers() }));
  }

  async updateSite(siteId: string, payload: SiteUpsertPayload): Promise<SiteRecord> {
    return firstValueFrom(this.http.patch<SiteRecord>(`${this.baseUrl}/sites/${siteId}`, payload, { headers: this.headers() }));
  }

  async listSites(): Promise<ItemsResponse<SiteRecord>> {
    return firstValueFrom(this.http.get<ItemsResponse<SiteRecord>>(`${this.baseUrl}/sites`, { headers: this.headers() }));
  }

  async listArticles(siteId: string): Promise<ItemsResponse<ArticleRecord>> {
    return firstValueFrom(this.http.get<ItemsResponse<ArticleRecord>>(`${this.baseUrl}/sites/${siteId}/articles`, { headers: this.headers() }));
  }

  async createAuthor(siteId: string, payload: AuthorUpsertPayload): Promise<AuthorRecord> {
    try {
      return await firstValueFrom(
        this.http.post<AuthorRecord>(`${this.baseUrl}/sites/${siteId}/authors`, payload, { headers: this.headers() }),
      );
    } catch (error) {
      throw this.toError(error);
    }
  }

  async updateAuthor(siteId: string, authorId: string, payload: AuthorUpsertPayload): Promise<AuthorRecord> {
    try {
      return await firstValueFrom(
        this.http.patch<AuthorRecord>(`${this.baseUrl}/sites/${siteId}/authors/${authorId}`, payload, { headers: this.headers() }),
      );
    } catch (error) {
      throw this.toError(error);
    }
  }

  async deleteAuthor(siteId: string, authorId: string): Promise<void> {
    try {
      await firstValueFrom(this.http.delete(`${this.baseUrl}/sites/${siteId}/authors/${authorId}`, { headers: this.headers() }));
    } catch (error) {
      throw this.toError(error);
    }
  }

  async createCategory(siteId: string, payload: CategoryUpsertPayload): Promise<CategoryRecord> {
    try {
      return await firstValueFrom(
        this.http.post<CategoryRecord>(`${this.baseUrl}/sites/${siteId}/categories`, payload, { headers: this.headers() }),
      );
    } catch (error) {
      throw this.toError(error);
    }
  }

  async updateCategory(
    siteId: string,
    categoryId: string,
    payload: CategoryUpsertPayload,
  ): Promise<CategoryRecord> {
    try {
      return await firstValueFrom(
        this.http.patch<CategoryRecord>(`${this.baseUrl}/sites/${siteId}/categories/${categoryId}`, payload, { headers: this.headers() }),
      );
    } catch (error) {
      throw this.toError(error);
    }
  }

  async deleteCategory(siteId: string, categoryId: string): Promise<void> {
    try {
      await firstValueFrom(this.http.delete(`${this.baseUrl}/sites/${siteId}/categories/${categoryId}`, { headers: this.headers() }));
    } catch (error) {
      throw this.toError(error);
    }
  }

  async createTag(siteId: string, payload: TagUpsertPayload): Promise<TagRecord> {
    try {
      return await firstValueFrom(
        this.http.post<TagRecord>(`${this.baseUrl}/sites/${siteId}/tags`, payload, { headers: this.headers() }),
      );
    } catch (error) {
      throw this.toError(error);
    }
  }

  async updateTag(siteId: string, tagId: string, payload: TagUpsertPayload): Promise<TagRecord> {
    try {
      return await firstValueFrom(
        this.http.patch<TagRecord>(`${this.baseUrl}/sites/${siteId}/tags/${tagId}`, payload, { headers: this.headers() }),
      );
    } catch (error) {
      throw this.toError(error);
    }
  }

  async deleteTag(siteId: string, tagId: string): Promise<void> {
    try {
      await firstValueFrom(this.http.delete(`${this.baseUrl}/sites/${siteId}/tags/${tagId}`, { headers: this.headers() }));
    } catch (error) {
      throw this.toError(error);
    }
  }

  async upsertArticle(siteId: string, payload: ArticleUpsertPayload): Promise<ArticleRecord> {
    try {
      // Temporarily disable update semantics while we trace the editor/article ID flow.
      // Always create a new article so PATCH cannot be triggered from the client.
      const createPayload = { ...payload };
      delete createPayload.id;
      return await firstValueFrom(this.http.post<ArticleRecord>(`${this.baseUrl}/sites/${siteId}/articles`, createPayload, { headers: this.headers() }));
    } catch (error) {
      throw this.toError(error);
    }
  }

  async updateArticle(articleId: string, payload: ArticleUpsertPayload): Promise<ArticleRecord> {
    try {
      return await firstValueFrom(
        this.http.patch<ArticleRecord>(`${this.baseUrl}/articles/${articleId}`, payload, { headers: this.headers() }),
      );
    } catch (error) {
      throw this.toError(error);
    }
  }

  async deleteArticle(articleId: string): Promise<void> {
    try {
      await firstValueFrom(this.http.delete(`${this.baseUrl}/articles/${articleId}`, { headers: this.headers() }));
    } catch (error) {
      throw this.toError(error);
    }
  }

  async listLandingSections(siteId: string): Promise<ItemsResponse<LandingSectionRecord>> {
    return firstValueFrom(this.http.get<ItemsResponse<LandingSectionRecord>>(`${this.baseUrl}/sites/${siteId}/landing-sections`, { headers: this.headers() }));
  }

  async updateLandingSection(siteId: string, sectionId: string, payload: LandingUpdatePayload): Promise<LandingSectionRecord> {
    return firstValueFrom(
      this.http.patch<LandingSectionRecord>(`${this.baseUrl}/sites/${siteId}/landing-sections/${sectionId}`, payload, {
        headers: this.headers(),
      }),
    );
  }

  async reorderLandingSections(siteId: string, sectionIds: string[]): Promise<ItemsResponse<LandingSectionRecord>> {
    return firstValueFrom(
      this.http.put<ItemsResponse<LandingSectionRecord>>(
        `${this.baseUrl}/sites/${siteId}/landing-sections`,
        {
          sectionIds,
        } satisfies LandingOrderPayload,
        { headers: this.headers() },
      ),
    );
  }

  async createBuild(siteId: string, buildType: 'preview' | 'published', articleIds: string[] = []): Promise<BuildRecord> {
    return firstValueFrom(
      this.http.post<BuildRecord>(
        `${this.baseUrl}/sites/${siteId}/builds`,
        { buildType, articleIds } satisfies BuildCreatePayload,
        {
          headers: this.headers(),
        },
      ),
    );
  }

  async createMediaAsset(siteId: string, payload: MediaCreatePayload): Promise<MediaAssetRecord> {
    return firstValueFrom(this.http.post<MediaAssetRecord>(`${this.baseUrl}/sites/${siteId}/media`, payload, { headers: this.headers() }));
  }

  async uploadMediaFile(siteId: string, file: File, altText: string): Promise<MediaAssetRecord> {
    const formData = new FormData();
    formData.append('file', file, file.name);
    formData.append('altText', altText);
    return firstValueFrom(this.http.post<MediaAssetRecord>(`${this.baseUrl}/sites/${siteId}/media`, formData, { headers: this.headers() }));
  }

  private headers(): HttpHeaders {
    const token = this.tokenStore.getToken();
    if (!token) {
      return new HttpHeaders();
    }

    return new HttpHeaders({ Authorization: `Bearer ${token}` });
  }

  private toError(error: unknown): Error {
    if (error instanceof HttpErrorResponse) {
      const message = this.extractApiErrorMessage(error);
      return new Error(message);
    }

    if (error instanceof Error) {
      return error;
    }

    return new Error('Request failed');
  }

  private extractApiErrorMessage(error: HttpErrorResponse): string {
    const body = error.error;
    if (typeof body === 'string' && body.trim()) {
      return body;
    }

    if (body && typeof body === 'object' && 'error' in body && typeof body.error === 'string' && body.error.trim()) {
      return body.error;
    }

    return error.message || 'Request failed';
  }
}
