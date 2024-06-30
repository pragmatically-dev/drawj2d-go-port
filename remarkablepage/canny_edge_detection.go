package remarkablepage

import (
	"image"
	"image/color"
	"math"
	"sync"
)

var sobelX = [][]float64{
	{-1, 0, 1},
	{-2, 0, 2},
	{-1, 0, 1},
}

var sobelY = [][]float64{
	{-1, -2, -1},
	{0, 0, 0},
	{1, 2, 1},
}

// Gaussian kernel for smoothing
var gaussianKernel = [][]float64{
	{1, 4, 7, 4, 1},
	{4, 16, 26, 16, 4},
	{7, 26, 41, 26, 7},
	{4, 16, 26, 16, 4},
	{1, 4, 7, 4, 1},
}

// Normalize the Gaussian kernel
func normalizeKernel(kernel [][]float64) {
	sum := 0.0
	for _, row := range kernel {
		for _, val := range row {
			sum += val
		}
	}
	for y := range kernel {
		for x := range kernel[y] {
			kernel[y][x] /= sum
		}
	}
}

func init() {
	normalizeKernel(gaussianKernel)
}

// ConvolveGray applies a convolution matrix (kernel) to a grayscale image.
func ConvolveGrayCanny(img *image.Gray, kernel [][]float64) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	kernelWidth := len(kernel[0])
	kernelHeight := len(kernel)

	for y := bounds.Min.Y + kernelHeight/2; y < bounds.Max.Y-kernelHeight/2; y++ {
		for x := bounds.Min.X + kernelWidth/2; x < bounds.Max.X-kernelWidth/2; x++ {
			var sum float64
			for ky := 0; ky < kernelHeight; ky++ {
				for kx := 0; kx < kernelWidth; kx++ {
					pixel := img.GrayAt(x+kx-kernelWidth/2, y+ky-kernelHeight/2).Y
					sum += float64(pixel) * kernel[ky][kx]
				}
			}
			gray.SetGray(x, y, color.Gray{uint8(sum)})
		}
	}
	return gray
}

// CannyEdgeDetection performs Canny edge detection on a grayscale image.
func CannyEdgeDetection(img *image.Gray, lowThreshold, highThreshold float64) *image.Gray {
	smoothed := ConvolveGrayCanny(img, gaussianKernel)
	gradientMagnitude, gradientDirection := computeGradient(smoothed)
	suppressed := nonMaximumSuppression(gradientMagnitude, gradientDirection)
	thresholded := doubleThreshold(suppressed, lowThreshold, highThreshold)
	edges := edgeTrackingByHysteresis(thresholded, lowThreshold, highThreshold)
	return edges
}

func computeGradient(img *image.Gray) (*image.Gray, *image.Gray) {
	bounds := img.Bounds()
	gradientMagnitude := image.NewGray(bounds)
	gradientDirection := image.NewGray(bounds)

	for y := 1; y < bounds.Max.Y-1; y++ {
		for x := 1; x < bounds.Max.X-1; x++ {
			var gx, gy float64
			for ky := 0; ky < 3; ky++ {
				for kx := 0; kx < 3; kx++ {
					pixel := img.GrayAt(x+kx-1, y+ky-1).Y
					gx += float64(pixel) * sobelX[ky][kx]
					gy += float64(pixel) * sobelY[ky][kx]
				}
			}
			magnitude := math.Sqrt(gx*gx + gy*gy)
			gradientMagnitude.SetGray(x, y, color.Gray{uint8(magnitude)})

			// Angle in radians, converted to degrees
			angle := math.Atan2(gy, gx) * (180.0 / math.Pi)
			if angle < 0 {
				angle += 180
			}
			gradientDirection.SetGray(x, y, color.Gray{uint8(angle)})
		}
	}
	return gradientMagnitude, gradientDirection
}

func nonMaximumSuppression(magnitude, direction *image.Gray) *image.Gray {
	bounds := magnitude.Bounds()
	suppressed := image.NewGray(bounds)

	for y := 1; y < bounds.Max.Y-1; y++ {
		for x := 1; x < bounds.Max.X-1; x++ {
			angle := float64(direction.GrayAt(x, y).Y)
			q := uint8(255)
			r := uint8(255)

			// Angle 0
			if (0 <= angle && angle < 22.5) || (157.5 <= angle && angle <= 180) {
				q = magnitude.GrayAt(x+1, y).Y
				r = magnitude.GrayAt(x-1, y).Y
			} else if 22.5 <= angle && angle < 67.5 { // Angle 45
				q = magnitude.GrayAt(x+1, y-1).Y
				r = magnitude.GrayAt(x-1, y+1).Y
			} else if 67.5 <= angle && angle < 112.5 { // Angle 90
				q = magnitude.GrayAt(x, y+1).Y
				r = magnitude.GrayAt(x, y-1).Y
			} else if 112.5 <= angle && angle < 157.5 { // Angle 135
				q = magnitude.GrayAt(x-1, y-1).Y
				r = magnitude.GrayAt(x+1, y+1).Y
			}

			if magnitude.GrayAt(x, y).Y >= q && magnitude.GrayAt(x, y).Y >= r {
				suppressed.SetGray(x, y, magnitude.GrayAt(x, y))
			} else {
				suppressed.SetGray(x, y, color.Gray{0})
			}
		}
	}
	return suppressed
}

func doubleThreshold(img *image.Gray, lowThreshold, highThreshold float64) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.GrayAt(x, y).Y
			if float64(pixel) >= highThreshold {
				result.SetGray(x, y, color.Gray{255})
			} else if float64(pixel) >= lowThreshold {
				result.SetGray(x, y, color.Gray{128})
			} else {
				result.SetGray(x, y, color.Gray{0})
			}
		}
	}
	return result
}

func edgeTrackingByHysteresis(img *image.Gray, lowThreshold, highThreshold float64) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)
	var wg sync.WaitGroup

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			if img.GrayAt(x, y).Y == 128 {
				wg.Add(1)
				go func(x, y int) {
					defer wg.Done()
					if isStrongEdge(img, x, y) {
						traceEdge(img, result, x, y)
					}
				}(x, y)
			}
		}
	}

	wg.Wait()
	return result
}

func isStrongEdge(img *image.Gray, x, y int) bool {
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			if img.GrayAt(x+i, y+j).Y == 255 {
				return true
			}
		}
	}
	return false
}

func traceEdge(src, dst *image.Gray, x, y int) {
	stack := []image.Point{{X: x, Y: y}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if src.GrayAt(p.X, p.Y).Y == 128 {
			dst.SetGray(p.X, p.Y, color.Gray{255})
			src.SetGray(p.X, p.Y, color.Gray{0})

			for j := -1; j <= 1; j++ {
				for i := -1; i <= 1; i++ {
					if src.GrayAt(p.X+i, p.Y+j).Y == 128 {
						stack = append(stack, image.Point{X: p.X + i, Y: p.Y + j})
					}
				}
			}
		}
	}
}
