package main

import (
	"log"
)

func div(a, b float32) (c float32) {
	log.Println("Input parameters:", a, b)

	if b == 0 {
		log.Println("Division by zero")
		return 0
	}

	c = a / b

	log.Println("Division result:", c)
	return
}

func main() {
	div(6.0, 3.0)
	div(6.0, 0.0)
}
