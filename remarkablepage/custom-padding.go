package remarkablepage

import (
	"errors"
	"image"
	"image/color"
	"sync"
)

// CBorder is an enum type for supported padding types
type CBorder int

const (
	// CBorderConstant - xxxabcdefghxxx - where x is a black ( color.Gray{0} ) pixel
	CBorderConstant CBorder = iota
	// CBorderReplicate - aaaabcdefghhhh - replicates the nearest pixel
	CBorderReplicate
	// CBorderReflect - cbabcdefgfed - reflects the nearest pixel group
	CBorderReflect
)

// Paddings struct holds the padding sizes for each padding
type Paddings struct {
	// PaddingLeft is the size of the left padding
	PaddingLeft int
	// PaddingRight is the size of the right padding
	PaddingRight int
	// PaddingTop is the size of the top padding
	PaddingTop int
	// PaddingBottom is the size of the bottom padding
	PaddingBottom int
}

func topPaddingReplicate(img image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for x := p.PaddingLeft; x < originalSize.X+p.PaddingLeft; x++ {
		firstPixel := img.At(x-p.PaddingLeft, p.PaddingTop)
		for y := 0; y < p.PaddingTop; y++ {
			setPixel(x, y, firstPixel)
		}
	}
}

func bottomPaddingReplicate(img image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for x := p.PaddingLeft; x < originalSize.X+p.PaddingLeft; x++ {
		lastPixel := img.At(x-p.PaddingLeft, originalSize.Y-1)
		for y := p.PaddingTop + originalSize.Y; y < originalSize.Y+p.PaddingTop+p.PaddingBottom; y++ {
			setPixel(x, y, lastPixel)
		}
	}
}

func leftPaddingReplicate(img image.Image, padded image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for y := 0; y < originalSize.Y+p.PaddingBottom+p.PaddingTop; y++ {
		firstPixel := padded.At(p.PaddingLeft, y)
		for x := 0; x < p.PaddingLeft; x++ {
			setPixel(x, y, firstPixel)
		}
	}
}

func rightPaddingReplicate(img image.Image, padded image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for y := 0; y < originalSize.Y+p.PaddingBottom+p.PaddingTop; y++ {
		lastPixel := padded.At(originalSize.X+p.PaddingLeft-1, y)
		for x := originalSize.X + p.PaddingLeft; x < originalSize.X+p.PaddingLeft+p.PaddingRight; x++ {
			setPixel(x, y, lastPixel)
		}
	}
}

func topPaddingReflect(img image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for x := p.PaddingLeft; x < originalSize.X+p.PaddingLeft; x++ {
		for y := 0; y < p.PaddingTop; y++ {
			pixel := img.At(x-p.PaddingLeft, p.PaddingTop-y)
			setPixel(x, y, pixel)
		}
	}
}

func bottomPaddingReflect(img image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for x := p.PaddingLeft; x < originalSize.X+p.PaddingLeft; x++ {
		for y := p.PaddingTop + originalSize.Y; y < originalSize.Y+p.PaddingTop+p.PaddingBottom; y++ {
			pixel := img.At(x-p.PaddingLeft, originalSize.Y-(y-p.PaddingTop-originalSize.Y)-2)
			setPixel(x, y, pixel)
		}
	}
}

func leftPaddingReflect(img image.Image, padded image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for y := 0; y < originalSize.Y+p.PaddingBottom+p.PaddingTop; y++ {
		for x := 0; x < p.PaddingLeft; x++ {
			pixel := padded.At(2*p.PaddingLeft-x, y)
			setPixel(x, y, pixel)
		}
	}
}

func rightPaddingReflect(img image.Image, padded image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for y := 0; y < originalSize.Y+p.PaddingBottom+p.PaddingTop; y++ {
		for x := originalSize.X + p.PaddingLeft; x < originalSize.X+p.PaddingLeft+p.PaddingRight; x++ {
			pixel := padded.At(originalSize.X+p.PaddingLeft-(x-originalSize.X-p.PaddingLeft)-2, y)
			setPixel(x, y, pixel)
		}
	}
}

// ParallelForEachPixel applies a function to each pixel in the image in parallel.
func ParallelForEachPixel(size image.Point, fn func(x, y int)) {
	var wg sync.WaitGroup
	wg.Add(size.Y)

	for j := 0; j < size.Y; j++ {
		go func(j int) {
			defer wg.Done()
			for i := 0; i < size.X; i++ {
				fn(i, j)
			}
		}(j)
	}

	wg.Wait()
}

