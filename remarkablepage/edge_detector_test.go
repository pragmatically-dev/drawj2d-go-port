package remarkablepage

import (
	"fmt"
	"os"
	"testing"
)

/* func TestBooleanMatrixBuilding(t *testing.T) {
	img, err := im.ImreadGray("image.png")
	if err != nil {
		DebugPrint("Error opening the file:", err)
		return
	}
	laplacianGray, _ := ed.LaplacianGray(img, padding.BorderReplicate, ed.K8)

	BuildBooleanMatrix(laplacianGray)

} */

func TestPNGConversion(t *testing.T) {
	rmData := LaplacianEdgeDetection("/home/nieva/Proyectos/drawj2d-rm/image.png", "/home/nieva/Proyectos/drawj2d-rm/")
	if rmData == nil {
		t.FailNow()
	}

	//fmt.Println("File testPNGConversion.rm generated successfully.")
	rmFile := GetFileNameWithoutExtension("/home/nieva/Proyectos/drawj2d-rm/image.png")
	rmFile = fmt.Sprintf("%s/%s.rm", "/home/nieva/Proyectos/drawj2d-rm/", rmFile)

	zipData, zipName := CreateRmDoc(rmFile, "/home/nieva/Proyectos/drawj2d-rm/", rmData)

	file, _ := os.Create(zipName)
	zipData.WriteTo(file)

}

/*
func TestDrawingPixel(t *testing.T) {
	file, err := os.Create("test.rm")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	page := NewReMarkablePage(file, 1872) // Asumiendo una altura de página de 1000 unidades

	page.AddPixel(500, 500)

	err = page.Export()
	if err != nil {
		fmt.Println("Error exporting page:", err)
		return
	}
	_ = CreateRmDoc("/home/nieva/Proyectos/drawj2d-rm/test.rm", "")

	fmt.Println("File testRemarkablePageSmiley.rm generated successfully.")
}
*/
