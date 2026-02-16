package main

import (
	"fmt"
	"math"
)

// Shape is an interface that defines a behavior: calculating area.
// Any type that has a method `Area() float64` will implicitly
// satisfy this interface.
type Shape interface {
	Area() float64
}

// Rectangle is a struct holding data for a rectangle.
type Rectangle struct {
	Width  float64
	Height float64
}

// Area is a method on Rectangle. Because it matches the signature
// required by the Shape interface, Rectangle now "is a" Shape.
// We did not need to write "type Rectangle struct implements Shape".
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Circle is a struct holding data for a circle.
type Circle struct {
	Radius float64
}

// Area is a method on Circle. Circle also satisfies the Shape interface.
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// PrintArea takes any Shape as an argument.
// It doesn't know or care if the concrete type is a Rectangle or Circle.
// It only cares that the value passed in has an Area() method.
// This is polymorphism in Go.
func PrintArea(s Shape) {
	fmt.Printf("Area of shape is: %0.2f\n", s.Area())
}

func main() {
	rect := Rectangle{Width: 10, Height: 5}
	circ := Circle{Radius: 7}

	// We can pass both rect and circ to PrintArea because both
	// types satisfy the Shape interface.
	PrintArea(rect)
	PrintArea(circ)
}