// PaddingGray appends padding to a given grayscale image. The size of the padding is calculated from the kernel size
// and the anchor point. Supported CBorder types are: CBorderConstant, CBorderReplicate, CBorderReflect.
// Example of usage:
//
//	res, err := padding.PaddingGray(img, {5, 5}, {1, 1}, CBorderReflect)
//
// Note: this will add a 1px padding for the top and left CBorders of the image and a 3px padding for the bottom and
// right CBorders of the image.
func PaddingGray(img *image.Gray, kernelSize image.Point, anchor image.Point, CBorder CBorder) (*image.Gray, error) {
	originalSize := img.Bounds().Size()
	p, err := calculatePaddings(kernelSize, anchor)
	if err != nil {
		return nil, err
	}
	rect := getRectangleFromPaddings(p, originalSize)
	padded := image.NewGray(rect)

	ParallelForEachPixel(originalSize, func(x, y int) {
		padded.Set(x+p.PaddingLeft, y+p.PaddingTop, img.GrayAt(x, y))
	})

	switch CBorder {
	case CBorderConstant:
		// do nothing
	case CBorderReplicate:
		topPaddingReplicate(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		bottomPaddingReplicate(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		leftPaddingReplicate(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		rightPaddingReplicate(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
	case CBorderReflect:
		topPaddingReflect(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		bottomPaddingReflect(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		leftPaddingReflect(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		rightPaddingReflect(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
	default:
		return nil, errors.New("unknown CBorder type")
	}
	return padded, nil
}

// PaddingRGBA appends padding to a given RGBA image. The size of the padding is calculated from the kernel size
// and the anchor point. Supported CBorder types are: CBorderConstant, CBorderReplicate, CBorderReflect.
// Example of usage:
//
//	res, err := padding.PaddingRGBA(img, {5, 5}, {1, 1}, CBorderReflect)
//
// Note: this will add a 1px padding for the top and left CBorders of the image and a 3px padding for the bottom and
// right CBorders of the image.
func PaddingRGBA(img *image.RGBA, kernelSize image.Point, anchor image.Point, CBorder CBorder) (*image.RGBA, error) {
	originalSize := img.Bounds().Size()
	p, err := calculatePaddings(kernelSize, anchor)
	if err != nil {
		return nil, err
	}
	rect := getRectangleFromPaddings(p, originalSize)
	padded := image.NewRGBA(rect)

	ParallelForEachPixel(originalSize, func(x, y int) {
		padded.Set(x+p.PaddingLeft, y+p.PaddingTop, img.RGBAAt(x, y))
	})

	switch CBorder {
	case CBorderConstant:
		// do nothing
	case CBorderReplicate:
		topPaddingReplicate(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		bottomPaddingReplicate(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		leftPaddingReplicate(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		rightPaddingReplicate(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
	case CBorderReflect:
		topPaddingReflect(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		bottomPaddingReflect(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		leftPaddingReflect(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		rightPaddingReflect(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
	default:
		return nil, errors.New("unknown CBorder type")
	}
	return padded, nil
}

// -------------------------------------------------------------------------------------------------------
func calculatePaddings(kernelSize image.Point, anchor image.Point) (Paddings, error) {
	var p Paddings
	if kernelSize.X < 0 || kernelSize.Y < 0 {
		return p, errors.New("negative size")
	}
	if anchor.X < 0 || anchor.Y < 0 {
		return p, errors.New("negative anchor value")
	}
	if anchor.X > kernelSize.X || anchor.Y > kernelSize.Y {
		return p, errors.New("anc" + "hor value outside of the kernel")
	}

	p = Paddings{PaddingLeft: anchor.X, PaddingRight: kernelSize.X - anchor.X - 1, PaddingTop: anchor.Y, PaddingBottom: kernelSize.Y - anchor.Y - 1}

	return p, nil
}

func getRectangleFromPaddings(p Paddings, imgSize image.Point) image.Rectangle {
	x := p.PaddingLeft + p.PaddingRight + imgSize.X
	y := p.PaddingTop + p.PaddingBottom + imgSize.Y
	return image.Rect(0, 0, x, y)
}
