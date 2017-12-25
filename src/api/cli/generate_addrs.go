package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	gcli "github.com/urfave/cli"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

func generateAddrsCmd(cfg Config) gcli.Command {
	name := "generateAddresses"
	return gcli.Command{
		Name:      name,
		Usage:     "Generate additional addresses for a wallet",
		ArgsUsage: " ",
		Description: fmt.Sprintf(`The default wallet (%s) will
		be used if no wallet and address was specified.

		Use caution when using the "-p" command. If you have command
		history enabled your wallet encryption password can be recovered from the
		history log. If you do not include the "-p" option you will be prompted to
		enter your password after you enter your command.`, cfg.fullWalletPath()),
		Flags: []gcli.Flag{
			gcli.UintFlag{
				Name:  "n",
				Value: 1,
				Usage: `[numberOfAddresses]	Number of addresses to generate`,
			},
			gcli.StringFlag{
				Name:  "f",
				Value: cfg.fullWalletPath(),
				Usage: `[wallet file or path] Generate addresses in the wallet`,
			},
			gcli.BoolFlag{
				Name:  "json,j",
				Usage: "Returns the results in JSON format",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action:       generateAddrs,
	}
	// Commands = append(Commands, cmd)
}

func generateAddrs(c *gcli.Context) error {
	cfg := configFromContext(c)

	// get number of address that are need to be generated.
	num := c.Uint64("n")
	if num == 0 {
		return errors.New("-n must > 0")
	}

	jsonFmt := c.Bool("json")

	w, err := resolveWalletPath(cfg, c.String("f"))
	if err != nil {
		return err
	}

	addrs, err := generateAddressesInFile(w, num)

	switch err.(type) {
	case nil:
	case walletLoadError:
		errorWithHelp(c, err)
		return nil
	case walletSaveError:
		return errors.New("save wallet failed")
	default:
		return err
	}

	if jsonFmt {
		s, err := formatAddressesAsJSON(addrs)
		if err != nil {
			return err
		}
		fmt.Println(s)
	} else {
		fmt.Println(formatAddressesAsJoinedArray(addrs))
	}

	return nil
}

func generateAddressesInFile(walletFile string, num uint64) ([]cipher.Address, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, walletLoadError(err)
	}

	addrs := wlt.GenerateAddresses(num)

	dir, err := filepath.Abs(filepath.Dir(walletFile))
	if err != nil {
		return nil, err
	}

	if err := wlt.Save(dir); err != nil {
		return nil, walletSaveError(err)
	}

	return addrs, nil
}

func formatAddressesAsJSON(addrs []cipher.Address) (string, error) {
	d, err := formatJSON(struct {
		Addresses []string `json:"addresses"`
	}{
		Addresses: addressesToStrings(addrs),
	})

	if err != nil {
		return "", err
	}

	return string(d), nil
}

func formatAddressesAsJoinedArray(addrs []cipher.Address) string {
	return strings.Join(addressesToStrings(addrs), ",")
}

func addressesToStrings(addrs []cipher.Address) []string {
	if addrs == nil {
		return nil
	}

	addrsStr := make([]string, len(addrs))
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}

	return addrsStr
}
