package main

import (
	"fmt"
	"sync"
)

func main() {
	a := []int{5, 8, 4, 9, 3}
	wg := &sync.WaitGroup{}
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(i int) {
			a[i] *= 3
			fmt.Println("Goroutine ", i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println(a)
}
