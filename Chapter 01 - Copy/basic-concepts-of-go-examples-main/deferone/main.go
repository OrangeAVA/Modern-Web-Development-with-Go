package main

import "fmt"

func main() {
	fmt.Print(1)
	defer fmt.Print(2)
	fmt.Print(3)
	defer fmt.Print(4)
}
