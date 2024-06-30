package remarkablepage

import "fmt"

type Point struct {
	X, Y float64
}

type AffineTransform struct {
	matrix [6]float64
}

func NewAffineTransform() *AffineTransform {
	return &AffineTransform{
		matrix: [6]float64{1, 0, 0, 1, 0, 0},
	}
}

func GetHorizontalLines(pointMatrix [][]bool, imageWidth, imageHeight int) [][]float64 {
	var list [][]int
	var from, to int
	var isLine bool
	var isPixel bool

	for z := 0; z < imageHeight; z++ {
		isLine = false
		for x := 0; x < imageWidth; x++ {
			isPixel = pointMatrix[x][z]
			if isLine {
				if isPixel {
					to = x
					if x+1 == imageWidth {
						list = append(list, []int{z, from, to})
						isLine = false
					}
				} else { // !isPixel, line ended
					list = append(list, []int{z, from, to})
					isLine = false
				}
			} else { // !isLine, line not started
				if isPixel {
					from = x
					to = x
					if x+1 == imageWidth { // single pixel at last row
						list = append(list, []int{z, from, to})
					} else {
						isLine = true
					}
				} // else do nothing
			}
		}
	}
	if isLine {
		fmt.Println("Assertion failed: !isLine")
	}

	horizontalLines := make([][]float64, len(list))
	var startPoint, endPoint Point
	for i, line := range list {
		if len(line) != 3 {
			fmt.Println("Assertion failed: line length is not 3")
		}
		startPoint = Point{float64(line[1]), float64(line[0])}

		if line[1] != line[2] { // separate start and end points
			endPoint = Point{float64(line[2]), float64(line[0])}

			horizontalLines[i] = []float64{startPoint.X, startPoint.Y, endPoint.X, endPoint.Y}
		} else { // start and end points are equal, thus it is just a point
			horizontalLines[i] = []float64{startPoint.X, startPoint.Y}
		}
	}
	return horizontalLines
}
