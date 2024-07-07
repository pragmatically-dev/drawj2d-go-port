package remarkablepage

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"
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
func DrawLines(lines *LineList, width, height float32) []byte {

	page := NewReMarkablePage()
	defer func() {
		page = nil
		lines.Lines = nil
		lines.Size = 0
		lines = nil
	}()
	for i := 0; i < lines.Size; i++ {

		p1, p2, p3, p4 := lines.Lines[i*4], lines.Lines[i*4+1], lines.Lines[i*4+2], lines.Lines[i*4+3]

		if p1 != p3 || p2 != p4 {
			ln := page.AddLine()
			ln.AddPoint(p1, p2)
			ln.AddPoint(p3, p4)

		} else {
			ln := page.AddLine()
			ln.AddPoint(p1, p2)

		}

	}

	return page.Export()
}

func BuildBooleanMatrix(img *image.Gray) [][]bool {
	bounds := img.Bounds().Size()
	width, height := bounds.X, bounds.Y

	boolImgMap := make([][]bool, width)
	for i := range boolImgMap {
		boolImgMap[i] = make([]bool, height)
	}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			boolImgMap[i][j] = img.GrayAt(i, j).Y > 0
		}
	}

	return boolImgMap
}

func LaplacianEdgeDetection(imagePath string) []byte {

	dir, filep := filepath.Dir(imagePath), filepath.Base(imagePath)

	horLines := HandleNewFile(dir, filep)

	return DrawLines(horLines, float32(X_MAX), float32(Y_MAX))

}

func DebugPrint(info string, opt ...error) {
	if debug {
		fmt.Println(info, opt)
	}
}
