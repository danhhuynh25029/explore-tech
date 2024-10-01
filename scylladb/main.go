package main

import (
	"context"
	"fmt"
	"github.com/scylladb/scylla-go-driver"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	cfg := scylla.DefaultSessionConfig("pets_clinic", "localhost")
	session, err := scylla.NewSession(ctx, cfg)
	if err != nil {
		log.Fatalf("func scylla.NewSession err :  %v", err)
	}
	defer session.Close()

	requestCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	q, err := session.Prepare(requestCtx, "SELECT book_id, book_name FROM pets_clinic.books WHERE book_id=?")
	if err != nil {
		log.Fatalf("func session.Prepare err : %v ", err)
	}

	res, err := q.BindInt64(0, 64).Exec(requestCtx)
	fmt.Println(res)
}
