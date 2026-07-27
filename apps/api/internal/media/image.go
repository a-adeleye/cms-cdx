package media

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"

	"github.com/gen2brain/webp"
)

const (
	maxImageDimension = 1600
	maxSourcePixels   = 40_000_000
	webPQuality       = 80
)

// OptimizedImage is the browser-ready representation stored for a static image.
type OptimizedImage struct {
	FileName string
	Contents []byte
	MimeType string
}

// OptimizeForUpload resizes supported static images and encodes them as WebP.
// GIF and ICO files remain unchanged to preserve animation and icon behaviour.
func OptimizeForUpload(fileName, mimeType string, contents []byte) (OptimizedImage, error) {
	original := OptimizedImage{FileName: fileName, Contents: contents, MimeType: mimeType}
	if mimeType == "image/gif" || mimeType == "image/x-icon" {
		return original, nil
	}
	imageValue, width, height, err := decodeImage(contents)
	if err != nil {
		return OptimizedImage{}, err
	}
	if mimeType == "image/webp" && width <= maxImageDimension && height <= maxImageDimension {
		return original, nil
	}
	return encodeWebP(fileName, resizeImage(imageValue, width, height))
}

func decodeImage(contents []byte) (image.Image, int, int, error) {
	config, _, err := image.DecodeConfig(bytes.NewReader(contents))
	if err != nil {
		return nil, 0, 0, errors.New("image cannot be decoded")
	}
	if config.Width <= 0 || config.Height <= 0 || config.Width > maxSourcePixels/config.Height {
		return nil, 0, 0, errors.New("image dimensions are too large")
	}
	imageValue, _, err := image.Decode(bytes.NewReader(contents))
	if err != nil {
		return nil, 0, 0, errors.New("image cannot be decoded")
	}
	bounds := imageValue.Bounds()
	return imageValue, bounds.Dx(), bounds.Dy(), nil
}

func resizeImage(source image.Image, width, height int) image.Image {
	newWidth, newHeight := resizedDimensions(width, height)
	if newWidth == width && newHeight == height {
		return source
	}
	resized := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight))
	bounds := source.Bounds()
	for y := 0; y < newHeight; y++ {
		sourceY := scaledCoordinate(y, height, newHeight)
		y0, y1 := int(sourceY), min(int(sourceY)+1, height-1)
		for x := 0; x < newWidth; x++ {
			sourceX := scaledCoordinate(x, width, newWidth)
			x0, x1 := int(sourceX), min(int(sourceX)+1, width-1)
			resized.SetNRGBA(x, y, interpolatedColor(
				colorAt(source, bounds.Min.X+x0, bounds.Min.Y+y0), colorAt(source, bounds.Min.X+x1, bounds.Min.Y+y0),
				colorAt(source, bounds.Min.X+x0, bounds.Min.Y+y1), colorAt(source, bounds.Min.X+x1, bounds.Min.Y+y1),
				sourceX-float64(x0), sourceY-float64(y0),
			))
		}
	}
	return resized
}

func scaledCoordinate(position, sourceSize, targetSize int) float64 {
	if targetSize <= 1 {
		return 0
	}
	return float64(position) * float64(sourceSize-1) / float64(targetSize-1)
}

func colorAt(source image.Image, x, y int) color.NRGBA {
	return color.NRGBAModel.Convert(source.At(x, y)).(color.NRGBA)
}

func interpolatedColor(topLeft, topRight, bottomLeft, bottomRight color.NRGBA, xWeight, yWeight float64) color.NRGBA {
	return color.NRGBA{
		R: interpolatedChannel(topLeft.R, topRight.R, bottomLeft.R, bottomRight.R, xWeight, yWeight),
		G: interpolatedChannel(topLeft.G, topRight.G, bottomLeft.G, bottomRight.G, xWeight, yWeight),
		B: interpolatedChannel(topLeft.B, topRight.B, bottomLeft.B, bottomRight.B, xWeight, yWeight),
		A: interpolatedChannel(topLeft.A, topRight.A, bottomLeft.A, bottomRight.A, xWeight, yWeight),
	}
}

func interpolatedChannel(topLeft, topRight, bottomLeft, bottomRight uint8, xWeight, yWeight float64) uint8 {
	top := float64(topLeft)*(1-xWeight) + float64(topRight)*xWeight
	bottom := float64(bottomLeft)*(1-xWeight) + float64(bottomRight)*xWeight
	return uint8(top*(1-yWeight) + bottom*yWeight)
}

func resizedDimensions(width, height int) (int, int) {
	if width <= maxImageDimension && height <= maxImageDimension {
		return width, height
	}
	if width >= height {
		return maxImageDimension, maxImageDimension * height / width
	}
	return maxImageDimension * width / height, maxImageDimension
}

func encodeWebP(fileName string, imageValue image.Image) (OptimizedImage, error) {
	var encoded bytes.Buffer
	if err := webp.Encode(&encoded, imageValue, webp.Options{Quality: webPQuality, Method: 4}); err != nil {
		return OptimizedImage{}, errors.New("image cannot be encoded as WebP")
	}
	return OptimizedImage{FileName: webPFileName(fileName), Contents: encoded.Bytes(), MimeType: "image/webp"}, nil
}

func webPFileName(fileName string) string {
	base := strings.TrimSuffix(filepath.Base(strings.TrimSpace(fileName)), filepath.Ext(fileName))
	if base == "" {
		base = "image"
	}
	return base + ".webp"
}
