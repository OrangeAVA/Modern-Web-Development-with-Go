package main

import "fmt"

func Contains[T comparable](arr []T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

func main() {
	intArr := []int{18, 27, 21, 1}
	stringArr := []string{"apple", "banana", "ananas", "orange"}

	fmt.Println(Contains(intArr, 27))
	fmt.Println(Contains(stringArr, "lemon"))
}
