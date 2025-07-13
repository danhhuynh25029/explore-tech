package main

import (
	"fmt"
	"sync"
)

var counter int
var rwMutex sync.RWMutex

func writeCounter(value int) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	counter = counter + value
}

func main() {
	var wg sync.WaitGroup

	// Tạo các goroutines để ghi dữ liệu
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			writeCounter(value)
		}(i)
	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}
