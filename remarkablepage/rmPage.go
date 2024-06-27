package remarkablepage

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"math"
	"os"
)

// -------------------------------
// Constants

// see https://github.com/bsdz/remarkable-layers/blob/master/rmlines/constants.py
/*     static String HEADER_V5 = "reMarkable .lines file, version=5          ";
static final float X_MAX = 1404f;
static final float Y_MAX = 1872f;

static final class rmColour {
    static final int BLACK = 0;
    static final int GREY = 1;
    static final int WHITE = 2;
    // see https://github.com/ricklupton/rmscene/blob/main/src/rmscene/scene_items.py
    static final int BLUE = 6;
    static final int RED = 7;
}

static final class rmPen {
    // see https://github.com/ax3l/lines-are-beautiful/blob/develop/include/rmlab/Line.hpp
    static final int BRUSH = 0;
    static final int PENCIL_TILT = 1;
    static final int BALLPOINT_PEN_1 = 2;
    static final int MARKER_1 = 3;
    static final int FINELINER_1 = 4;
    static final int HIGHLIGHTER = 5;
    static final int RUBBER = 6;  // used in version 5
    static final int PENCIL_SHARP = 7;
    static final int RUBBER_AREA = 8;
    static final int ERASE_ALL = 9;
    static final int SELECTION_BRUSH_1 = 10;
    static final int SELECTION_BRUSH_2 = 11;
    // below used for version 5;
    static final int PAINT_BRUSH_1 = 12;
    static final int MECHANICAL_PENCIL_1 = 13;
    static final int PENCIL_2 = 14;
    static final int BALLPOINT_PEN_2 = 15;
    static final int MARKER_2 = 16;
    static final int FINELINER_2 = 17;
    static final int HIGHLIGHTER_2 = 18;
    static final int DEFAULT = FINELINER_2;
}

static final class rmWidth {
    static final float SMALL = 1.875f;
    static final float MEDIUM = 2.0f;
    static final float LARGE = 2.125f;
} */

// Constantes
const (
	VERSION   = "0.3"
	REV_DATE  = "2023-04-08"
	X_MAX     = 1404.0
	Y_MAX     = 1872.0
	HEADER_V5 = "reMarkable .lines file, version=5          "
)

// ReMarkablePage representa una página para la tableta reMarkable

type ReMarkablePage struct {
	lines      []*rmLine
	debug      bool
	out        *os.File
	colors     map[string]color.RGBA
	pageHeight float32
}

// rmLine representa una línea en la página reMarkable

type rmLine struct {
	brushType            int32
	color                int32
	padding              int32 // ?
	brushBaseSize        float32
	unknownLineAttribute float32 // ?
	pointList            []*rmPoint
}

// rmPoint representa un punto en una línea reMarkable

type rmPoint struct {
	x, y      float32
	speed     float32
	direction float32 // 3.1415f; // ev. tilt
	width     float32 // 0.3f * 226.85f/25.4f; for 0.30mm
	pressure  float32
}

// NewReMarkablePage crea una nueva página reMarkable

