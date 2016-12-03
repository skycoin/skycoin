package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
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
		Usage:     "Check the balance of a wallet or specific address",
		ArgsUsage: "[wallet or address]",
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "f",
				Usage: "[wallet file or path] List balance of all addresses in a wallet",
			},
		},
		Action: checkBalance,
	}
	Commands = append(Commands, cmd)
}

func checkBalance(c *gcli.Context) error {
	addrs, err := gatherAddrs(c)
	if err != nil {
		return err
	}

	balRlt, err := getAddrsBalance(addrs)
	if err != nil {
		return err
	}

	var d []byte
	d, err = json.MarshalIndent(balRlt, "", "    ")
	if err != nil {
		return errJSONMarshal
	}
	fmt.Println(string(d))
	return nil
}

func gatherAddrs(c *gcli.Context) ([]string, error) {
	w := c.String("f")
	var a string
	if c.NArg() > 0 {
		a = c.Args().First()
		if _, err := cipher.DecodeBase58Address(a); err != nil {
			return []string{}, err
		}
	}

	addrs := []string{}
	if w == "" && a == "" {
		// use default wallet
		w = filepath.Join(cfg.WalletDir, cfg.DefaultWalletName)
	}

	if w != "" {
		if !strings.HasSuffix(w, walletExt) {
			return []string{}, fmt.Errorf("error wallet file name, must has %v extension", walletExt)
		}

		if filepath.Base(w) == w {
			w = filepath.Join(cfg.WalletDir, w)
		} else {
			var err error
			w, err = filepath.Abs(w)
			if err != nil {
				return []string{}, err
			}
		}

		wlt, err := wallet.Load(w)
		if err != nil {
			return []string{}, err
		}

		addresses := wlt.GetAddresses()
		for _, a := range addresses {
			addrs = append(addrs, a.String())
		}
	}

	if a != "" {
		addrs = append(addrs, a)
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
