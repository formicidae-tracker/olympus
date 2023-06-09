package main

import (
	"os"
	"path/filepath"
)

var datapath string

func init() {
	datapath = os.Getenv("OLYMPUS_DATA_HOME")
	if len(datapath) == 0 {
		datapath = filepath.Join(os.TempDir(), "fort", "olympus")
	}
}
