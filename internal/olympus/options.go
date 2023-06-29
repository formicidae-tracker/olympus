package olympus

import "github.com/jessevdk/go-flags"

type Options struct {
}

var opts = &Options{}

var parser = flags.NewParser(opts, flags.Default)

func init() {
	parser.SubcommandsOptional = true
}
