package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

type unspentOut struct {
	visor.ReadableOutput
}

type unspentOutSet struct {
	visor.ReadableOutputSet
}

type balance struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

type balanceResult struct {
	TotalAmount uint64    `json:"total_amount"`
	Addresses   []balance `json:"addresses"`
}

func walletBalanceCMD() gcli.Command {
	name := "walletBalance"
	return gcli.Command{
		Name:      name,
		Usage:     "Check the balance of a wallet",
		ArgsUsage: "[wallet]",
		Description: fmt.Sprintf(`Check balance of specific wallet, the default 
		wallet(%s/%s) will be 
		used if no wallet was specificed, use ENV 'WALLET_NAME' 
		to update default wallet file name, and 'WALLET_DIR' to update 
		the default wallet directory`, cfg.WalletDir, cfg.DefaultWalletName),
		OnUsageError: onCommandUsageError(name),
		Action:       checkWltBalance,
	}
}

func addressBalanceCMD() gcli.Command {
	name := "addressBalance"
	return gcli.Command{
		Name:      name,
		Usage:     "Check the balance of specific addresses",
		ArgsUsage: "[addresses]",
		Description: `Check balance of specific addresses, join multiple addresses with space.
		example: addressBalance "$addr1 $addr2 $addr3"`,
		OnUsageError: onCommandUsageError(name),
		Action:       addrBalance,
	}
}

func checkWltBalance(c *gcli.Context) error {
	var w string
	if c.NArg() == 0 {
		w = filepath.Join(cfg.WalletDir, cfg.DefaultWalletName)
	} else {
		w = c.Args().First()
		if !strings.HasSuffix(w, walletExt) {
			return errWalletName
		}

		var err error
		if filepath.Base(w) == w {
			w = filepath.Join(cfg.WalletDir, w)
		} else {
			w, err = filepath.Abs(w)
			if err != nil {
				return err
			}
		}
	}

	wlt, err := wallet.Load(w)
	if err != nil {
		return err
	}

	var addrs []string
	addresses := wlt.GetAddresses()
	for _, a := range addresses {
		// validate the address
		addrs = append(addrs, a.String())
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

func addrBalance(c *gcli.Context) error {
	addrs := make([]string, c.NArg())
	var err error
	for i := 0; i < c.NArg(); i++ {
		addrs[i] = c.Args().Get(i)
		if _, err = cipher.DecodeBase58Address(addrs[i]); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", addrs[i], err)
		}
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

	for _, o := range outs.HeadOutputs {
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
