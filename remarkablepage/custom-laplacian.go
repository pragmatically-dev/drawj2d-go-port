package remarkablepage

import (
	"errors"

	"image"
)

var kernel4 = Kernel{Content: [][]float64{
	{1, 4, 1},
	{4, -20, 4},
	{1, 4, 1},
}, Width: 3, Height: 3}

var kernel8 = Kernel{Content: [][]float64{
	{1, 1, 1},
	{1, -8, 1},
	{1, 1, 1},
}, Width: 3, Height: 3}

var sharpen = Kernel{Content: [][]float64{
	{0, -1, 0},
	{-1, 5, -1},
	{0, -1, 0},
}, Width: 3, Height: 3}

var kernel9 = Kernel{Content: [][]float64{
	{1, 1, 1},
	{1, -9, 1},
	{1, 1, 1},
}, Width: 3, Height: 3}

var kernel13 = Kernel{Content: [][]float64{
	{2, 2, 2},
	{2, -13, 2},
	{2, 2, 2},
}, Width: 3, Height: 3}

const a, b, c = 1.0 / 16.0, 2.0 / 16.0, 4.0 / 16.0

var gaussianBlur = Kernel{Content: [][]float64{
	{a, b, a},
	{b, c, b},
	{a, b, a},
}, Width: 3, Height: 3}

var sobel5x5X = Kernel{
	Content: [][]float64{
		{-1, -2, 0, 2, 1},
		{-2, -3, 0, 3, 2},
		{-3, -5, 0, 5, 3},
		{-2, -3, 0, 3, 2},
		{-1, -2, 0, 2, 1},
	},
	Width: 5, Height: 5,
}

var sobel5x5Y = Kernel{
	Content: [][]float64{
		{1, 2, 3, 2, 1},
		{2, 3, 5, 3, 2},
		{0, 0, 0, 0, 0},
		{-2, -3, -5, -3, -2},
		{-1, -2, -3, -2, -1},
	},
	Width: 5, Height: 5,
}

// LaplacianKernel - constant type for differentiating Laplacian kernels
type LaplacianKernel int

const (
	// K4 Laplacian kernel:
	//	{0, 1, 0},
	//	{1, -4, 1},
	//	{0, 1, 0},
	K4 LaplacianKernel = iota
	// K8 Laplacian kernel:
	//	{0, 1, 0},
	//	{1, -8, 1},
	//	{0, 1, 0},
	K8

	Sharpen

	K9

	K13
	SobelY
	SobelX
	Gaussian
	Sobel5x5X
	Sobel5x5Y
)

// LaplacianGray applies Laplacian filter to a grayscale image. The kernel types are: K4 and K8 (see LaplacianKernel)
// Example of usage:
//
//	res, err := edgedetection.LaplacianGray(img, paddding.BorderReflect, edgedetection.K8)
func LaplacianGray(gray *image.Gray, border CBorder, kernel LaplacianKernel) (*image.Gray, error) {
	var laplacianKernel Kernel
	switch kernel {
	case K4:
		laplacianKernel = kernel4
	case K8:
		laplacianKernel = kernel8
	case Sharpen:
		laplacianKernel = sharpen
	case K9:
		laplacianKernel = kernel9
	case K13:
		laplacianKernel = kernel13

	case Gaussian:
		laplacianKernel = gaussianBlur

	case Sobel5x5X:
		laplacianKernel = sobel5x5X
	case Sobel5x5Y:
		laplacianKernel = sobel5x5Y

	default:
		return nil, errors.New("invalid kernel")
	}
	return ConvolveGray(gray, &laplacianKernel, image.Point{X: 1, Y: 1}, border)
}
