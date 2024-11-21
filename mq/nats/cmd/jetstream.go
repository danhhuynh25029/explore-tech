package main

import (
	"context"
	"fmt"
	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"os"
	"time"
)

func main() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, _ := nats.Connect(url)

	defer nc.Drain()

	js, _ := jetstream.New(nc)

	cfg := jetstream.StreamConfig{
		Name:     "EVENTS",
		Subjects: []string{"events.>"},
	}

	cfg.Storage = jetstream.FileStorage

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, _ = js.CreateStream(ctx, cfg)
	fmt.Println("created the stream")

	js.Publish(ctx, "ORDERS.page_loaded", []byte("hello page_loaded"))
	js.Publish(ctx, "ORDERS.mouse_clicked", []byte("hello mouse_clicked"))
	js.Publish(ctx, "ORDERS.mouse_clicked", []byte("hello mouse_clicked"))
	js.Publish(ctx, "ORDERS.page_loaded", []byte("hello page_loaded"))
	js.Publish(ctx, "ORDERS.mouse_clicked", []byte("hello mouse_clicked"))
	js.Publish(ctx, "ORDERS.input_focused", []byte("hello input_focused"))
	fmt.Println(ctx, "published 6 messages")

	js.Publish(ctx, "events.page_loaded", []byte("hello page_loaded"))
	js.Publish(ctx, "events.mouse_clicked", []byte("hello mouse_clicked"))
	js.Publish(ctx, "events.mouse_clicked", []byte("hello mouse_clicked"))
	js.Publish(ctx, "events.page_loaded", []byte("hello page_loaded"))
	js.Publish(ctx, "events.mouse_clicked", []byte("hello mouse_clicked"))
	js.Publish(ctx, "events.input_focused", []byte("hello input_focused"))
	fmt.Println(ctx, "published 6 messages")
	select {
	case <-js.PublishAsyncComplete():
		fmt.Println("published 6 messages")
	case <-time.After(time.Second):
		log.Fatal("publish took too long")
	}

	// getMessageFromStream

	getMessageFromStream()
}

func getMessageFromStream() {
	url := "nats://127.0.0.1:4222" // default url

	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("cannot connect nats %v", err)
	}
	mgr, _ := jsm.New(nc, jsm.WithTimeout(10*time.Second))
	str, err := mgr.LoadStream("EVENTS")
	if err != nil {
		log.Fatalf("cannot load stream %v", err)
	}

	pops := []jsm.PagerOption{
		jsm.PagerSize(3),
	}

	pgr, err := str.PageContents(pops...)
	if err != nil {
		log.Fatal(err)
	}
	defer pgr.Close()
	for {
		msg, last, err := pgr.NextMsg(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		meta, err := jsm.ParseJSMsgMetadata(msg)
		if err == nil {
			fmt.Printf("[%d] Subject: %s Received: %s\n", meta.StreamSequence(), msg.Subject, meta.TimeStamp().Format(time.RFC3339))
		} else {
			fmt.Printf("Subject: %s Reply: %s\n", msg.Subject, msg.Reply)
		}
		fmt.Println(string(msg.Data))
		if last {
			return
		}
	}
}
