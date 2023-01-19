package main

import "fmt"

func main() {
	code := "fr"
	var country string
	switch code {
	case "fr":
		country = "France"
	case "uk":
		country = "United Kingdom"
	default:
		country = "India"
	}
	fmt.Println(country)

	number := 5
	switch {
	case number%2 == 0:
		fmt.Println("Even Number")
	case number%2 == 1:
		fmt.Println("Odd Number")
	default:
		fmt.Println("Invalid Number")
	}
}
