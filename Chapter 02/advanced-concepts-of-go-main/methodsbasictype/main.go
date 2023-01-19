package main

import (
	"fmt"
)

type CustomInt int

func (ci CustomInt) Double() int {
	return int(ci * 2)
}

func main() {
	c := CustomInt(5)
	fmt.Println(c.Double())
}
