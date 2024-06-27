/* package main

import (
	_ "image/jpeg"
	_ "image/png"

	t "github.com/pragmatically-dev/drawj2d-rm/remarkablepage"
)

func main() {
	t.TestRmDoc()
	/*
		file, err := os.Open("image.png")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}

		parser := NewClParserIMG(img)
		parser.bildeinlesen()

		//	for parser.hasMoreElements() {
		//		entity := parser.getEntity()
		//		fmt.Printf("Entity: %+v\n", entity)
		//	}

		horLines := parser.getHorLines()
		fmt.Printf("Horizontal lines: %+v\n", horLines)

		//pixelSize := parser.getPixelSize()
		//fmt.Printf("Pixel size: %+v\n", pixelSize)

}
*/

package main

import (
	_ "image/jpeg"
	_ "image/png"

	t "github.com/pragmatically-dev/drawj2d-rm/remarkablepage"
)

func main() {
	//t.Test()
	//t.TestRmDoc()
	t.TestCannyEdgeDetection("image.png")
	t.TestInvert("PostSobelResult.png")
	t.TestInvert("PostVertSobelResult.png")
	t.TestInvert("PostHorSobelResult.png")
	t.TestInvert("PostLaplacianResult.png")

}
