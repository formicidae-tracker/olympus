package main

import (
	"os"

	"github.com/formicidae-tracker/olympus/internal/olympus"
)

func main() {
	if err := olympus.Execute(); err != nil {
		os.Exit(3)
	}
}
