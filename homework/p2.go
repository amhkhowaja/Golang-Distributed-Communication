package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var ints = []int{1, 2, 3, -4, 5}
	var intMap = make(map[int]string)
	var mut sync.Mutex
	for _, i := range ints {
		wg.Add(1)
		mut.Lock()
		go func() {
			intMap[i] = strconv.Itoa(handleInt(i))
			mut.Unlock()
			wg.Done()
		}()
		wg.Wait()
	}

	fmt.Printf("%v ", intMap)
}
func handleInt(n int) int {
	return (n + 2) * 3
}
