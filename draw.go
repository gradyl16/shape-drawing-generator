package main

import (
	"errors"
	"fmt"
	"os"
)

// TYPES

// Point struct to represent a point in 2D space
type Point struct {
	x, y int
}

// Rectangle struct representing a rectangle
type Rectangle struct {
	ll Point
	ur Point
	c  Color
}

// Circle struct representing a circle
type Circle struct {
	cp Point
	r  int
	c  Color
}

// Triangle struct representing a triangle
type Triangle struct {
	pt0, pt1, pt2 Point
	c             Color
}

// Screen interface
type screen interface {
	initialize(maxX, maxY int)
	getMaxXY() (int, int)
	drawPixel(x, y int, c Color) error
	getPixel(x, y int) (Color, error)
	clearScreen()
	screenShot(f string) error
}

// Geometry interface
type geometry interface {
	draw(scn screen) error
	shape() string
}

// Display struct implementing the screen interface
type Display struct {
	maxX, maxY int
	matrix     [][]Color
}

type RGB struct {
	R int
	G int
	B int
}

type Color int

// GLOBAL CONSTANTS

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

// Error constants
var outOfBoundsErr = errors.New("geometry out of bounds")
var colorUnknownErr = errors.New("color unknown")

// HELPER FUNCTIONS

// Function to check if a color is valid
func colorUnknown(c Color) bool {
	return !(c == red || c == green || c == blue || c == yellow ||
		c == orange || c == purple || c == brown || c == black || c == white)
}

// Function to check if a point is out of bounds
func outOfBounds(pt Point, scn screen) bool {
	maxX, maxY := scn.getMaxXY()
	return pt.x < 0 || pt.x >= maxX || pt.y < 0 || pt.y >= maxY
}

// INTERFACE IMPLEMENTATIONS

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

// Function to draw a rectangle on the screen
func (rect Rectangle) draw(scn screen) (err error) {
	if outOfBounds(rect.ll, scn) || outOfBounds(rect.ur, scn) {
		return outOfBoundsErr
	}
	if colorUnknown(rect.c) {
		return colorUnknownErr
	}

	for x := rect.ll.x; x <= rect.ur.x; x++ {
		for y := rect.ll.y; y <= rect.ur.y; y++ {
			scn.drawPixel(x, y, rect.c)
		}
	}
	return nil
}

// https://stackoverflow.com/questions/51626905/drawing-circles-with-two-radius-in-golang
// Function to draw a filled circle with rings on the screen
func (circ Circle) draw(scn screen) error {
	if outOfBounds(circ.cp, scn) {
		return outOfBoundsErr
	}
	if colorUnknown(circ.c) {
		return colorUnknownErr
	}

	x0, y0, r := circ.cp.x, circ.cp.y, circ.r

	// Draw rings for each possible radius size
	for ringRadius := 0; ringRadius <= r; ringRadius++ {
		x := ringRadius
		y := 0
		err := 0

		for x >= y {
			// Draw points in all octants to form a ring
			scn.drawPixel(x0+x, y0-y, circ.c)
			scn.drawPixel(x0+y, y0-x, circ.c)
			scn.drawPixel(x0-y, y0-x, circ.c)
			scn.drawPixel(x0-x, y0-y, circ.c)
			scn.drawPixel(x0-x, y0+y, circ.c)
			scn.drawPixel(x0-y, y0+x, circ.c)
			scn.drawPixel(x0+y, y0+x, circ.c)
			scn.drawPixel(x0+x, y0+y, circ.c)

			if err <= 0 {
				y++
				err += 2*y + 1
			}
			if err > 0 {
				x--
				err -= 2*x + 1
			}
		}
	}
	return nil
}

// https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func (tri Triangle) draw(scn screen) (err error) {
	if outOfBounds(tri.pt0, scn) || outOfBounds(tri.pt1, scn) || outOfBounds(tri.pt2, scn) {
		return outOfBoundsErr
	}
	if colorUnknown(tri.c) {
		return colorUnknownErr
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

// Function to initialize the display
func (d *Display) initialize(maxX, maxY int) {
	d.maxX = maxX
	d.maxY = maxY
	d.matrix = make([][]Color, maxY)
	for i := range d.matrix {
		d.matrix[i] = make([]Color, maxX)
	}
	d.clearScreen()
}

// Function to get the maximum dimensions of the screen
func (d *Display) getMaxXY() (int, int) {
	return d.maxX, d.maxY
}

// Function to draw a pixel on the screen
func (d *Display) drawPixel(x, y int, c Color) error {
	if x < 0 || x >= d.maxX || y < 0 || y >= d.maxY {
		return outOfBoundsErr
	}
	d.matrix[y][x] = c
	return nil
}

// Function to get the color of a pixel on the screen
func (d *Display) getPixel(x, y int) (Color, error) {
	if x < 0 || x >= d.maxX || y < 0 || y >= d.maxY {
		return -1, outOfBoundsErr
	}
	return d.matrix[y][x], nil
}

// Function to clear the whole screen
func (d *Display) clearScreen() {
	for i := range d.matrix {
		for j := range d.matrix[i] {
			d.matrix[i][j] = white
		}
	}
}

// Function to take a screenshot and save it as a ppm file
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
