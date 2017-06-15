package cli

import (
	"fmt"

	"bytes"
	"encoding/json"

	"github.com/skycoin/skycoin/src/api/webrpc"
	gcli "github.com/urfave/cli"
)

func broadcastTxCMD() gcli.Command {
	name := "broadcastTransaction"
	return gcli.Command{
		Name:         name,
		Usage:        "Broadcast a raw transaction to the network",
		ArgsUsage:    "[raw transaction]",
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			rawtx := c.Args().First()
			if rawtx == "" {
				gcli.ShowSubcommandHelp(c)
				return nil
			}
			txid, err := broadcastTx(rawtx)
			if err != nil {
				return err
			}

			fmt.Println(txid)
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}

func broadcastTx(rawtx string) (string, error) {
	params := []string{rawtx}
	req, err := webrpc.NewRequest("inject_transaction", params, "1")
	if err != nil {
		return "", fmt.Errorf("create rpc request failed, %v", err)
	}

	rsp, err := webrpc.Do(req, cfg.RPCAddress)
	if err != nil {
		return "", fmt.Errorf("do rpc request failed, %v", err)
	}

	if rsp.Error != nil {
		return "", fmt.Errorf("rpc request failed, %+v", *rsp.Error)
	}

	var rlt webrpc.TxIDJson
	if err := json.NewDecoder(bytes.NewBuffer(rsp.Result)).Decode(&rlt); err != nil {
		return "", fmt.Errorf("decode inject result failed")
	}

	return rlt.Txid, nil
}
