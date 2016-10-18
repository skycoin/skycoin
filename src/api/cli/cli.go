package cli

import (
	"os"

	"github.com/skycoin/skycoin/src/util"

	gcli "gopkg.in/urfave/cli.v1"
)

// Commands all cmds that we support
var Commands []gcli.Command
var (
	nodeAddress       = os.Getenv("SKYCOIN_NODE_ADDR")
	walletDir         = os.Getenv("SKYCOIN_WLT_DIR")
	defaultWalletName = "skycoin_cli.wlt"
)

func stringPtr(v string) *string {
	return &v
}

func httpGet(url string, v interface{}) error {
	return nil
}

func init() {
	if nodeAddress == "" {
		nodeAddress = "http://localhost:6421"
	}

	if walletDir == "" {
		home := util.UserHome()
		walletDir = home + "/.skycoin-cli/wallet/"
	}
}
