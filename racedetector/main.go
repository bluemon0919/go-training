package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	F()
}

func F() {

	var m sync.RWMutex
	var n int

	go func() {
		time.Sleep(50 * time.Millisecond)
		m.Lock()
		n++
		m.Unlock()
		time.Sleep(100 * time.Millisecond)
	}()

	time.Sleep(100 * time.Millisecond)

	m.RLock()
	fmt.Println(n)
	m.RUnlock()

	time.Sleep(100 * time.Millisecond)
}
