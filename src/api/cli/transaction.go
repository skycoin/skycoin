package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"

	gcli "gopkg.in/urfave/cli.v1"
)

func init() {
	cmd := gcli.Command{
		Name:      "transaction",
		Usage:     "Lists details of specific transaction",
		ArgsUsage: "[option] [transaction id]",
		// Flags:     []gcli.Flag{
		// gcli.StringFlag{
		// 	Name:  "j,json",
		// 	Usage: "Returns the results in JSON format.",
		// },
		// },
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

			url := fmt.Sprintf("%v/transaction?txid=%v", nodeAddress, txid)
			rsp, err := http.Get(url)
			if err != nil {
				return errConnectNodeFailed
			}

			defer rsp.Body.Close()
			d, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				return err
			}
			fmt.Println(string(d))
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
