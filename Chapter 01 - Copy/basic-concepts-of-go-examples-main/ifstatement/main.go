package main

import "fmt"

func main() {
	a := 5
	if a < 100 {
		a += 1
	}
	fmt.Println(a)

	if b := a * a; b < 100 {
		a += 1
	}
	fmt.Println(a)

	code := "fr"
	var country string
	if code == "fr" {
		country = "France"
	} else if code == "uk" {
		country = "United Kingdom"
	} else {
		country = "India"
	}
	fmt.Println(country)
}
