package main

import "fmt"

func modify(s []int) {
	s[0] = 4
	s = append(s, 5)
	fmt.Println(s)
}

func main() {
	s := []int{1, 2, 3}
	modify(s)
	fmt.Println(s)
}
