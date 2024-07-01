package remarkablepage

import (
	"testing"
)

func TestBooleanMatrixBuilding(t *testing.T) {
	/* img, err := DecodeToGray("/home/nieva/Proyectos/drawj2d-rm/images/image.png")
	if err != nil {
		DebugPrint("Error opening the file:", err)
		return
	}

	config, _ := DecodeAndConfig("/home/nieva/Proyectos/drawj2d-rm/images/image.png")
	laplacianGray, _ := LaplacianGray(img, CBorderReplicate, K8)

	boolMap := BuildBooleanMatrix(laplacianGray)
	horLines := GetHorizontalLines(boolMap, config.Width, config.Height)

	rmFile := GetFileNameWithoutExtension("/home/nieva/Proyectos/drawj2d-rm/image.png")
	rmFile = fmt.Sprintf("%s/%s.rm", "/home/nieva/Proyectos/drawj2d-rm/", rmFile)

	rmRawData := DrawLines(horLines, float32(config.Width), float32(config.Height))
	zipData, zipName := CreateRmDoc(rmFile, rmRawData)

	file, _ := os.Create(zipName)
	zipData.WriteTo(file) */
}

func TestPNGConversion(t *testing.T) {
	/* 	rmData := LaplacianEdgeDetection("/home/nieva/Proyectos/drawj2d-rm/test-3-book.png")
	   	if rmData == nil {
	   		t.FailNow()
	   	}

	   	//fmt.Println("File testPNGConversion.rm generated successfully.")
	   	rmFile := GetFileNameWithoutExtension("/home/nieva/Proyectos/drawj2d-rm/test-3-book.png")
	   	rmFile = fmt.Sprintf("%s/%s.rm", "/home/nieva/Proyectos/drawj2d-rm/", rmFile)

	   	zipData, zipName := CreateRmDoc(rmFile, rmData)

	   	file, _ := os.Create(zipName)
	   	zipData.WriteTo(file) */

}

func TestDrawingPixel(t *testing.T) {

	/*
		 	page := NewReMarkablePage() // Asumiendo una altura de página de 1000 unidades
			var lineas [][]float32 = [][]float32{
				{384, 41, 419, 41}, {610, 41, 629, 41}, {796, 41, 831, 41}, {1344, 41, 1351, 41},
				{384, 42}, {419, 42}, {610, 42}, {629, 42}, {796, 42}, {831, 42},
				{1341, 42, 1354, 42}, {384, 43}, {388, 43, 416, 43}, {419, 43}, {610, 43},
				{629, 43}, {796, 43}, {800, 43, 828, 43}, {831, 43}, {1339, 43, 1343, 43},
				{1352, 43, 1356, 43}, {384, 44}, {387, 44, 417, 44}, {419, 44}, {610, 44, 629, 44},
				{796, 44}, {799, 44, 828, 44}, {831, 44}, {1338, 44, 1357, 44}, {33, 45, 72, 45},
				{384, 45}, {387, 45}, {416, 45, 417, 45}, {419, 45}, {796, 45}, {799, 45},
				{828, 45}, {831, 45}, {1337, 45, 1344, 45}, {1351, 45, 1358, 45}, {33, 46},
				{72, 46}, {384, 46}, {387, 46}, {416, 46, 417, 46}, {419, 46}, {796, 46},
				{799, 46}, {828, 46}, {831, 46}, {1336, 46, 1341, 46}, {1354, 46, 1359, 46},
				{33, 47}, {72, 47}, {384, 47}, {387, 47}, {416, 47, 417, 47}, {419, 47},
				{796, 47}, {799, 47}, {828, 47}, {831, 47}, {1335, 47, 1340, 47}, {1355, 47, 1360, 47},
				{33, 48, 72, 48}, {384, 48}, {387, 48}, {416, 48, 417, 48}, {419, 48}, {796, 48},
				{799, 48}, {828, 48}, {831, 48}, {1334, 48, 1339, 48}, {1356, 48, 1361, 48},
				{380, 49, 384, 49}, {387, 49, 389, 49}, {416, 49, 417, 49}, {419, 49},
				{610, 49, 649, 49}, {796, 49}, {799, 49}, {828, 49}, {831, 49},
				{1334, 49, 1338, 49}, {1357, 49, 1361, 49}, {380, 50, 383, 50},
				{388, 50, 389, 50}, {416, 50, 417, 50}, {419, 50}, {501, 50, 504, 50},
				{550, 50, 553, 50}, {610, 50, 649, 50}, {699, 50, 702, 50},
				{716, 50, 719, 50}, {796, 50}, {799, 50}, {828, 50}, {831, 50},
				{890, 50, 891, 50}, {913, 50, 916, 50}, {950, 50, 953, 50},
				{1333, 50, 1337, 50}, {1358, 50, 1362, 50}, {90, 51, 95, 51},
				{105, 51, 110, 51}, {380, 51}, {389, 51}, {416, 51, 417, 51},
				{419, 51}, {437, 51, 442, 51}, {450, 51, 452, 51}, {476, 51, 478, 51},
				{501, 51}, {503, 51, 504, 51}, {550, 51}, {552, 51, 553, 51},
				{610, 51}, {649, 51}, {667, 51, 679, 51}, {699, 51},
				{701, 51, 702, 51}, {716, 51}, {718, 51, 719, 51},
			}
			page.AddPixel(500, 500)

			//fmt.Println("File testPNGConversion.rm generated successfully.")
			rmFile := "test-pixel"
			rmFile = fmt.Sprintf("%s/%s.rm", "/home/nieva/Proyectos/drawj2d-rm/", rmFile)

			for _, line := range lineas {
				if len(line) == 4 {
					ln := page.AddLine()
					ln.AddPoint(line[0], line[1])
					ln.AddPoint(line[2], line[3])
				}
				if len(line) == 2 {
					page.AddPixel(line[0], line[1])
				}
			}

			rawData := page.Export()
			zipData, zipName := CreateRmDoc(rmFile, rawData)
			file, _ := os.Create(zipName)
			zipData.WriteTo(file)

			fmt.Println("File testRemarkablePageSmiley.rm generated successfully.")
	*/
}
