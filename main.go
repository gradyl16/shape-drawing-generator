package main

import (
	"fmt"
)

func main() {
	fmt.Println("starting ...")

	// Create a Display instance
	var display Display
	display.initialize(1024, 1024)

	// Draw a rectangle
	rect := Rectangle{Point{100, 300}, Point{600, 900}, red}
	err := rect.draw(&display)
	if err != nil {
		fmt.Println("rect: ", err)
	}

	// Draw another rectangle
	rect2 := Rectangle{Point{0, 0}, Point{100, 1024}, green}
	err = rect2.draw(&display)
	if err != nil {
		fmt.Println("rect2: ", err)
	}

	// Draw a circle
	circ := Circle{Point{500, 500}, 200, blue}
	err = circ.draw(&display)
	if err != nil {
		fmt.Println("circ: ", err)
	}

	// Draw a triangle
	tri := Triangle{Point{100, 100}, Point{600, 300}, Point{859, 850}, yellow}
	err = tri.draw(&display)
	if err != nil {
		fmt.Println("tri: ", err)
	}

	// Save the screenshot
	display.screenShot("output")
}
