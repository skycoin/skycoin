package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/util"

	"encoding/json"

	"github.com/skycoin/skycoin/src/api/webrpc"
	gcli "github.com/urfave/cli"
)

// Commands all cmds that we support
var Commands []gcli.Command
var (
	rpcAddress        = os.Getenv("SKYCOIN_RPC_ADDR")
	walletDir         = os.Getenv("SKYCOIN_WLT_DIR")
	walletExt         = ".wlt"
	defaultWalletName = "skycoin_cli.wlt"
)

var (
	errConnectNodeFailed = errors.New("connect to node failed")
	errWalletName        = fmt.Errorf("error wallet file name, must has %v extension", walletExt)
	errLoadWallet        = errors.New("load wallet failed")
	errAddress           = errors.New("invalidate address")
	errReadResponse      = errors.New("read response body failed")
	errJSONMarshal       = errors.New("json marshal failed")
	errJSONUnmarshal     = errors.New("json unmarshal failed")
)

func stringPtr(v string) *string {
	return &v
}

func httpGet(url string, v interface{}) error {
	return nil
}

func init() {
	if rpcAddress == "" {
		rpcAddress = "127.0.0.1:6422"
	}

	if walletDir == "" {
		home := util.UserHome()
		walletDir = home + "/.skycoin/wallets/"
	}
}

func getUnspent(addrs []string) ([]unspentOut, error) {
	req, err := webrpc.NewRequest("get_outputs", addrs, "1")
	if err != nil {
		return []unspentOut{}, fmt.Errorf("create webrpc request failed:%v", err)
	}

	rsp, err := webrpc.Do(req, rpcAddress)
	if err != nil {
		return []unspentOut{}, fmt.Errorf("do rpc request failed:%v", err)
	}
	var rlt webrpc.OutputsResult
	if err := json.NewDecoder(bytes.NewBuffer(rsp.Result)).Decode(&rlt); err != nil {
		return nil, errJSONUnmarshal
	}

	ret := make([]unspentOut, len(rlt.Outputs))
	for i, o := range rlt.Outputs {
		ret[i] = unspentOut{
			Hash:              o.Hash,
			SourceTransaction: o.SourceTransaction,
			Address:           o.Address,
			Coins:             o.Coins,
			Hours:             o.Hours,
		}
	}

	return ret, nil
}
