package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"

	gcli "gopkg.in/urfave/cli.v1"
)

var defaultAddrNum = 1

func init() {
	cmd := gcli.Command{
		Name:      "generateAddresses",
		Usage:     "Generate additional addresses for a wallet.",
		ArgsUsage: "[options]",
		Description: `
        Use caution when using the “-p” command. If you have command 
        history enabled your wallet encryption password can be recovered from the history log. 
        If you do not include the “-p” option you will be prompted to enter your password after 
        you enter your command.`,
		Flags: []gcli.Flag{
			gcli.IntFlag{
				Name: "m",
				Usage: "[numberOfAddresses]	Number of addresses to generate. By default 1 address is generated.",
			},
			gcli.StringFlag{
				Name:  "w",
				Usage: "[wallet file or path] In wallet. If no path is specified your default wallet path will be used.",
			},
			gcli.StringFlag{
				Name:  "j,json",
				Usage: "Returns the results in JSON format.",
			},
		},
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
