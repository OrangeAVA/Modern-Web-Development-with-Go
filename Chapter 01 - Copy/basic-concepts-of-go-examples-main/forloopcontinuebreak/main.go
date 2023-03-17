package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4, 5}
	for i := 0; i < 5; i++ {
		if i == 3 {
			continue
		}
		fmt.Println(a[i])
	}

	for i := 0; i < 5; i++ {
		if i == 3 {
			break
		}
		fmt.Println(a[i])
	}
}
