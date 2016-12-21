package cli

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"

	"bytes"
	"encoding/json"

	"github.com/skycoin/skycoin/src/api/webrpc"
	gcli "github.com/urfave/cli"
)

func transactionCMD() gcli.Command {
	name := "transaction"
	return gcli.Command{
		Name:         name,
		Usage:        "Show detail info of specific transaction",
		ArgsUsage:    "[transaction id]",
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			txid := c.Args().First()
			if txid == "" {
				return errors.New("txid is empty")
			}

			// validate the txid
			_, err := cipher.SHA256FromHex(txid)
			if err != nil {
				return errors.New("error txid")
			}

			tx, err := getTransactionByID(txid)
			if err != nil {
				return err
			}

			v, err := json.MarshalIndent(tx, "", "    ")
			if err != nil {
				return errors.New("invalid tx result")
			}

			fmt.Println(string(v))
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}

func getTransactionByID(txid string) (*webrpc.TxnResult, error) {
	req, err := webrpc.NewRequest("get_transaction", []string{txid}, "1")
	if err != nil {
		return nil, fmt.Errorf("create rpc request failed:%v", err)
	}

	rsp, err := webrpc.Do(req, cfg.RPCAddress)
	if err != nil {
		return nil, fmt.Errorf("do rpc request failed:%v", err)
	}

	if rsp.Error != nil {
		return nil, fmt.Errorf("do rpc request failed:%+v", *rsp.Error)
	}

	rlt := webrpc.TxnResult{}
	if err := json.NewDecoder(bytes.NewReader(rsp.Result)).Decode(&rlt); err != nil {
		return nil, err
	}
	return &rlt, nil
}
