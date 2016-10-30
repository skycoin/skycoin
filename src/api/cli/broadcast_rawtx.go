package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:      "broadcastTransaction",
		ArgsUsage: "Broadcast a raw transaction to the network.",
		Usage:     "[raw transaction]",
		Action: func(c *gcli.Context) error {
			rawtx := c.Args().First()
			if rawtx == "" {
				return errors.New("raw transaction is empty")
			}

			v, err := broadcastTx(rawtx)
			if err != nil {
				return err
			}
			fmt.Println(v)
			return nil
		},
	}
	Commands = append(Commands, cmd)
}

func broadcastTx(rawtx string) (string, error) {
	var tx = struct {
		Rawtx string `json:"rawtx"`
	}{
		rawtx,
	}
	d, err := json.Marshal(tx)
	if err != nil {
		return "", errors.New("error raw transaction")
	}
	url := fmt.Sprintf("http://%s/injectTransaction", rpcAddress)
	rsp, err := http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return "", errConnectNodeFailed
	}
	defer rsp.Body.Close()
	v, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", errReadResponse
	}

	return strings.Trim(string(v), "\""), nil
}
