package main

import (
	"log"

	"github.com/formicidae-tracker/olympus/internal/olympus"
)

func main() {
	if err := olympus.Execute(); err != nil {
		log.Fatalf("%s", err)
	}
}
