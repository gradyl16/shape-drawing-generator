// ## IMPORTS ##

package main

import (
	"errors"
	"fmt"
	"os"
)


// ## TYPE DEFINITIONS ##

type Color int

type RGB struct {
	R Color
	G Color
	B Color
}

// Represents a coordinate pair in 2D space
type Point struct {
	x, y int
}

type Rectangle struct {
	ll Point  // Lower left corner
	ur Point  // Upper right corner
	c  Color
}

type Circle struct {
	cp Point  // Center
	r  int    // Radius
	c  Color
}

type Triangle struct {
	pt0, pt1, pt2 Point  // Vertices
	c             Color
}

// Represents the screen
type Display struct {
	maxX, maxY int
	matrix     [][]Color
}


// ## INTERFACE DEFINITIONS ##

type screen interface {
	initialize(maxX, maxY int)          // Initializes with max dimensions
	getMaxXY() (int, int)               // Returns max dimensions 
	drawPixel(x, y int, c Color) error  // Draws a pixel on the screen
	getPixel(x, y int) (Color, error)   // Returns the color of a pixel
	clearScreen()                       // Sets all pixels to white
	screenShot(f string) error          // Generates .ppm image of screen
}

type geometry interface {
	draw(scn screen) error
	// shape() string --> TODO: Implement! :)
}


// ## GLOBAL CONSTANTS ##

const (
	red Color = iota
	green
	blue
	yellow
	orange
	purple
	brown
	black
	white
)

var cmap = [...]RGB{
	{255, 0, 0},     // red
	{0, 255, 0},     // green
	{0, 0, 255},     // blue
	{255, 255, 0},   // yellow
	{255, 164, 0},   // orange
	{128, 0, 128},   // purple
	{165, 42, 42},   // brown
	{0, 0, 0},       // black
	{255, 255, 255}, // white
}

var errOutOfBounds = errors.New("geometry out of bounds")
var errColorUnknown = errors.New("color unknown")


// ## INTERFACE IMPLEMENTATIONS ##
// HELPER FUNCTIONS

func isColorUnknown(c Color) bool {
	return !(c == red || c == green || c == blue || c == yellow ||
		c == orange || c == purple || c == brown || c == black || c == white)
}

func isOutOfBounds(pt Point, scn screen) bool {
	maxX, maxY := scn.getMaxXY()
	return pt.x < 0 || pt.x >= maxX || pt.y < 0 || pt.y >= maxY
}

func isInsideCircle(center Point, tile Point, radius int) bool {
	dx := center.x - tile.x
	dy := center.y - tile.y
	distanceSquared := dx*dx + dy*dy
	return distanceSquared <= radius*radius
}

// SCREEN

func (d *Display) getMaxXY() (int, int) {
	return d.maxX, d.maxY
}

func (d *Display) getPixel(x, y int) (Color, error) {
	if x < 0 || x >= d.maxX || y < 0 || y >= d.maxY {
		return -1, errOutOfBounds
	}
	return d.matrix[y][x], nil
}

func (d *Display) drawPixel(x, y int, c Color) error {
	if x < 0 || x >= d.maxX || y < 0 || y >= d.maxY {
		return errOutOfBounds
	}
	d.matrix[y][x] = c
	return nil
}

func (d *Display) clearScreen() {
	for i := range d.matrix {
		for j := range d.matrix[i] {
			d.matrix[i][j] = white
		}
	}
}

func (d *Display) initialize(maxX, maxY int) {
	d.maxX = maxX
	d.maxY = maxY
	d.matrix = make([][]Color, maxY)
	for i := range d.matrix {
		d.matrix[i] = make([]Color, maxX)
	}
	d.clearScreen()
}

func (d *Display) screenShot(f string) error {
	file, err := os.Create(f + ".ppm")
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "P3\n%d %d\n255\n", d.maxX, d.maxY)

	for _, row := range d.matrix {
		for _, color := range row {
			pixel := cmap[color]
			fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B)
		}
		// Add a newline after each row
		fmt.Fprint(file, "\n")
	}
	return nil
}

// GEOMETRY

// Source:
// https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func interpolate(l0, d0, l1, d1 int) (values []int) {
	a := float64(d1-d0) / float64(l1-l0)
	d := float64(d0)

	count := l1 - l0 + 1
	for ; count > 0; count-- {
		values = append(values, int(d))
		d = d + a
	}
	return
}

func (rect Rectangle) draw(scn screen) (err error) {
	if isOutOfBounds(rect.ll, scn) || isOutOfBounds(rect.ur, scn) {
		return errOutOfBounds
	}
	if isColorUnknown(rect.c) {
		return errColorUnknown
	}

	for x := rect.ll.x; x <= rect.ur.x; x++ {
		for y := rect.ll.y; y <= rect.ur.y; y++ {
			scn.drawPixel(x, y, rect.c)
		}
	}
	return nil
}

// Source:
// https://www.redblobgames.com/grids/circle-drawing/
func (circ Circle) draw(scn screen) error {
	if isOutOfBounds(circ.cp, scn) {
		return errOutOfBounds
	}
	if isColorUnknown(circ.c) {
		return errColorUnknown
	}

	x0, y0, r := circ.cp.x, circ.cp.y, circ.r

	// Calculate the bounding box
	top := int(float64(y0 - r))
	bottom := int(float64(y0 + r))
	left := int(float64(x0 - r))
	right := int(float64(x0 + r))

	// Draw the circle within the bounding box
	for y := top; y <= bottom; y++ {
		for x := left; x <= right; x++ {
			tile := Point{x, y}
			if isInsideCircle(circ.cp, tile, circ.r) {
				scn.drawPixel(x, y, circ.c)
			}
		}
	}

	return nil
}

// Source:
// https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func (tri Triangle) draw(scn screen) (err error) {
	if isOutOfBounds(tri.pt0, scn) || isOutOfBounds(tri.pt1, scn) || isOutOfBounds(tri.pt2, scn) {
		return errOutOfBounds
	}
	if isColorUnknown(tri.c) {
		return errColorUnknown
	}

	y0 := tri.pt0.y
	y1 := tri.pt1.y
	y2 := tri.pt2.y

	// Sort the points so that y0 <= y1 <= y2
	if y1 < y0 {
		tri.pt1, tri.pt0 = tri.pt0, tri.pt1
	}
	if y2 < y0 {
		tri.pt2, tri.pt0 = tri.pt0, tri.pt2
	}
	if y2 < y1 {
		tri.pt2, tri.pt1 = tri.pt1, tri.pt2
	}

	x0, y0, x1, y1, x2, y2 := tri.pt0.x, tri.pt0.y, tri.pt1.x, tri.pt1.y, tri.pt2.x, tri.pt2.y

	x01 := interpolate(y0, x0, y1, x1)
	x12 := interpolate(y1, x1, y2, x2)
	x02 := interpolate(y0, x0, y2, x2)

	// Concatenate the short sides

	x012 := append(x01[:len(x01)-1], x12...)

	// Determine which is left and which is right
	var x_left, x_right []int
	m := len(x012) / 2
	if x02[m] < x012[m] {
		x_left = x02
		x_right = x012
	} else {
		x_left = x012
		x_right = x02
	}

	// Draw the horizontal segments
	for y := y0; y <= y2; y++ {
		for x := x_left[y-y0]; x <= x_right[y-y0]; x++ {
			scn.drawPixel(x, y, tri.c)
		}
	}
	return nil
}