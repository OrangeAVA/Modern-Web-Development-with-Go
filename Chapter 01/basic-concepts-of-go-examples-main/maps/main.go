package main

import "fmt"

func main() {
	var countryMap = make(map[string]string)
	countryMap["fr"] = "France"
	country := countryMap["fr"]
	fmt.Println("Country in map is:", country)
	delete(countryMap, "fr")
	if _, ok := countryMap["fr"]; ok {
		fmt.Println("Country is still in map")
	} else {
		fmt.Println("Country is not in map")
	}
}
