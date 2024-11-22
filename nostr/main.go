package main

import (
	"context"
	"fmt"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"log"
	"time"
)

func main() {
	sk := nostr.GeneratePrivateKey()
	pk, _ := nostr.GetPublicKey(sk)
	nsec, _ := nip19.EncodePrivateKey(sk)
	npub, _ := nip19.EncodePublicKey(pk)

	fmt.Println("sk:", sk)
	fmt.Println("pk:", pk)
	fmt.Println(nsec)
	fmt.Println(npub)
	pub, err := nostr.GetPublicKey(sk)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pub)
	//publish(sk, pub)

	time.Sleep(5 * time.Second)

	subscribe(npub)
}

func subscribe(npub string) {
	ctx := context.Background()
	for _, v := range []string{"wss://relay.stoner.com", "wss://nostr-pub.wellorder.net"} {
		relay, err := nostr.RelayConnect(ctx, v)
		if err != nil {
			panic(err)
		}

		var filters nostr.Filters
		if _, v, err := nip19.Decode(npub); err == nil {
			pub := v.(string)
			filters = []nostr.Filter{{
				Kinds:   []int{nostr.KindTextNote},
				Authors: []string{pub},
				Limit:   1,
			}}
		} else {
			panic(err)
		}

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		sub, err := relay.Subscribe(ctx, filters)
		if err != nil {
			panic(err)
		}

		for ev := range sub.Events {
			// handle returned event.
			// channel will stay open until the ctx is cancelled (in this case, context timeout)
			fmt.Println(ev.ID, ev.Content)
		}
	}
}

func publish(sk, pub string) {

	ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindTextNote,
		Tags:      nil,
		Content:   "Hello World!",
	}

	// calling Sign sets the event ID field and the event Sig field
	ev.Sign(sk)

	// publish the event to two relays
	ctx := context.Background()
	for _, url := range []string{"wss://relay.stoner.com", "wss://nostr-pub.wellorder.net"} {
		relay, err := nostr.RelayConnect(ctx, url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := relay.Publish(ctx, ev); err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("published to %s\n", url)
	}
}
