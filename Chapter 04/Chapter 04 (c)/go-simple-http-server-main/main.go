package main

import (
	"fmt"
	"net/http"
)

func helloWorld(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello World\n")
}

func main() {
	http.HandleFunc("/hello", helloWorld)

	http.ListenAndServe(":8080", nil)
}
