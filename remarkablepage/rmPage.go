package remarkablepage

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"math"
	"os"
	"sync"
)

// Constants
const (
	VERSION   = "0.3"
	REV_DATE  = "2023-04-08"
	X_MAX     = 1404.0
	Y_MAX     = 1872.0
	HEADER_V5 = "reMarkable .lines file, version=5          "
)

// ReMarkablePage represents a page for the reMarkable tablet
type ReMarkablePage struct {
	lines      []*rmLine
	debug      bool
	out        *os.File
	colors     map[string]color.RGBA
	pageHeight float32
	mu         sync.Mutex // Add a mutex for thread safety
}

// rmLine represents a line on the reMarkable page
type rmLine struct {
	brushType            int32
	color                int32
	padding              int32 // ?
	brushBaseSize        float32
	unknownLineAttribute float32 // ?
	pointList            []*rmPoint
}

// rmPoint represents a point on a reMarkable line
type rmPoint struct {
	x, y      float32
	speed     float32
	direction float32 // 3.1415f; // ev. tilt
	width     float32 // 0.3f * 226.85f/25.4f; for 0.30mm
	pressure  float32
}

// NewReMarkablePage creates a new reMarkable page
func NewReMarkablePage(out *os.File, pageHeight float32) *ReMarkablePage {
	return &ReMarkablePage{
		lines: make([]*rmLine, 0),
		debug: false,
		out:   out,
		colors: map[string]color.RGBA{
			"red":   {R: 217, G: 7, B: 7, A: 255},
			"blue":  {R: 0, G: 98, B: 204, A: 255},
			"black": {R: 0, G: 0, B: 0, A: 255},
		},
		pageHeight: pageHeight,
	}
}

// AddLine adds a new line to the page
func (page *ReMarkablePage) AddLine() *rmLine {
	page.mu.Lock() // Lock the mutex before accessing shared resources
	defer page.mu.Unlock()

	line := &rmLine{
		pointList: make([]*rmPoint, 0),
	}
	page.lines = append(page.lines, line)
	if page.debug {
		fmt.Printf("[RemarkablePage] line added. Nb lines: %d\n", len(page.lines))
	}
	return line
}

// AddPoint adds a point to a line
func (line *rmLine) AddPoint(x, y float32) *rmPoint {
	point := &rmPoint{
		x:         x,
		y:         y,
		speed:     0.1,
		direction: 0,
		width:     2.125,
		pressure:  1.0,
	}
	line.pointList = append(line.pointList, point)
	return point
}

// Export writes the content of the page to the output file
func (page *ReMarkablePage) Export() error {
	defer page.out.Close()

	// Write the header
	header := []byte(HEADER_V5)
	_, err := page.out.Write(header)
	if err != nil {
		return err
	}

	// Write the page
	nbLayers := int32(1)
	err = binary.Write(page.out, binary.LittleEndian, nbLayers)
	if err != nil {
		return err
	}

	// Write the layers
	err = page.writeLayer()
	if err != nil {
		return err
	}

	return nil
}

