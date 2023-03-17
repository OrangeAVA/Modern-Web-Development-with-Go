package main

import "fmt"

func main() {
	var i int = 18
	var pi *int

	pi = &i
	*pi = 27
	fmt.Println(i)
}
