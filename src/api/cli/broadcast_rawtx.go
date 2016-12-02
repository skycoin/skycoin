package cli

import (
	"errors"
	"fmt"

	"bytes"
	"encoding/json"

	"github.com/skycoin/skycoin/src/api/webrpc"
	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:      "broadcastTransaction",
		Usage:     "Broadcast a raw transaction to the network",
		ArgsUsage: "[raw transaction]",
		Action: func(c *gcli.Context) error {
			rawtx := c.Args().First()
			if rawtx == "" {
				return errors.New("raw transaction is empty")
			}
			txid, err := broadcastTx(rawtx)
			if err != nil {
				return err
			}

			fmt.Println(txid)
			return nil
		},
	}
	Commands = append(Commands, cmd)
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
		return "", fmt.Errorf("rpc request failed, %v", rsp.Error)
	}

	var rlt webrpc.InjectResult
	if err := json.NewDecoder(bytes.NewBuffer(rsp.Result)).Decode(&rlt); err != nil {
		return "", fmt.Errorf("decode inject result failed")
	}

	return rlt.Txid, nil
}
