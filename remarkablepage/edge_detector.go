package remarkablepage

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"
	"sync"
)

const (
	debug = false
)

func GetFileNameWithoutExtension(filePath string) string {
	// Get the base name of the file
	base := filepath.Base(filePath)
	// Get the extension
	ext := filepath.Ext(base)
	// Remove the extension from the base name
	return strings.TrimSuffix(base, ext)
}

// DetectWhitePixels detects white pixels in a grayscale image and adds them to a reMarkable page
func DrawLines(coordinates [][]float64, width, height float32) []byte {

	page := NewReMarkablePage()

	for _, line := range coordinates {
		isLine := len(line) == 4

		if isLine {
			ln := page.AddLine()
			ln.AddPoint(float32(line[0]), float32(line[1]))
			ln.AddPoint(float32(line[2]), float32(line[3]))

		}

		if !isLine {
			ln := page.AddLine()
			ln.AddPoint(float32(line[0]), float32(line[1]))
			ln.AddPoint(float32(line[0]-0.1), float32(line[1]))

		}
	}

	return page.Export()
}

func BuildBooleanMatrix(img *image.Gray) [][]bool {
	bounds := img.Bounds().Size()
	width, height := bounds.X, bounds.Y

	// Inicializar la matriz booleana con el tamaño adecuado
	boolImgMap := make([][]bool, width)

	for i := range boolImgMap {
		boolImgMap[i] = make([]bool, height)
	}

	var wg sync.WaitGroup

	// self documented bruh
	processRow := func(y int) {
		defer wg.Done()
		for x := 0; x < width; x++ {

			boolImgMap[x][y] = img.GrayAt(x, y).Y > 0

		}
	}

	// spawn go routines for each row
	for y := 0; y < height; y++ {
		wg.Add(1)
		go processRow(y)
	}

	wg.Wait()
	return boolImgMap
}

// DetectWhitePixels detects white pixels in a grayscale image and adds them to a reMarkable page
func DetectWhitePixels(img *image.Gray, filename, dirToSave string) []byte {

	page := NewReMarkablePage()

	size := img.Bounds().Max
	var wg sync.WaitGroup

	// self documented bruh
	processRow := func(y int) {
		defer wg.Done()
		for x := 0; x < size.X; x++ {
			if img.GrayAt(x, y).Y > 0 {
				page.AddPixel(float32(x), float32(y))
			}
		}
	}

	// spawn go routines for each row
	for y := 0; y < size.Y; y++ {
		wg.Add(1)
		go processRow(y)
	}

	wg.Wait()

	return page.Export()

}

func LaplacianEdgeDetection(imagePath string) []byte {

	img, err := DecodeToGray(imagePath)
	if err != nil {
		DebugPrint("Error opening the file:", err)
		return nil
	}

	img, _ = LaplacianGray(img, CBorderReplicate, K8)

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	horLines := GetHorizontalLines(BuildBooleanMatrix(img), width, height)

	return DrawLines(horLines, float32(width), float32(height))

}

func DebugPrint(info string, opt ...error) {
	if debug {
		fmt.Println(info, opt)
	}
}
