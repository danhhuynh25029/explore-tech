package main

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"log"
	"sync"
	"time"
)

var Cache = make(map[string]string)
var group = singleflight.Group{}

// MultipleRequest
// cache penetration :  call with id is maximum after call to database high IO
// cache breakdown :  not have in cache after call to database high IO

func MultipleRequest(n int) {
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(req int) {
			defer wg.Done()
			log.Printf("resp : %v, req : %v", SingleFlight(req, "1"), req)
		}(i)
	}
	wg.Wait()
}

func SingleFlight(req int, key string) string {

	if data, ok := Cache[key]; ok {
		log.Printf("cache hit req: %v", req)
		return data
	}
	row, err, _ := group.Do(key, func() (interface{}, error) {
		log.Printf("missing cache : %v ,request : %v", key, req)
		Cache[key] = "mickey"
		time.Sleep(100 * time.Millisecond)
		return "mickey", nil
	})
	if err != nil {
		fmt.Println("error : ", err)
		return ""
	}
	return row.(string)
}
