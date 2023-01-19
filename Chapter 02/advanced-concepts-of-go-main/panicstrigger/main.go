package main

func checkName(name string) {
	if name == "" {
		panic("Invalid name!!!")
	}
}

func main() {
	checkName("test")
	checkName("")
}
