package cli

import gcli "gopkg.in/urfave/cli.v1"

// Commands all cmds that we support
var Commands []gcli.Command
var skycoinNodeAddr = "http://localhost:6420"

func stringPtr(v string) *string {
	return &v
}

func httpGet(url string, v interface{}) error {
	return nil
}
