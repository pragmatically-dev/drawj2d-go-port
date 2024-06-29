package remarkablepage

import (
	"image"
	"image/color"
	"sync"
)

// MaxUint8 - maximum value which can be held in an uint8
const MaxUint8 = ^uint8(0)

// MinUint8 - minimum value which can be held in an uint8
const MinUint8 = 0

// MaxUint16 - maximum value which can be held in an uint16
const MaxUint16 = ^uint16(0)

// MinUint16 - minimum value which can be held in an uint8
const MinUint16 = 0

// ConvolveGray applies a convolution matrix (kernel) to a grayscale image.
// Example of usage:
//
//	res, err := convolution.ConvolveGray(img, kernel, {1, 1}, BorderReflect)
//
// Note: the anchor represents a point inside the area of the kernel. After every step of the convolution, the position
// specified by the anchor point gets updated on the result image.
func ConvolveGray(img *image.Gray, kernel *Kernel, anchor image.Point, border CBorder) (*image.Gray, error) {
	kernelSize := kernel.Size()
	padded, err := PaddingGray(img, kernelSize, anchor, border)
	if err != nil {
		return nil, err
	}
	originalSize := img.Bounds().Size()
	resultImage := image.NewGray(img.Bounds())

	// Parallel processing using multiple goroutines
	var wg sync.WaitGroup
	wg.Add(originalSize.Y)

	for y := 0; y < originalSize.Y; y++ {
		go func(y int) {
			defer wg.Done()
			for x := 0; x < originalSize.X; x++ {
				sum := float64(0)
				for ky := 0; ky < kernelSize.Y; ky++ {
					for kx := 0; kx < kernelSize.X; kx++ {
						pixel := padded.GrayAt(x+kx, y+ky)
						kE := kernel.At(kx, ky)
						sum += float64(pixel.Y) * kE
					}
				}
				sum = ClampF64(sum, MinUint8, float64(MaxUint8))
				resultImage.Set(x, y, color.Gray{uint8(sum)})
			}
		}(y)
	}

	wg.Wait()
	return resultImage, nil
}

// ConvolveRGBA applies a convolution matrix (kernel) to an RGBA image.
// Example of usage:
//
//	res, err := convolution.ConvolveRGBA(img, kernel, {1, 1}, BorderReflect)
//
// Note: the anchor represents a point inside the area of the kernel. After every step of the convolution, the position
// specified by the anchor point gets updated on the result image.
func ConvolveRGBA(img *image.RGBA, kernel *Kernel, anchor image.Point, border CBorder) (*image.RGBA, error) {
	kernelSize := kernel.Size()
	padded, err := PaddingRGBA(img, kernelSize, anchor, border)
	if err != nil {
		return nil, err
	}
	originalSize := img.Bounds().Size()
	resultImage := image.NewRGBA(img.Bounds())

	// Parallel processing using multiple goroutines
	var wg sync.WaitGroup
	wg.Add(originalSize.Y)

	for y := 0; y < originalSize.Y; y++ {
		go func(y int) {
			defer wg.Done()
			for x := 0; x < originalSize.X; x++ {
				sumR, sumG, sumB := 0.0, 0.0, 0.0
				for kx := 0; kx < kernelSize.X; kx++ {
					for ky := 0; ky < kernelSize.Y; ky++ {
						pixel := padded.RGBAAt(x+kx, y+ky)
						sumR += float64(pixel.R) * kernel.At(kx, ky)
						sumG += float64(pixel.G) * kernel.At(kx, ky)
						sumB += float64(pixel.B) * kernel.At(kx, ky)
					}
				}
				sumR = ClampF64(sumR, MinUint8, float64(MaxUint8))
				sumG = ClampF64(sumG, MinUint8, float64(MaxUint8))
				sumB = ClampF64(sumB, MinUint8, float64(MaxUint8))
				rgba := img.RGBAAt(x, y)
				resultImage.Set(x, y, color.RGBA{uint8(sumR), uint8(sumG), uint8(sumB), rgba.A})
			}
		}(y)
	}

	wg.Wait()
	return resultImage, nil
}

func ClampF64(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
