package media

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"testing"
)

func TestOptimizeForUploadResizesAndEncodesWebP(t *testing.T) {
	input := makePNG(t, 1800, 900)
	optimized, err := OptimizeForUpload("cover.png", "image/png", input)
	if err != nil {
		t.Fatalf("OptimizeForUpload returned error: %v", err)
	}
	if optimized.FileName != "cover.webp" || optimized.MimeType != "image/webp" {
		t.Fatalf("unexpected optimized metadata: %#v", optimized)
	}
	decoded, format, err := image.Decode(bytes.NewReader(optimized.Contents))
	if err != nil {
		t.Fatalf("decode optimized image: %v", err)
	}
	if format != "webp" || decoded.Bounds().Dx() != 1600 || decoded.Bounds().Dy() != 800 {
		t.Fatalf("unexpected optimized image: format=%q bounds=%v", format, decoded.Bounds())
	}
}

func TestOptimizeForUploadPreservesGIF(t *testing.T) {
	var input bytes.Buffer
	if err := gif.Encode(&input, image.NewPaletted(image.Rect(0, 0, 1, 1), color.Palette{color.Black}), nil); err != nil {
		t.Fatalf("encode GIF: %v", err)
	}
	optimized, err := OptimizeForUpload("animation.gif", "image/gif", input.Bytes())
	if err != nil {
		t.Fatalf("OptimizeForUpload returned error: %v", err)
	}
	if optimized.FileName != "animation.gif" || optimized.MimeType != "image/gif" || !bytes.Equal(optimized.Contents, input.Bytes()) {
		t.Fatalf("expected GIF to remain unchanged, got %#v", optimized)
	}
}

func makePNG(t *testing.T, width, height int) []byte {
	t.Helper()
	imageValue := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			imageValue.SetRGBA(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: uint8((x + y) % 256), A: 255})
		}
	}
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, imageValue); err != nil {
		t.Fatalf("encode PNG: %v", err)
	}
	return encoded.Bytes()
}
