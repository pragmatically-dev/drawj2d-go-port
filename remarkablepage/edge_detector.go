package remarkablepage

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"

	bl "github.com/ernyoke/imger/blur"
	ed "github.com/ernyoke/imger/edgedetection"
	im "github.com/ernyoke/imger/imgio"
	"github.com/ernyoke/imger/padding"
	re "github.com/ernyoke/imger/resize"
)

func TestCannyEdgeDetection(imagePath string) {

	img, err := im.ImreadGray(imagePath)
	if err != nil {
		fmt.Println("Error openning file:", err)
		return
	}
	img, _ = bl.GaussianBlurGray(img, 1, 2, padding.BorderReplicate)
	img, _ = re.ResizeGray(img, .5, .5, re.InterLinear)

	sobelGray, _ := ed.SobelGray(img, padding.BorderReplicate)
	vertSobelGray, _ := ed.VerticalSobelGray(img, padding.BorderReplicate)
	horSobelGray, _ := ed.HorizontalSobelGray(img, padding.BorderReplicate)
	laplacianGray, _ := ed.LaplacianGray(img, padding.BorderReplicate, ed.K8)

	im.Imwrite(sobelGray, "PostSobelResult.png")
	im.Imwrite(vertSobelGray, "PostVertSobelResult.png")
	im.Imwrite(horSobelGray, "PostHorSobelResult.png")
	im.Imwrite(laplacianGray, "PostLaplacianResult.png") //BEST RESULT

}

func invertColor(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(255 - r>>8),
		G: uint8(255 - g>>8),
		B: uint8(255 - b>>8),
		A: uint8(a >> 8),
	}
}

func invertImage(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	inverted := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			invertedColor := invertColor(originalColor)
			inverted.Set(x, y, invertedColor)
		}
	}

	return inverted
}

func TestInvert(path string) {
	inputFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	invertedImage := invertImage(img)

	outputFile, err := os.Create("INVERTED-" + path)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	if err := jpeg.Encode(outputFile, invertedImage, nil); err != nil {
		log.Fatalf("failed to encode image: %v", err)
	}
}
