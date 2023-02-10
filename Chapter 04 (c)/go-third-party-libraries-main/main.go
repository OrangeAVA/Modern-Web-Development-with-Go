package main

import (
	"fmt"

	set "github.com/deckarep/golang-set"
)

func main() {
	yellowFruit := set.NewSet()
	yellowFruit.Add("banana")
	yellowFruit.Add("lemon")
	yellowFruit.Add("pineapple")
	fmt.Println(yellowFruit)

	redFruit := set.NewSetFromSlice([]interface{}{"apple", "cherry", "strawberry"})
	fmt.Println(redFruit)

	fruit := yellowFruit.Union(redFruit)
	fmt.Println(fruit)
}
