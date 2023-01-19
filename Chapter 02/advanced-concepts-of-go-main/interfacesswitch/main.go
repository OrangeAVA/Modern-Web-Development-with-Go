package main

import "fmt"

type Geometry interface {
	Perimeter() int
	Area() int
}

type Square struct {
	a int
}

func (s Square) Perimeter() int {
	return s.a * 4
}
func (s Square) Area() int {
	return s.a * s.a
}

type Rectangle struct {
	a, b int
}

func (r Rectangle) Perimeter() int {
	return r.a*2 + r.b*2
}
func (r Rectangle) Area() int {
	return r.a * r.b
}

func main() {
	var g Geometry = Square{5}

	switch g.(type) {
	case Square:
		fmt.Printf("Square")
	case Rectangle:
		fmt.Printf("Rectangle")
	default:
		fmt.Printf("Unknown type")
	}
}
