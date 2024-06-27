package remarkablepage

import (
	"fmt"
	"image"
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

func DetectWhitePixels(img *image.Gray, filename string) {
	rmFile := GetFileNameWithoutExtension(filename)

	file, err := os.Create(fmt.Sprintf("%s.rm", rmFile))
	if err != nil {
		debugPrint("Error creating file:", err)
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
		debugPrint("Error exporting page:", err)
	}

	debugPrint("File testPNGConversion.rm generated successfully.")
}

func TestCannyEdgeDetection(imagePath string) {

	// Check the file size
	fileInfo, err := os.Stat(imagePath)
	if err != nil {
		debugPrint("Error getting file information:", err)
		return
	}

	// If the file size is greater than 200 KB, perform resizing
	if fileInfo.Size() > 50*1024 {
		debugPrint("The file size is greater than 200 KB, resizing will be performed.")
		img, err := im.ImreadGray(imagePath)
		if err != nil {
			debugPrint("Error opening the file:", err)
			return
		}

		img, _ = rz.ResizeGray(img, 0.7, 0.7, rz.InterLinear)
		laplacianGray, _ := ed.LaplacianGray(img, padding.BorderReplicate, ed.K8)

		DetectWhitePixels(laplacianGray, imagePath)

	} else {
		debugPrint("The file size is less than 200 KB, resizing will not be performed.")
		img, err := im.ImreadGray(imagePath)
		if err != nil {
			debugPrint("Error opening the file:", err)
			return
		}

		laplacianGray, _ := ed.LaplacianGray(img, padding.BorderReplicate, ed.K8)

		DetectWhitePixels(laplacianGray, imagePath)

	}
}

func debugPrint(info string, opt ...error) {
	if debug {
		fmt.Println(info, opt)
	}
}
