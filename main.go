package main

import (
  "fmt"
  "math"
  "os"
  "errors"
)

// GLOBAL CONSTANTS

// Color constants
var (
	red    = Color{255, 0, 0}
	green  = Color{0, 255, 0}
	blue   = Color{0, 0, 255}
	yellow = Color{255, 255, 0}
	orange = Color{255, 164, 0}
	purple = Color{128, 0, 128}
	brown  = Color{165, 42, 42}
	black  = Color{0, 0, 0}
	white  = Color{255, 255, 255}
)

// Error constants
var outOfBoundsErr = errors.New("geometry out of bounds")
var colorUnknownErr = errors.New("color unknown")

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

// HELPER FUNCTIONS

// Function to check if a color is valid
func colorUnkown(c Color) bool {
	return !(c == red || c == green || c == blue || c == yellow ||
		c == orange || c == purple || c == brown || c == black || c == white)
}

// Function to check if a point is out of bounds
func outOfBounds(pt Point, scn screen) bool {
	maxX, maxY := scn.getMaxXY()
	return pt.x < 0 || pt.x >= maxX || pt.y < 0 || pt.y >= maxY
}

//  https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func interpolate (l0, d0, l1, d1 int) (values []int) {
  a := float64(d1 - d0) / float64(l1 - l0)
  d  := float64(d0)

  count := l1-l0+1
  for ; count>0; count-- {
    values = append(values, int(d))
    d = d+a
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

// Function to draw a circle on the screen
func (circ Circle) draw(scn screen) (err error) {
	if outOfBounds(circ.cp, scn) {
		return outOfBoundsErr
	}
	if colorUnknown(circ.c) {
		return colorUnknownErr
	}

	x0, y0, r := circ.cp.x, circ.cp.y, circ.r
	x := r
	y := 0
	err := 0

	for x >= y {
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
	return nil
}

//  https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func (tri Triangle) draw(scn screen) (err error) {
  if outOfBounds(tri.pt0,scn) || outOfBounds(tri.pt1,scn)  || outOfBounds(tri.pt2,scn){
    return outOfBoundsErr
  }
  if colorUnknown(tri.c) {
    return colorUnknownErr
  }

  y0 := tri.pt0.y
  y1 := tri.pt1.y
  y2 := tri.pt2.y

  // Sort the points so that y0 <= y1 <= y2
  if y1 < y0 { tri.pt1, tri.pt0 = tri.pt0, tri.pt1 }
  if y2 < y0 { tri.pt2, tri.pt0 = tri.pt0, tri.pt2 }
  if y2 < y1 { tri.pt2, tri.pt1 = tri.pt1, tri.pt2 }

  x0,y0,x1,y1,x2,y2 := tri.pt0.x, tri.pt0.y, tri.pt1.x, tri.pt1.y, tri.pt2.x, tri.pt2.y

  x01 := interpolate(y0, x0, y1, x1)
  x12 := interpolate(y1, x1, y2, x2)
  x02 := interpolate(y0, x0, y2, x2)

  // Concatenate the short sides

  x012 := append(x01[:len(x01)-1],  x12...)

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
  for y := y0; y<= y2; y++  {
    for x := x_left[y - y0]; x <=x_right[y - y0]; x++ {
      scn.drawPixel(x, y, tri.c)
    }
  }
  return nil
}

// display 
// TODO: you must implement the struct for this variable, and the interface it implements (screen)
var display Display


func main() {
  fmt.Println("starting ...")
  display.initialize(1024,1024)

  rect :=  Rectangle{Point{100,300}, Point{600,900}, red}
  err := rect.draw(&display)
  if err != nil {
    fmt.Println("rect: ", err)
  }

  rect2 := Rectangle{Point{0,0}, Point{100, 1024}, green}
  err = rect2.draw(&display)
  if err != nil {
    fmt.Println("rect2: ", err)
  }

  rect3 := Rectangle{Point{0,0}, Point{100, 1022}, 102}
  err = rect3.draw(&display)
  if err != nil {
    fmt.Println("rect3: ", err)
  }

  circ := Circle{Point{500,500}, 200, green}
  err = circ.draw(&display)
  if err != nil {
    fmt.Println("circ: ", err)
  }

  tri := Triangle{Point{100, 100}, Point{600, 300},  Point{859,850}, yellow}
  err = tri.draw(&display)
  if err != nil {
    fmt.Println("tri: ", err)
  }

  display.screenShot("output")
}
