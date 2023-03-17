package main

import "fmt"

func testDefer(val int) {
	fmt.Print(1)
	defer fmt.Print(2)
	fmt.Print(3)

	if val == 5 {
		return
	}

	defer fmt.Print(4)
}

func main() {
	testDefer(3)
	testDefer(5)
}
