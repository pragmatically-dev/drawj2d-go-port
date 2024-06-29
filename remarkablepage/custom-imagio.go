package remarkablepage

import (
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

// DecodeAndConfig decodes an image from the given path and returns its configuration (width, height) without fully decoding it.
func DecodeAndConfig(path string) (image.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return image.Config{}, err
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	return config, err
}

// DecodeToGray decodes an image to grayscale from the given path.
func DecodeToGray(path string) (*image.Gray, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, img, bounds.Min, draw.Src)
	return gray, nil
}

// DecodeToGray16 decodes an image to grayscale16 from the given path.
func DecodeToGray16(path string) (*image.Gray16, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	gray16 := image.NewGray16(bounds)
	draw.Draw(gray16, bounds, img, bounds.Min, draw.Src)
	return gray16, nil
}

// DecodeToRGBA decodes an image to RGBA from the given path.
func DecodeToRGBA(path string) (*image.RGBA, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba, nil
}

// DecodeToRGBA64 decodes an image to RGBA64 from the given path.
func DecodeToRGBA64(path string) (*image.RGBA64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	rgba64 := image.NewRGBA64(bounds)
	draw.Draw(rgba64, bounds, img, bounds.Min, draw.Src)
	return rgba64, nil
}

// Encode encodes an image to the given format and writes it to the specified path.
func Encode(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	extension := filepath.Ext(path)
	switch extension {
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, nil)
	case ".png":
		return png.Encode(file, img)
	default:
		return errors.New("unsupported extension")
	}
}

// Imread reads an image from the given path and returns it in the desired format.
// Supported formats: "gray", "gray16", "rgba", "rgba64"
func Imread(path string, format string) (interface{}, error) {
	switch format {
	case "gray":
		return DecodeToGray(path)
	case "gray16":
		return DecodeToGray16(path)
	case "rgba":
		return DecodeToRGBA(path)
	case "rgba64":
		return DecodeToRGBA64(path)
	default:
		return nil, errors.New("unsupported format")
	}
}

// Imwrite saves the image to the specified path.
func Imwrite(img image.Image, path string) error {
	return Encode(img, path)
}
