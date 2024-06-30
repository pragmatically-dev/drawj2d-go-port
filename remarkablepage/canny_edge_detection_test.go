package remarkablepage

import (
	"image/png"
	"os"
	"testing"
)

func TestCannyEdge(t *testing.T) {
	img, _ := DecodeToGray("/home/nieva/Proyectos/drawj2d-rm/test-3-book.png")
	lowThreshold := 40.0
	highThreshold := 564.0
	img, _ = LaplacianGray(img, CBorderReflect, K8)
	edges := CannyEdgeDetection(img, lowThreshold, highThreshold)
	wr, _ := os.Create("/home/nieva/Proyectos/drawj2d-rm/image-of-test.png")

	png.Encode(wr, edges.SubImage(edges.Bounds()))
	
}
