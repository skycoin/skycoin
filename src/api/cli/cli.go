package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/skycoin/skycoin/src/util"

	gcli "github.com/urfave/cli"
)

// Commands all cmds that we support
var Commands []gcli.Command
var (
	nodeAddress       = os.Getenv("SKYCOIN_NODE_ADDR")
	walletDir         = os.Getenv("SKYCOIN_WLT_DIR")
	walletExt         = ".wlt"
	defaultWalletName = "skycoin_cli.wlt"
)

var (
	errConnectNodeFailed = errors.New("connect to node failed")
	errWalletName        = fmt.Errorf("error wallet file name, must has %v extension", walletExt)
	errLoadWallet        = errors.New("load wallet failed")
	errAddress           = errors.New("error address")
	errReadResponse      = errors.New("read response body failed")
	errJSONMarshal       = errors.New("json marshal failed")
)

func stringPtr(v string) *string {
	return &v
}

func httpGet(url string, v interface{}) error {
	return nil
}

func init() {
	if nodeAddress == "" {
		nodeAddress = "http://localhost:6420"
	}

	if walletDir == "" {
		home := util.UserHome()
		walletDir = home + "/.skycoin-cli/wallet/"
	}
}

func getUnspent(addrs []string) ([]unspentOut, error) {
	url := fmt.Sprintf("%v/outputs?addrs=%s", nodeAddress, strings.Join(addrs, ","))
	rsp, err := http.Get(url)
	if err != nil {
		return []unspentOut{}, errConnectNodeFailed
	}
	defer rsp.Body.Close()
	outs := []unspentOut{}
	if err := json.NewDecoder(rsp.Body).Decode(&outs); err != nil {
		return []unspentOut{}, errors.New("decode json failed")
	}
	return outs, nil
}
