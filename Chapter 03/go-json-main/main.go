package main

import (
	"encoding/json"
	"fmt"
)

type spouse struct {
	Name string
	Age  int
}

type person struct {
	Name     string
	Age      int
	Married  bool
	Spouse   spouse
	Children []string
}

func main() {
	jsonString := `{"name": "Albert Einstein", "age": 76, "married": true, "spouse": {"name": "Mileva", "age": 72}, "children": ["Lieserl", "Hans Albert", "Eduard"]}`
	var p person
	err := json.Unmarshal([]byte(jsonString), &p)
	if err != nil {
		fmt.Println("Failed to decode JSON")
	}
	fmt.Println(p)

	jsonByte, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Failed to encode JSON")
	}
	fmt.Println(string(jsonByte))
}
