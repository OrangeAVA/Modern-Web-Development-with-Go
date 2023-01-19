package main

import "fmt"

func main() {
	result := 1
	sum := 1
	for result < 500 {
		result *= sum * 2
	}
	fmt.Println(result)
}
