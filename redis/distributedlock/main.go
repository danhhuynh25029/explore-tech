package main

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"sync"
)

// Optimistic lock : read heavy,low update
// Pessimistic lock : high update

var count = 0

func incrUseDisLock(i int) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	mutexname := "my-global-mutex"
	mutex := rs.NewMutex(mutexname)

	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	count = count + 1

	if ok, err := mutex.Unlock(); !ok || err != nil {
		panic("unlock failed")
	}
}

func incrNotUseDisLock(i int) {
	count += 1
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			incrUseDisLock(i)
		}()
	}
	wg.Wait()
	fmt.Println(count)
}
