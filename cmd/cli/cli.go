package main

import (
	"fmt"
	"os"

	"strings"

	"github.com/skycoin/skycoin/src/api/cli"
	"github.com/skycoin/skycoin/src/util"
)

var (
	commandHelpTemplate = `USAGE:
		{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Category}}
		
CATEGORY:
		{{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
		{{.Description}}{{end}}{{if .VisibleFlags}}

OPTIONS:
		{{range .VisibleFlags}}{{.}}
		{{end}}{{end}}
	`
)

func main() {
	// get rpc address from env
	rpcAddr := os.Getenv("RPC_ADDR")
	if rpcAddr == "" {
		rpcAddr = "127.0.0.1:6430"
	}

	// get wallet dir from env
	wltDir := os.Getenv("WALLET_DIR")
	if wltDir == "" {
		home := util.UserHome()
		wltDir = home + "/.skycoin/wallets"
	}

	// get wallet name from env
	wltName := os.Getenv("WALLET_NAME")
	if wltName == "" {
		wltName = "skycoin_cli.wlt"
	} else {
		if !strings.HasSuffix(wltName, ".wlt") {
			fmt.Println("value of 'WALLET_NAME' env is not correct, must has .wlt extenstion")
			return
		}
	}

	// init the skycli
	app := cli.NewApp(cli.RPCAddr(rpcAddr),
		cli.WalletDir(wltDir),
		cli.DefaultWltName(wltName))
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
