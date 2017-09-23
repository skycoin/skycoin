package cli

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/droplet"
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
		wallet (%s) will be
		used if no wallet was specified, use ENV 'WALLET_NAME'
		to update default wallet file name, and 'WALLET_DIR' to update
		the default wallet directory`, cfg.FullWalletPath()),
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
	cfg := ConfigFromContext(c)
	rpcClient := RpcClientFromContext(c)

	var w string
	if c.NArg() > 0 {
		w = c.Args().First()
	}

	var err error
	w, err = resolveWalletPath(cfg, w)
	if err != nil {
		return err
	}

	balRlt, err := CheckWalletBalance(rpcClient, w)
	if err != nil {
		return err
	}

	return printJson(balRlt)
}

func addrBalance(c *gcli.Context) error {
	rpcClient := RpcClientFromContext(c)

	addrs := make([]string, c.NArg())
	var err error
	for i := 0; i < c.NArg(); i++ {
		addrs[i] = c.Args().Get(i)
		if _, err = cipher.DecodeBase58Address(addrs[i]); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", addrs[i], err)
		}
	}

	balRlt, err := GetBalanceOfAddresses(rpcClient, addrs)
	if err != nil {
		return err
	}

	return printJson(balRlt)
}

// PUBLIC

func CheckWalletBalance(c *webrpc.Client, walletFile string) (BalanceResult, error) {
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

	return GetBalanceOfAddresses(c, addrs)
}

func GetBalanceOfAddresses(c *webrpc.Client, addrs []string) (BalanceResult, error) {
	balRlt := BalanceResult{
		Addresses: make([]Balance, len(addrs)),
	}

	for i, a := range addrs {
		balRlt.Addresses[i] = Balance{
			Address: a,
		}
	}

	outs, err := c.GetUnspentOutputs(addrs)
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

	for _, o := range outs.Outputs.HeadOutputs {
		amt, err := droplet.FromString(o.Coins)
		if err != nil {
			return BalanceResult{}, fmt.Errorf("error coins string: %v", err)
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
