package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/urfave/cli"
)

type unspentOut struct {
	Hash              string `json:"txid"` //hash uniquely identifies transaction
	SourceTransaction string `json:"src_tx"`
	Address           string `json:"address"`
	Coins             string `json:"coins"`
	Hours             uint64 `json:"hours"`
}

type balance struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

type balanceResult struct {
	TotalAmount uint64    `json:"total_amount"`
	Addresses   []balance `json:"addresses"`
}

func init() {
	cmd := gcli.Command{
		Name:      "checkBalance",
		ArgsUsage: "Check the balance of a wallet or specific address.",
		Usage:     "[option] [wallet path or address]",
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "w",
				Usage: "[wallet file or path], List balance of all addresses in a wallet.",
			},
			gcli.StringFlag{
				Name:  "a",
				Usage: "[address] List balance of specific address.",
			},
			// gcli.StringFlag{
			// 	Name:  "j,json",
			// 	Usage: "Returns the results in JSON format.",
			// },
		},
		Action: checkBalance,
	}
	Commands = append(Commands, cmd)
}

func checkBalance(c *gcli.Context) error {
	// get w option
	w := c.String("w")

	// get a option
	a := c.String("a")

	if w != "" && a != "" {
		// 1 1
		return errors.New("specify wallet or address, cannot set both")
	}

	addrs, err := gatherAddrs(w, a)
	if err != nil {
		return err
	}

	balRlt, err := getAddrsBalance(addrs)
	if err != nil {
		return err
	}

	var d []byte
	if a != "" {
		d, err = json.MarshalIndent(balRlt.Addresses[0], "", "    ")
	} else {
		d, err = json.MarshalIndent(balRlt, "", "    ")
	}

	if err != nil {
		return errJsonMarshal
	}
	fmt.Println(string(d))
	return nil
}

func gatherAddrs(w, a string) ([]string, error) {
	var addrs []string
	if a != "" {
		// 1 0
		addrs = append(addrs, a)
	} else {
		if w == "" {
			// 0 0
			w = filepath.Join(walletDir, defaultWalletName)
		} else {
			// 0 1
			if !strings.HasSuffix(w, walletExt) {
				return []string{}, fmt.Errorf("error wallet file name, must has %v extension", walletExt)
			}

			if filepath.Base(w) == w {
				w = filepath.Join(walletDir, w)
			} else {
				var err error
				w, err = filepath.Abs(w)
				if err != nil {
					return []string{}, err
				}
			}
		}

		wlt, err := wallet.Load(w)
		if err != nil {
			return []string{}, errLoadWallet
		}

		addresses := wlt.GetAddresses()
		for _, a := range addresses {
			addrs = append(addrs, a.String())
		}
	}
	return addrs, nil
}

func getAddrsBalance(addrs []string) (balanceResult, error) {
	balRlt := balanceResult{
		Addresses: make([]balance, len(addrs)),
	}

	for i, a := range addrs {
		balRlt.Addresses[i] = balance{
			Address: a,
		}
	}

	outs, err := getUnspent(addrs)
	if err != nil {
		return balanceResult{}, err
	}

	find := func(bals []balance, addr string) (int, error) {
		for i, b := range bals {
			if b.Address == addr {
				return i, nil
			}
		}
		return -1, errors.New("not exist")
	}

	for _, o := range outs {
		amt, err := strconv.ParseUint(o.Coins, 10, 64)
		if err != nil {
			return balanceResult{}, errors.New("error coins string")
		}

		i, err := find(balRlt.Addresses, o.Address)
		if err != nil {
			return balanceResult{}, fmt.Errorf("output belongs to no address")
		}
		balRlt.Addresses[i].Amount += amt
		balRlt.TotalAmount += amt
	}
	return balRlt, nil
}
