package main

import "fmt"

func sendMessage(message string, ch chan string) {
	ch <- message
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	go sendMessage("Channel 1", ch1)

	for {
		select {
		case <-ch1:
			fmt.Println("Channel One")
			return
		case <-ch2:
			fmt.Println("Channel Two")
		default:
			fmt.Println("Waiting")
		}
	}
}
