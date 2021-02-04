package main

import (
	"log"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Basepath   string `long:"basepath" default:"/var/www/olympus"`
	RPCAddress string `long:"rpc-address" default:"localhost"`
	RPCPort    int    `long:"rpc-address" default:"3001"`
}

var opts = &Options{}
var parser = flags.NewParser(opts, flags.Default)

func Execute() error {
	if _, err := parser.Parse(); err != nil {
		if ferr, ok := err.(*flags.Error); ok == true && ferr.Type == flags.ErrHelp {
			return nil
		}
		return err
	}
	return nil
}

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("unhandled error: %s", err)
	}
}
