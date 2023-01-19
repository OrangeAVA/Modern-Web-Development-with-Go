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

func IsNil(a interface{}) bool {
	return a == nil
}

func main() {
	var s *Square
	fmt.Println(IsNil(s))
}
