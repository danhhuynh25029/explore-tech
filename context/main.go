package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

func main() {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	//for {
	//	select {
	//	case <-ctx.Done():
	//		fmt.Println("timeout")
	//		return
	//	default:
	//		fmt.Println("not timeout")
	//		time.Sleep(1 * time.Second)
	//	}
	//}
	//cancel()

	Value()
}

const keyID = "id"

func Value() {
	// create a request ID as a random number
	rand.Seed(time.Now().Unix())
	requestID := rand.Intn(1000)

	// create a new context variable with a key value pair
	ctx := context.WithValue(context.Background(), keyID, requestID)
	operation1(ctx)
	operation2(ctx)
}

func operation1(ctx context.Context) {
	// do some work

	// we can get the value from the context by passing in the key
	ctx = context.WithValue(ctx, keyID, 1111)
	ctx = context.WithValue(ctx, 1111, 2222)
	log.Println("operation1 for id:", ctx.Value(1111), " completed")
	log.Println("operation1 for id:", ctx.Value(keyID), " completed")
}

func operation2(ctx context.Context) {
	// do some work
	log.Println("operation2 for id:", ctx.Value(1111), " completed")
	// this way, the same ID is passed from one function call to the next
	log.Println("operation2 for id:", ctx.Value(keyID), " completed")
}
