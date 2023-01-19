package main

import (
	"fmt"
	"math"
)

type Circle struct {
	r float64
}

func (c Circle) Perimeter() float64 {
	return 2 * c.r * math.Pi
}

func (c Circle) Area() float64 {
	return c.r * c.r * math.Pi
}

func main() {
	c := Circle{5.0}
	fmt.Println(c.Perimeter())
	fmt.Println(c.Area())
}
