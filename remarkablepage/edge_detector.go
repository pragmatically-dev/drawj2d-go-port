package remarkablepage

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"
	"sync"

	"os"
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

func BuildBooleanMatrix(img *image.Gray) [][]bool {
	bounds := img.Bounds().Size()
	width, height := bounds.X, bounds.Y

	// Inicializar la matriz booleana con el tamaño adecuado
	boolImgMap := make([][]bool, width)
	for i := range boolImgMap {
		boolImgMap[i] = make([]bool, height)
	}

	ParallelForEachPixel(bounds,

		func(x, y int) {

			// Get pixel color
			originalColor := img.At(x, y)

			// Convert to grayscale
			r, _, _, _ := originalColor.RGBA()
			grayValue := uint8(r >> 8)

			// Si el valor de gris es mayor a 0, asignar true, sino false
			boolImgMap[x][y] = grayValue > 0
		},
	)

	return boolImgMap
}

// DetectWhitePixels detects white pixels in a grayscale image and adds them to a reMarkable page
func DetectWhitePixels(img *image.Gray, filename, dirToSave string) []byte {

	page := NewReMarkablePage()

	size := img.Bounds().Max
	var wg sync.WaitGroup

	// Función que procesa una fila de píxeles
	processRow := func(y int) {
		defer wg.Done()
		for x := 0; x < size.X; x++ {
			if img.GrayAt(x, y).Y > 0 {
				page.AddPixel(float32(x), float32(y))
			}
		}
	}

	// Lanzar gorutinas para procesar cada fila
	for y := 0; y < size.Y; y++ {
		wg.Add(1)
		go processRow(y)
	}

	wg.Wait()

	return page.Export()

}
func LaplacianEdgeDetection(imagePath, dirToSave string) []byte {

	fileInfo, err := os.Stat(imagePath)
	if err != nil {
		DebugPrint("Error getting file information:", err)
		return nil
	}

	fileSize := fileInfo.Size()
	resizeFactor := 0.85
	if fileSize > 50*1024 {
		DebugPrint("The file size is greater than 50 KB, aggressive resizing will be performed.")
		resizeFactor = 0.75
	} else {
		DebugPrint("The file size is less than 50 KB, lite resizing will be performed.")
	}

	img, err := DecodeToGray(imagePath)
	if err != nil {
		DebugPrint("Error opening the file:", err)
		return nil
	}

	img, _ = ResizeGray(img, resizeFactor, resizeFactor, InterLinear)
	laplacianGray, _ := LaplacianGray(img, CBorderReplicate, K8)
	return DetectWhitePixels(laplacianGray, imagePath, dirToSave)

}

func DebugPrint(info string, opt ...error) {
	if debug {
		fmt.Println(info, opt)
	}
}
