package main

import (
	"fmt"
)

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

func Details(g Geometry) {
	fmt.Println("Perimeter:", g.Perimeter())
	fmt.Println("Area:", g.Area())
}

func main() {
	s := Square{5}
	Details(s)
}