// writeLayer writes a layer of lines to the output file
func (page *ReMarkablePage) writeLayer() error {
	page.mu.Lock() // Lock the mutex before accessing shared resources
	defer page.mu.Unlock()

	err := binary.Write(page.out, binary.LittleEndian, int32(len(page.lines)))
	if err != nil {
		return err
	}

	for _, line := range page.lines {
		err := page.writeLine(line)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeLine writes a line and its points to the output file
func (page *ReMarkablePage) writeLine(line *rmLine) error {
	err := binary.Write(page.out, binary.LittleEndian, line.brushType)
	if err != nil {
		return err
	}
	err = binary.Write(page.out, binary.LittleEndian, line.color)
	if err != nil {
		return err
	}
	err = binary.Write(page.out, binary.LittleEndian, line.padding)
	if err != nil {
		return err
	}
	err = binary.Write(page.out, binary.LittleEndian, line.brushBaseSize)
	if err != nil {
		return err
	}
	err = binary.Write(page.out, binary.LittleEndian, line.unknownLineAttribute)
	if err != nil {
		return err
	}

	nbPoints := int32(len(line.pointList))
	err = binary.Write(page.out, binary.LittleEndian, nbPoints)
	if err != nil {
		return err
	}

	for _, point := range line.pointList {
		err := page.writePoint(point)
		if err != nil {
			return err
		}
	}

	if page.debug {
		fmt.Printf("                 line with points: %d\n", len(line.pointList))
	}
	return nil
}

// writePoint writes a point to the output file
func (page *ReMarkablePage) writePoint(point *rmPoint) error {
	err := binary.Write(page.out, binary.LittleEndian, point.x)
	if err != nil {
		return err
	}
	err = binary.Write(page.out, binary.LittleEndian, point.y)
	if err != nil {
		return err
	}
	err = binary.Write(page.out, binary.LittleEndian, point.speed)
	if err != nil {
		return err
	}
	err = binary.Write(page.out, binary.LittleEndian, point.direction)
	if err != nil {
		return err
	}
	err = binary.Write(page.out, binary.LittleEndian, point.width)
	if err != nil {
		return err
	}
	err = binary.Write(page.out, binary.LittleEndian, point.pressure)
	if err != nil {
		return err
	}
	return nil
}

// transformPoint transforms a point to the new coordinate system
func (page *ReMarkablePage) transformPoint(x, y float32) (float32, float32) {
	return x, page.pageHeight - y
}

// DrawCircle draws a circle on the page
func (page *ReMarkablePage) DrawCircle(centerX, centerY, radius float32) {
	line := page.AddLine()
	numSegments := 360
	for i := 0; i <= numSegments; i++ {
		theta := float64(i) * 2.0 * math.Pi / float64(numSegments)
		x := centerX + radius*float32(math.Cos(theta))
		y := centerY + radius*float32(math.Sin(theta))
		//tx, ty := page.transformPoint(x, y)
		line.AddPoint(x, y)
	}
}

// DrawBezierCurve draws a Bezier curve on the page
func (page *ReMarkablePage) DrawBezierCurve(p0, p1, p2, p3 rmPoint) {
	line := page.AddLine()
	numSegments := 100
	for i := 0; i <= numSegments; i++ {
		t := float32(i) / float32(numSegments)
		x := (1-t)*(1-t)*(1-t)*p0.x + 3*(1-t)*(1-t)*t*p1.x + 3*(1-t)*t*t*p2.x + t*t*t*p3.x
		y := (1-t)*(1-t)*(1-t)*p0.y + 3*(1-t)*(1-t)*t*p1.y + 3*(1-t)*t*t*p2.y + t*t*t*p3.y
		tx, ty := page.transformPoint(x, y)
		line.AddPoint(tx, ty)
	}
}

// AddPixel adds a pixel to the page
func (page *ReMarkablePage) AddPixel(x, y float32) {
	line := page.AddLine()
	const c = 0.1

	line.AddPoint(x-c, y)
	line.AddPoint(x, y)
	line.AddPoint(x+c, y)
}

// DrawFilledRectangle draws a filled rectangle on the page
func (page *ReMarkablePage) DrawFilledRectangle(x1, y1, x2, y2 float32) {
	line := page.AddLine()

	// Bottom left to bottom right
	for x := x1; x <= x2; x++ {
		//tx, ty := page.transformPoint(x, y1)
		line.AddPoint(x, y1)
	}

	// Bottom right to top right
	for y := y1; y <= y2; y++ {
		//tx, ty := page.transformPoint(x2, y)
		line.AddPoint(x2, y)
	}

	// Top right to top left
	for x := x2; x >= x1; x-- {
		//tx, ty := page.transformPoint(x, y2)
		line.AddPoint(x, y2)
	}

	// Top left to bottom left (closing the rectangle)
	for y := y2; y >= y1; y-- {
		//tx, ty := page.transformPoint(x1, y)
		line.AddPoint(x1, y)
	}
}
