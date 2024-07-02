package remarkablepage

import (
	"image/png"
	"os"
	"testing"
)

func TestHoughTransform(t *testing.T) {
	f, _ := DecodeToGray("/home/nieva/Proyectos/drawj2d-rm/images/image.png")
	config, _ := DecodeAndConfig("/home/nieva/Proyectos/drawj2d-rm/images/image.png")
	tmp, _ := os.Create("/home/nieva/Proyectos/drawj2d-rm/images/hough.png")

	h := Hough(f, config.Width, config.Height)

	png.Encode(tmp, h)

}
