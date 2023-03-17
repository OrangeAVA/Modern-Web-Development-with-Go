package main

import "fmt"

func double(a int) {
	a = a * 2
	fmt.Println(a)
}

func doublePointer(a *int) {
	*a = *a * 2
	fmt.Println(*a)
}

func main() {
	a := 5
	double(a)
	fmt.Println(a)

	doublePointer(&a)
	fmt.Println(a)
}