func NewReMarkablePage(out *os.File, pageHeight float32) *ReMarkablePage {
	page := &ReMarkablePage{
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
	return page
}

// AddLine agrega una nueva línea a la página

func (page *ReMarkablePage) AddLine() *rmLine {
	line := &rmLine{
		pointList: make([]*rmPoint, 0),
	}
	page.lines = append(page.lines, line)
	if page.debug {
		fmt.Printf("[RemarkablePage] line added. Nb lines: %d\n", len(page.lines))
	}
	return line
}

// AddPoint agrega un punto a una línea

func (line *rmLine) AddPoint(x, y float32) *rmPoint {
	point := &rmPoint{
		x: x,
		y: y,
	}
	line.pointList = append(line.pointList, point)
	return point
}

// Export escribe el contenido de la página al archivo de salida

func (page *ReMarkablePage) Export() error {
	defer page.out.Close()

	// Escribir el encabezado
	header := []byte(HEADER_V5)
	_, err := page.out.Write(header)
	if err != nil {
		return err
	}

	// Escribir la página
	nbLayers := int32(1)
	err = binary.Write(page.out, binary.LittleEndian, nbLayers)
	if err != nil {
		return err
	}

	// Escribir las capas
	err = page.writeLayer()
	if err != nil {
		return err
	}

	return nil
}

// writeLayer escribe una capa de líneas en el archivo de salida

func (page *ReMarkablePage) writeLayer() error {
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

// writeLine escribe una línea y sus puntos en el archivo de salida

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

// writePoint escribe un punto en el archivo de salida

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

// transformPoint transforma un punto al nuevo sistema de coordenadas

func (page *ReMarkablePage) transformPoint(x, y float32) (float32, float32) {
	return x, page.pageHeight - y
}

// DrawCircle dibuja un círculo en la página

func (page *ReMarkablePage) DrawCircle(centerX, centerY, radius float32) {
	line := page.AddLine()
	numSegments := 360
	for i := 0; i <= numSegments; i++ {
		theta := float64(i) * 2.0 * math.Pi / float64(numSegments)
		x := centerX + radius*float32(math.Cos(theta))
		y := centerY + radius*float32(math.Sin(theta))
		tx, ty := page.transformPoint(x, y)
		line.AddPoint(tx, ty)
	}
}

// DrawBezierCurve dibuja una curva de Bézier en la página

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

// DrawFilledRectangle dibuja un rectángulo relleno en la página
func (page *ReMarkablePage) DrawFilledRectangle(x1, y1, x2, y2 float32) {
	line := page.AddLine()

	// Bottom left to bottom right
	for x := x1; x <= x2; x++ {
		tx, ty := page.transformPoint(x, y1)
		line.AddPoint(tx, ty)
	}

	// Bottom right to top right
	for y := y1; y <= y2; y++ {
		tx, ty := page.transformPoint(x2, y)
		line.AddPoint(tx, ty)
	}

	// Top right to top left
	for x := x2; x >= x1; x-- {
		tx, ty := page.transformPoint(x, y2)
		line.AddPoint(tx, ty)
	}

	// Top left to bottom left (closing the rectangle)
	for y := y2; y >= y1; y-- {
		tx, ty := page.transformPoint(x1, y)
		line.AddPoint(tx, ty)
	}
}

// Test genera un archivo .rm de prueba con una carita feliz

func Test() {
	file, err := os.Create("testRemarkablePageSmiley.rm")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	page := NewReMarkablePage(file, 1872) // Asumiendo una altura de página de 1000 unidades

	center := &rmPoint{X_MAX / 2, Y_MAX / 2, 0, 0, 0, 0}
	// Dibujar la cara
	page.DrawCircle(center.x, center.y, 400)

	// Dibujar los ojos
	page.DrawCircle(center.x*1.2, center.y*1.25, 70)
	page.DrawCircle(center.x*.85, center.y*1.25, 70)

	for i := 0; i < 6; i++ {

		page.DrawCircle(center.x*1.2, center.y*1.25, float32(i*10))
		page.DrawCircle(center.x*.85, center.y*1.25, float32(i*10))

	}

	p0 := rmPoint{center.x * .65, center.x * 1.25, 0, 0, 0, 0}
	p1 := rmPoint{center.x * 1, center.x * 1, 0, 0, 0, 0}
	p2 := rmPoint{center.x * 1, center.x * 1, 0, 0, 0, 0}
	p3 := rmPoint{center.x * 1.37, center.x * 1.25, 0, 0, 0, 0}

	line := page.AddLine()
	line.AddPoint(center.x*.65, center.x*1.25)
	line.AddPoint(center.x*1.37, center.x*1.25)

	page.DrawBezierCurve(p0, p1, p2, p3)

	// Tamaño y padding de los cuadrados
	squareSize := float32(200)
	padding := float32(50) // Espacio entre cuadrados

	// Generar 10 cuadrados con padding entre ellos
	for i := 0; i < 4; i++ {
		x := float32(100 + i*(int(squareSize)+int(padding)))
		y := float32(100)
		page.DrawFilledRectangle(x, y, x+squareSize, y+squareSize)
	}
	err = page.Export()
	if err != nil {
		fmt.Println("Error exporting page:", err)
		return
	}

	fmt.Println("File testRemarkablePageSmiley.rm generated successfully.")
}
