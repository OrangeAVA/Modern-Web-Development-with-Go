package main

import "fmt"

type Node[T any] struct {
	leftChild  *Node[T]
	rightChild *Node[T]
	value      T
}

func main() {
	nodeInt := Node[int]{value: 5}
	nodeFloat := Node[float64]{value: 5.2}

	fmt.Println(nodeInt)
	fmt.Println(nodeFloat)
}
