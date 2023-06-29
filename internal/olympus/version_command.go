package olympus

import "fmt"

type VersionCommand struct{}

func (c *VersionCommand) Execute([]string) error {
	fmt.Println(OLYMPUS_VERSION)
	return nil
}

func init() {
	parser.AddCommand("version",
		"prints olympus version on stdout.",
		"prints olympus version on stdout and exit",
		&VersionCommand{})
}
