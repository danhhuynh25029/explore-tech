package main

import "fmt"

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// using hyperloglog for problem count-distinct
func main() {
	// Make redis clusters
	cluster := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if cluster == nil {
		fmt.Println("Error when creating cluster client")
		return
	}
	cids := []string{"campaign_id1"}
	pipe := cluster.Pipeline()
	// Add new set of user to hll data structure with key campaign_id
	for _, cid := range cids {
		pipe.PFAdd(context.TODO(), cid, "user1", "user2", "user3", "user4", "user1", "user2", "user3")
	}
	_, err := pipe.Exec(context.TODO())
	if err != nil {
		fmt.Println("Error when executing command", err)
		return
	}
	res := cluster.PFCount(context.TODO(), cids...)
	v, err := res.Uint64()
	if err != nil {
		fmt.Println("Error when executing command", err)
		return
	}
	fmt.Println("Cardinality value: ", v)
}
