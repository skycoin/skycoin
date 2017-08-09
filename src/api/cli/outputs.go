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
		Action:       getWalletOutputsCmd,
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
		Action:       getAddressOutputsCmd,
	}

}

func getWalletOutputsCmd(c *gcli.Context) error {
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

	outputs, err := GetWalletOutputsFromFile(w)
	if err != nil {
		return err
	}

	return printJson(outputs)
}

func getAddressOutputsCmd(c *gcli.Context) error {
	addrs := make([]string, c.NArg())
	var err error
	for i := 0; i < c.NArg(); i++ {
		addrs[i] = c.Args().Get(i)
		if _, err = cipher.DecodeBase58Address(addrs[i]); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", addrs[i], err)
		}
	}

	outputs, err := GetAddressOutputs(addrs)
	if err != nil {
		return err
	}

	return printJson(outputs)
}

// PUBLIC

func GetWalletOutputsFromFile(walletFile string) (*webrpc.OutputsResult, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, err
	}

	return GetWalletOutputs(wlt)
}

func GetWalletOutputs(wlt *wallet.Wallet) (*webrpc.OutputsResult, error) {
	cipherAddrs := wlt.GetAddresses()
	addrs := make([]string, len(cipherAddrs))
	for i := range cipherAddrs {
		addrs[i] = cipherAddrs[i].String()
	}

	return GetAddressOutputs(addrs)
}

func GetAddressOutputs(addrs []string) (*webrpc.OutputsResult, error) {
	outputs := webrpc.OutputsResult{}
	if err := DoRpcRequest(&outputs, "get_outputs", addrs, "1"); err != nil {
		return nil, err
	}

	return &outputs, nil
}
