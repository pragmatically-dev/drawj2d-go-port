package remarkablepage

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"
	"sync"

	"os"

	ed "github.com/ernyoke/imger/edgedetection"
	im "github.com/ernyoke/imger/imgio"
	"github.com/ernyoke/imger/padding"
	rz "github.com/ernyoke/imger/resize"
	"github.com/ernyoke/imger/utils"
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

	utils.ParallelForEachPixel(bounds,

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

/*
func DetectWhitePixels(img *image.Gray, filename, dirToSave string) {
	rmFile := GetFileNameWithoutExtension(filename)

	file, err := os.Create(fmt.Sprintf("%s/%s.rm", dirToSave, rmFile))
	if err != nil {
		DebugPrint("Error creating file:", err)
		return
	}
	defer file.Close()

	page := NewReMarkablePage(file, float32(Y_MAX))

	size := img.Bounds().Max
	utils.ParallelForEachPixel(size,

		func(x, y int) {
			// If the pixel is not black, assign the value to the output image and add a point to the reMarkable page
			if img.GrayAt(x, y).Y > 0 {
				page.AddPixel(float32(x), float32(y))
			}

		},
	)

	err = page.Export()
	if err != nil {
		log.Fatalln("Error exporting page:", err)
		return
	}

} */

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
func LaplacianEdgeDetection(imagePath, DirToSave string) []byte {

	// Check the file size
	fileInfo, err := os.Stat(imagePath)
	if err != nil {
		DebugPrint("Error getting file information:", err)
		return nil
	}

	predicateFilesize := fileInfo.Size() > 50*1024
	// If the file size is greater than 50 KB, perform resizing
	if predicateFilesize {

		DebugPrint("The file size is greater than 50 KB, agressive resizing will be performed.")
		img, err := im.ImreadGray(imagePath)
		if err != nil {
			DebugPrint("Error opening the file:", err)
			return nil
		}

		img, _ = rz.ResizeGray(img, 0.75, 0.75, rz.InterLinear)
		laplacianGray, _ := ed.LaplacianGray(img, padding.BorderReplicate, ed.K8)

		return DetectWhitePixels(laplacianGray, imagePath, DirToSave)

	}

	if !predicateFilesize {

		DebugPrint("The file size is less than 50 KB, lite resizing will be performed.")
		img, err := im.ImreadGray(imagePath)
		if err != nil {
			DebugPrint("Error opening the file:", err)
			return nil
		}
		img, _ = rz.ResizeGray(img, 0.8, 0.8, rz.InterLinear)
		laplacianGray, _ := ed.LaplacianGray(img, padding.BorderReplicate, ed.K8)
		return DetectWhitePixels(laplacianGray, imagePath, DirToSave)

	}

	return nil
}

func DebugPrint(info string, opt ...error) {
	if debug {
		fmt.Println(info, opt)
	}
}
