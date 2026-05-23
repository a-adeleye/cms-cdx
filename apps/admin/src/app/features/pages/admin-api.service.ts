import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import {
  ArticleRecord,
  AuthSession,
  BuildRecord,
  LandingSectionRecord,
  MediaAssetRecord,
  SiteRecord,
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
  aiConfig: string;
  storageConfig: string;
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

  async upsertArticle(siteId: string, payload: ArticleUpsertPayload): Promise<ArticleRecord> {
    try {
      if (payload.id) {
        return await firstValueFrom(this.http.patch<ArticleRecord>(`${this.baseUrl}/articles/${payload.id}`, payload, { headers: this.headers() }));
      }

      return await firstValueFrom(this.http.post<ArticleRecord>(`${this.baseUrl}/sites/${siteId}/articles`, payload, { headers: this.headers() }));
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

  async createBuild(siteId: string, buildType: 'preview' | 'published'): Promise<BuildRecord> {
    return firstValueFrom(
      this.http.post<BuildRecord>(`${this.baseUrl}/sites/${siteId}/builds`, { buildType } satisfies BuildCreatePayload, {
        headers: this.headers(),
      }),
    );
  }

  async createMediaAsset(siteId: string, payload: MediaCreatePayload): Promise<MediaAssetRecord> {
    return firstValueFrom(this.http.post<MediaAssetRecord>(`${this.baseUrl}/sites/${siteId}/media`, payload, { headers: this.headers() }));
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
