package main

import (
  "fmt"
  "math"
  "os"
  "errors"
)


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
  return
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
