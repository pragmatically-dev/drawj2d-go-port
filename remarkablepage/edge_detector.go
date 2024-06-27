package remarkablepage

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"

	ed "github.com/ernyoke/imger/edgedetection"
	ef "github.com/ernyoke/imger/effects"
	im "github.com/ernyoke/imger/imgio"
	"github.com/ernyoke/imger/padding"
)

// DetectWhitePixels detecta los píxeles blancos en una imagen y los exporta a un archivo reMarkable
func DetectWhitePixels(img *image.Gray) *image.Gray {
	file, err := os.Create("testPNGConversion.rm")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil
	}
	defer file.Close()

	page := NewReMarkablePage(file, float32(Y_MAX)) // Asumiendo una altura de página de 1872 unidades

	// Obtener dimensiones de la imagen
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Crear una imagen en escala de grises para almacenar los píxeles negros
	grayImg := image.NewGray(image.Rect(0, 0, width, height))
	grayImg = ef.InvertGray(grayImg)
	// Iterar sobre cada píxel de la imagen original
	for y := 0; y < height; y = y + 1 {

		for x := 0; x < width; x = x + 1 {
			// Obtener color del píxel
			originalColor := img.At(x, y)

			// Convertir a escala de grises
			r, _, _, _ := originalColor.RGBA()
			grayValue := uint8(r >> 8)

			// Si es un píxel blanco, asignar el valor a la imagen de salida y agregar un punto a la página reMarkable
			if grayValue > 0 {
				grayImg.SetGray(x, y, color.Gray{0})
				line := page.AddLine()
				var c float32 = 0.3
				line.AddPoint(float32(x)-c, float32(y)-c)
				line.AddPoint(float32(x), float32(y))
				line.AddPoint(float32(x)+c, float32(y)+c)

			}
		}
	}

	err = page.Export()
	if err != nil {
		fmt.Println("Error exporting page:", err)
		return nil
	}

	fmt.Println("File testPNGConversion.rm generated successfully.")
	return grayImg
}

func TestCannyEdgeDetection(imagePath string) {

	img, err := im.ImreadGray(imagePath)
	if err != nil {
		fmt.Println("Error openning file:", err)
		return
	}
	//img, _ = bl.GaussianBlurGray(img, 1, 2, padding.BorderReplicate)

	//sobelGray, _ := ed.SobelGray(img, padding.BorderReplicate)
	//vertSobelGray, _ := ed.VerticalSobelGray(img, padding.BorderReplicate)
	//horSobelGray, _ := ed.HorizontalSobelGray(img, padding.BorderReplicate)
	laplacianGray, _ := ed.LaplacianGray(img, padding.BorderReplicate, ed.K8)

	//im.Imwrite(sobelGray, "PostSobelResult.png")
	//im.Imwrite(vertSobelGray, "PostVertSobelResult.png")
	//im.Imwrite(horSobelGray, "PostHorSobelResult.png")
	im.Imwrite(laplacianGray, "PostLaplacianResult.png") //BEST RESULT

	cleanerImg := DetectWhitePixels(laplacianGray)
	im.Imwrite(cleanerImg, "PostCleaning.png")
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
