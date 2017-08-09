package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

type Balance struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

type BalanceResult struct {
	TotalAmount uint64    `json:"total_amount"`
	Addresses   []Balance `json:"addresses"`
}

func walletBalanceCmd(cfg Config) gcli.Command {
	name := "walletBalance"
	return gcli.Command{
		Name:      name,
		Usage:     "Check the balance of a wallet",
		ArgsUsage: "[wallet]",
		Description: fmt.Sprintf(`Check balance of specific wallet, the default
		wallet(%s/%s) will be
		used if no wallet was specificed, use ENV 'WALLET_NAME'
		to update default wallet file name, and 'WALLET_DIR' to update
		the default wallet directory`, cfg.WalletDir, cfg.WalletName),
		OnUsageError: onCommandUsageError(name),
		Action:       checkWltBalance,
	}
}

func addressBalanceCmd() gcli.Command {
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
	cfg := c.App.Metadata["config"].(Config)
	rpcClient := c.App.Metadata["rpc"].(*RpcClient)

	var w string
	if c.NArg() > 0 {
		w = c.Args().First()
	}

	var err error
	w, err = resolveWalletPath(cfg, w)
	if err != nil {
		return err
	}

	balRlt, err := rpcClient.CheckWalletBalance(w)
	if err != nil {
		return err
	}

	return printJson(balRlt)
}

func addrBalance(c *gcli.Context) error {
	rpcClient := c.App.Metadata["rpc"].(*RpcClient)

	addrs := make([]string, c.NArg())
	var err error
	for i := 0; i < c.NArg(); i++ {
		addrs[i] = c.Args().Get(i)
		if _, err = cipher.DecodeBase58Address(addrs[i]); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", addrs[i], err)
		}
	}

	balRlt, err := rpcClient.GetBalanceOfAddresses(addrs)
	if err != nil {
		return err
	}

	return printJson(balRlt)
}

// PUBLIC

func (c *RpcClient) CheckWalletBalance(walletFile string) (BalanceResult, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return BalanceResult{}, err
	}

	var addrs []string
	addresses := wlt.GetAddresses()
	for _, a := range addresses {
		// validate the address
		addrs = append(addrs, a.String())
	}

	return c.GetBalanceOfAddresses(addrs)
}

func (c *RpcClient) GetBalanceOfAddresses(addrs []string) (BalanceResult, error) {
	balRlt := BalanceResult{
		Addresses: make([]Balance, len(addrs)),
	}

	for i, a := range addrs {
		balRlt.Addresses[i] = Balance{
			Address: a,
		}
	}

	outs, err := c.GetUnspent(addrs)
	if err != nil {
		return BalanceResult{}, err
	}

	find := func(bals []Balance, addr string) (int, error) {
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
			return BalanceResult{}, errors.New("error coins string")
		}

		i, err := find(balRlt.Addresses, o.Address)
		if err != nil {
			return BalanceResult{}, fmt.Errorf("output belongs to no address")
		}
		balRlt.Addresses[i].Amount += amt
		balRlt.TotalAmount += amt
	}
	return balRlt, nil
}
