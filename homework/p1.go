package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var receivedCount uint64
	var wg sync.WaitGroup
	receivedCount = 0
	c := make(chan int)
	wg.Add(2)
	go func() {
		for i := 0; i <= 9; i++ {
			fmt.Printf("g1 sent <%v>\n", i)
			c <- i
		}
		close(c)
		wg.Done()
	}()
	go func() {
		for msg := range c {
			fmt.Printf("g2 received <%v>\n", msg)
			atomic.AddUint64(&receivedCount, 1)
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Printf("Received count:  %v", receivedCount)
}
