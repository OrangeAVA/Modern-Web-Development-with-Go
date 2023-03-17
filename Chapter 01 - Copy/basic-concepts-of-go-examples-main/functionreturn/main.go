package main

import "fmt"

func multiply(param int) func(int) int {
	if param%2 == 0 {
		return func(a int) int {
			return a * 2
		}
	} else {
		return func(a int) int {
			return a * 3
		}
	}
}

func main() {
	double := multiply(2)
	triple := multiply(3)
	fmt.Println(double(5), triple(5))
}
