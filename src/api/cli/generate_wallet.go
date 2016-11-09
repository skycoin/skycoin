package cli

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	bip39 "github.com/skycoin/skycoin/src/api/cli/go-bip39"
	"github.com/skycoin/skycoin/src/cipher"
	secp256k1 "github.com/skycoin/skycoin/src/cipher/secp256k1-go"
	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:      "generateWallet",
		Usage:     "Generate a new wallet",
		ArgsUsage: " ",
		Description: `Use caution when using the "-p" command. If you have command 
		history enabled your wallet encryption password can be recovered 
		from the history log. If you do not include the "-p" option you will 
		be prompted to enter your password after you enter your command. 
		
		All results are returned in JSON format.`,
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "s",
				Usage: "Your seed",
			},
			gcli.BoolFlag{
				Name:  "r",
				Usage: "A random alpha numeric seed will be generated for you",
			},
			gcli.BoolFlag{
				Name:  "rd",
				Usage: "A random seed consisting of 12 dictionary words will be generated for you",
			},
			gcli.IntFlag{
				Name:  "n",
				Value: 1,
				Usage: `[numberOfAddresses] Number of addresses to generate 
						By default 1 address is generated.`,
			},
			gcli.StringFlag{
				Name:  "f",
				Value: defaultWalletName,
				Usage: `[walletName] Name of wallet. The final format will be "yourName.wlt". 
						 If no wallet name is specified a generic name will be selected.`,
			},
			gcli.StringFlag{
				Name:  "l",
				Usage: "[label] Label used to idetify your wallet.",
			},
		},
		Action: generateWallet,
	}
	Commands = append(Commands, cmd)
}

func generateWallet(c *gcli.Context) error {
	// create wallet dir if not exist
	if _, err := os.Stat(walletDir); os.IsNotExist(err) {
		if err := os.MkdirAll(walletDir, 0755); err != nil {
			return errors.New("create dir failed")
		}
	}

	// get wallet name
	wltName := c.String("f")

	// check if the wallet name has wlt extension.
	if !strings.HasSuffix(wltName, ".wlt") {
		return errWalletName
	}

	// wallet file should not be a path.
	if filepath.Base(wltName) != wltName {
		return fmt.Errorf("wallet file name must not contain path")
	}

	// check if the wallet file does exist
	if _, err := os.Stat(filepath.Join(walletDir, wltName)); err == nil {
		return fmt.Errorf("%v already exist", wltName)
	}

	// get number of address that are need to be generated, if m is 0, set to 1.
	num := c.Int("n")
	if num == 0 {
		num = 1
	}

	// get label
	label := c.String("l")

	// get seed
	s := c.String("s")
	r := c.Bool("r")
	rd := c.Bool("rd")

	sd, err := makeSeed(s, r, rd)
	if err != nil {
		return err
	}
	wlt := wallet.NewWallet(sd, wltName, label)
	wlt.GenerateAddresses(num)

	// check if the wallet dir does exist.
	if _, err := os.Stat(walletDir); os.IsNotExist(err) {
		return err
	}

	if err := wlt.Save(walletDir); err != nil {
		return err
	}

	rwlt := wallet.NewReadableWallet(wlt)
	d, err := json.MarshalIndent(rwlt, "", "    ")
	if err != nil {
		return errJSONMarshal
	}
	fmt.Println(string(d))
	return nil
}

func makeSeed(s string, r, rd bool) (string, error) {
	if s != "" {
		// 111, 101, 110
		if r || rd {
			return "", errors.New("seed already specified, must not use -r or -rd again")
		}
		// 100
		return s, nil
	}

	// 011
	if r && rd {
		return "", errors.New("for -r and -rd, only one option can be used")
	}

	// 010
	if r {
		seedRaw := cipher.SumSHA256(secp256k1.RandByte(64))
		return hex.EncodeToString(seedRaw[:]), nil
	}

	// 001, 000
	// https://github.com/tyler-smith/go-bip39 generate mnemonic.
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic, nil
}
