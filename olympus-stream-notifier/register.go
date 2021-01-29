package main

type RegisterCommand struct {
	Args struct {
		Hostname string
		URL      string
	} `positional-args:"yes" required:"yes"`
}

type UnregisterCommand struct {
	Args struct {
		Hostname string
	} `positional-args:"yes" required:"yes"`
}

func (c *RegisterCommand) Execute(args []string) error {
	if err := RegisterPlaylist(opts.Basepath, c.Args.Hostname); err != nil {
		return err
	}
	return RegisterOlympus(c.Args.Hostname, c.Args.URL)
}

func (c *UnregisterCommand) Execute(args []string) error {
	if err := UnregisterPlaylist(opts.Basepath, c.Args.Hostname); err != nil {
		return err
	}
	return UnregisterOlympus(c.Args.Hostname)
}

func init() {
	parser.AddCommand("register", "register an host", "register a new host", &RegisterCommand{})
	parser.AddCommand("unregister", "unregister an host", "unregister an host", &UnregisterCommand{})
}
