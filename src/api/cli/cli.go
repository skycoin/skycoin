package cli

import gcli "gopkg.in/urfave/cli.v1"

// Commands all cmds that we support
var Commands []gcli.Command

func stringPtr(v string) *string {
	return &v
}
