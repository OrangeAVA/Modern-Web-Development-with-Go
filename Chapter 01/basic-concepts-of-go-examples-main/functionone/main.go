package main

import "fmt"

func pi() float64 {
	return 3.14159
}

func inc(a int) int {
	return a + 1
}

func add(a int, b int) int {
	return a + b
}

func main() {
	a := 1
	b := 5

	fmt.Println(pi())
	fmt.Println(inc(a))
	fmt.Println(add(a, b))
}
