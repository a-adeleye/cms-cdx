import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { MediaLibraryPageComponent } from './media-library-page.component';

describe('MediaLibraryPageComponent', () => {
  let fixture: ComponentFixture<MediaLibraryPageComponent>;
  let router: Router;
  const asset = {
    id: 'media-1', siteId: 'site-1', fileName: 'cover.jpg', fileUrl: 'https://cdn.example/cover.jpg',
    mimeType: 'image/jpeg', sizeBytes: 2048, storageProvider: 'r2', storageKey: 'media/cover.jpg', altText: 'Article cover',
  };
  const state = {
    mediaAssets: () => [asset],
    articles: () => [],
    error: () => null,
    clearError: jasmine.createSpy('clearError'),
    reportError: jasmine.createSpy('reportError'),
    uploadMediaFile: jasmine.createSpy('uploadMediaFile').and.resolveTo(asset),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RouterTestingModule, MediaLibraryPageComponent],
      providers: [{ provide: WorkspaceStateService, useValue: state }],
    }).compileComponents();
    fixture = TestBed.createComponent(MediaLibraryPageComponent);
    router = TestBed.inject(Router);
    spyOn(router, 'navigate').and.resolveTo(true);
    fixture.detectChanges();
  });

  it('renders, filters, and opens media assets', () => {
    expect(fixture.nativeElement.textContent).toContain('cover.jpg');
    fixture.componentInstance.onSearch('missing');
    fixture.detectChanges();
    expect(fixture.nativeElement.textContent).toContain('No media found');
    fixture.componentInstance.openAsset('media-1');
    expect(router.navigate).toHaveBeenCalledWith(['/content/media', 'media-1']);
  });

  it('uploads the chosen file through workspace state', async () => {
    const file = new File(['image'], 'cover-image.jpg', { type: 'image/jpeg' });
    await fixture.componentInstance.uploadFile({ target: { files: [file], value: 'cover-image.jpg' } } as unknown as Event);
    expect(state.clearError).toHaveBeenCalled();
    expect(state.uploadMediaFile).toHaveBeenCalledWith(file, 'cover image');
  });
});
