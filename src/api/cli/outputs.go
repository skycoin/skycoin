package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

func walletOutputsCMD() gcli.Command {
	name := "walletOutputs"
	return gcli.Command{
		Name:      name,
		Usage:     "Display outputs of specific wallet",
		ArgsUsage: "[wallet file]",
		Description: fmt.Sprintf(`Display outputs of specific wallet, the default 
		wallet(%s/%s) will be 
		used if no wallet was specificed, use ENV 'WALLET_NAME' 
		to update default wallet file name, and 'WALLET_DIR' to update 
		the default wallet directory`, cfg.WalletDir, cfg.DefaultWalletName),
		OnUsageError: onCommandUsageError(name),
		Action:       wltOutputs,
	}
}

func addressOutputsCMD() gcli.Command {
	name := "addressOutputs"
	return gcli.Command{
		Name:      name,
		Usage:     "Display outputs of specific addresses",
		ArgsUsage: "[address list]",
		Description: `Display outputs of specific addresses, join multiple addresses with space,
        example: addressOutputs $addr1 $addr2 $addr3`,
		OnUsageError: onCommandUsageError(name),
		Action:       addrOutputs,
	}

}

func wltOutputs(c *gcli.Context) error {
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

	cipherAddrs := wlt.GetAddresses()
	addrs := make([]string, len(cipherAddrs))
	for i := range cipherAddrs {
		addrs[i] = cipherAddrs[i].String()
	}

	rlt, err := getAddrOutputs(addrs)
	if err != nil {
		return err
	}

	fmt.Println(rlt)
	return nil
}

func addrOutputs(c *gcli.Context) error {
	addrs := make([]string, c.NArg())
	var err error
	for i := 0; i < c.NArg(); i++ {
		addrs[i] = c.Args().Get(i)
		if _, err = cipher.DecodeBase58Address(addrs[i]); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", addrs[i], err)
		}
	}

	rlt, err := getAddrOutputs(addrs)
	if err != nil {
		return err
	}
	fmt.Println(rlt)
	return nil
}

func getAddrOutputs(addrs []string) (string, error) {
	req, err := webrpc.NewRequest("get_outputs", addrs, "1")
	if err != nil {
		return "", fmt.Errorf("do rpc request failed: %v", err)
	}

	rsp, err := webrpc.Do(req, cfg.RPCAddress)
	if err != nil {
		return "", err
	}

	if rsp.Error != nil {
		return "", fmt.Errorf("do rpc request failed: %+v", *rsp.Error)
	}

	return string(rsp.Result), nil
}
