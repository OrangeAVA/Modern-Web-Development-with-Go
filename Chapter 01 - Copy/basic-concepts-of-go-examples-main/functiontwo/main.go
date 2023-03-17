package main

import "fmt"

func add(a, b int) int {
	return a + b
}

func swap(a, b int) (int, int) {
	return b, a
}

func main() {
	a := 1
	b := 5
	fmt.Println(add(a, b))
	fmt.Println(swap(a, b))
}
