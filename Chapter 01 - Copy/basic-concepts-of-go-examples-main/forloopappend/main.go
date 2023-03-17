package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	var b []*int

	for _, v := range a {
		b = append(b, &v)
	}

	for _, v := range b {
		fmt.Println(*v)
	}
}
