package remarkablepage

import (
	"fmt"
	"image"
	"log"
	"path/filepath"
	"strings"

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
			// Get pixel color
			originalColor := img.At(x, y)

			// Convert to grayscale
			r, _, _, _ := originalColor.RGBA()
			grayValue := uint8(r >> 8)

			// If the pixel is not black, assign the value to the output image and add a point to the reMarkable page
			if grayValue > 0 {
				page.AddPixel(float32(x), float32(y))
			}

		},
	)

	err = page.Export()
	if err != nil {
		log.Fatalln("Error exporting page:", err)
		return
	}

}

func CannyEdgeDetection(imagePath, DirToSave string) {

	// Check the file size
	fileInfo, err := os.Stat(imagePath)
	if err != nil {
		DebugPrint("Error getting file information:", err)
		return
	}

	predicateFilesize := fileInfo.Size() > 50*1024
	// If the file size is greater than 200 KB, perform resizing
	if predicateFilesize {

		DebugPrint("The file size is greater than 50 KB, agressive resizing will be performed.")
		img, err := im.ImreadGray(imagePath)
		if err != nil {
			DebugPrint("Error opening the file:", err)
			return
		}

		img, _ = rz.ResizeGray(img, 0.7, 0.7, rz.InterLinear)
		laplacianGray, _ := ed.LaplacianGray(img, padding.BorderReplicate, ed.K8)

		DetectWhitePixels(laplacianGray, imagePath, DirToSave)

	}

	if !predicateFilesize {

		DebugPrint("The file size is less than 50 KB, lite resizing will be performed.")
		img, err := im.ImreadGray(imagePath)
		if err != nil {
			DebugPrint("Error opening the file:", err)
			return
		}
		img, _ = rz.ResizeGray(img, 0.8, 0.8, rz.InterLinear)
		laplacianGray, _ := ed.LaplacianGray(img, padding.BorderReplicate, ed.K8)
		DetectWhitePixels(laplacianGray, imagePath, DirToSave)

	}

}

func DebugPrint(info string, opt ...error) {
	if debug {
		fmt.Println(info, opt)
	}
}
