package main

import (
	"fmt"
	"sync"
	"time"
)

type MutualExclusion struct {
	mutex sync.Mutex
	value int
}

func (me *MutualExclusion) Double() {
	me.mutex.Lock()
	me.value *= 2
	me.mutex.Unlock()
}

func main() {
	me := MutualExclusion{value: 5}
	go me.Double()
	go me.Double()

	time.Sleep(time.Second)
	fmt.Println(me.value)
}
