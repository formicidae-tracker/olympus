package main

import (
	"fmt"
	"log"

	"github.com/SherClockHolmes/webpush-go"
)

func main() {
	if err := execute(); err != nil {
		log.Fatalf("unhandled error: %s", err)
	}
}

func execute() error {
	private, public, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		return err
	}

	fmt.Printf("OLYMPUS_VAPID_PRIVATE=%s\n", private)
	fmt.Printf("OLYMPUS_VAPID_PUBLIC=%s\n", public)

	return nil
}
