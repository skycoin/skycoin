package cli

import gcli "github.com/urfave/cli"

// Commands all cmds that we support
var Commands []gcli.Command

func stringPtr(v string) *string {
	return &v
}
