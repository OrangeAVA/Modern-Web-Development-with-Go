package main

import "fmt"

func calc() func() int {
	a := 0
	return func() int {
		a += 1
		return a
	}
}

func main() {
	res := calc()
	fmt.Println(res())
}
