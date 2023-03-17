package main

import "fmt"

func calc(fn func(int, int) int) int {
	return fn(7, 18)
}

func main() {
	add := func(a, b int) int {
		return a + b
	}

	fmt.Println(calc(add))
}
