package cli

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"

	"github.com/skycoin/skycoin/src/api/webrpc"
	gcli "github.com/urfave/cli"
)

func init() {
	name := "transaction"
	cmd := gcli.Command{
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

			req, err := webrpc.NewRequest("get_transaction", []string{txid}, "1")
			if err != nil {
				return fmt.Errorf("create rpc request failed:%v", err)
			}

			rsp, err := webrpc.Do(req, cfg.RPCAddress)
			if err != nil {
				return fmt.Errorf("do rpc request failed:%v", err)
			}

			if rsp.Error != nil {
				return fmt.Errorf("do rpc request failed:%+v", *rsp.Error)
			}

			fmt.Println(string(rsp.Result))
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
