import { DOCUMENT } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { convertToParamMap } from '@angular/router';
import { of } from 'rxjs';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';
import { MediaAssetDetailsPageComponent } from './media-asset-details-page.component';

describe('MediaAssetDetailsPageComponent', () => {
  let fixture: ComponentFixture<MediaAssetDetailsPageComponent>;
  let router: Router;
  const asset = {
    id: 'media-1', siteId: 'site-1', fileName: 'cover.jpg', fileUrl: 'https://cdn.example/cover.jpg',
    mimeType: 'image/jpeg', sizeBytes: 2048, storageProvider: 'r2', storageKey: 'site-1/media/cover.jpg', altText: 'Original cover',
  };
  const state = {
    mediaAssets: () => [asset],
    articles: () => [],
    updateMediaAltText: jasmine.createSpy('updateMediaAltText').and.resolveTo({ ...asset, altText: 'Updated cover' }),
    replaceMediaFile: jasmine.createSpy('replaceMediaFile').and.resolveTo({ ...asset, fileName: 'replacement.jpg' }),
    deleteMediaAsset: jasmine.createSpy('deleteMediaAsset').and.resolveTo(),
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RouterTestingModule, MediaAssetDetailsPageComponent],
      providers: [
        { provide: ActivatedRoute, useValue: { paramMap: of(convertToParamMap({ assetId: 'media-1' })) } },
        { provide: WorkspaceStateService, useValue: state },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(MediaAssetDetailsPageComponent);
    router = TestBed.inject(Router);
    spyOn(router, 'navigate').and.resolveTo(true);
    fixture.detectChanges();
  });

  it('saves changed alt text for the selected asset', async () => {
    fixture.componentInstance.altTextForm.controls.altText.setValue('Updated cover');

    await fixture.componentInstance.saveAltText();

    expect(state.updateMediaAltText).toHaveBeenCalledWith('media-1', 'Updated cover');
  });

  it('replaces the selected asset file', async () => {
    const file = new File(['image'], 'replacement.jpg', { type: 'image/jpeg' });

    await fixture.componentInstance.replaceFile({ target: { files: [file], value: 'replacement.jpg' } } as unknown as Event);

    expect(state.replaceMediaFile).toHaveBeenCalledWith('media-1', file, 'Original cover');
  });

  it('deletes the selected asset only after confirmation', async () => {
    spyOn(TestBed.inject(DOCUMENT).defaultView!, 'confirm').and.returnValue(true);

    await fixture.componentInstance.deleteAsset();

    expect(state.deleteMediaAsset).toHaveBeenCalledWith('media-1');
    expect(router.navigate).toHaveBeenCalledWith(['/content/media']);
  });
});
