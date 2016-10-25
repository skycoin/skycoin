package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

var defaultAddrNum = 1

func init() {
	cmd := gcli.Command{
		Name:      "generateAddresses",
		ArgsUsage: "Generate additional addresses for a wallet.",
		Usage:     "[options]",
		Description: `
        Use caution when using the “-p” command. If you have command 
        history enabled your wallet encryption password can be recovered from the 
		history log. If you do not include the “-p” option you will be prompted to 
		enter your password after you enter your command.`,
		Flags: []gcli.Flag{
			gcli.IntFlag{
				Name:  "m",
				Value: 1,
				Usage: `[numberOfAddresses]	Number of addresses to generate. 
						By default 1 address is generated.`,
			},
			gcli.StringFlag{
				Name: "w",
				Usage: `[wallet file or path] In wallet. If no path is specified your default 
					wallet path will be used.`,
			},
			// gcli.StringFlag{
			// 	Name:  "p",
			// 	Usage: "[password] Password used to encrypt the wallet locally.",
			// },
			gcli.BoolFlag{
				Name:  "j,json",
				Usage: "Returns the results in JSON format.",
			},
		},
		Action: generateAddrs,
	}
	Commands = append(Commands, cmd)
}

func generateAddrs(c *gcli.Context) error {
	// get number of address that are need to be generated.
	num := c.Int("m")
	if num == 0 {
		num = defaultAddrNum
	}

	jsonFmt := c.Bool("json")

	w := c.String("w")
	if w == "" {
		w = filepath.Join(walletDir, defaultWalletName)
	}

	if !strings.HasSuffix(w, walletExt) {
		return errWalletName
	}

	// only wallet file name, no path.
	if filepath.Base(w) == w {
		w = filepath.Join(walletDir, w)
	}

	wlt, err := wallet.Load(w)
	if err != nil {
		return errLoadWallet
	}

	addrs := wlt.GenerateAddresses(num)
	dir, err := filepath.Abs(filepath.Dir(w))
	if err != nil {
		return err
	}

	if err := wlt.Save(dir); err != nil {
		return errors.New("save wallet failed")
	}

	s, err := addrResult(addrs, jsonFmt)
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}

func addrResult(addrs []cipher.Address, jsonFmt bool) (string, error) {
	if jsonFmt {
		var rlt = struct {
			Addresses []string `json:"addresses"`
		}{
			make([]string, len(addrs)),
		}

		for i, a := range addrs {
			rlt.Addresses[i] = a.String()
		}
		d, err := json.MarshalIndent(rlt, "", "    ")
		if err != nil {
			return "", errors.New("json marshal failed")
		}
		return string(d), nil
	}

	addrArray := make([]string, len(addrs))
	for i, a := range addrs {
		addrArray[i] = a.String()
	}
	return strings.Join(addrArray, ","), nil
}
