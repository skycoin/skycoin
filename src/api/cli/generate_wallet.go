package cli

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	secp256k1 "github.com/skycoin/skycoin/src/cipher/secp256k1-go"
	"github.com/skycoin/skycoin/src/wallet"

	gcli "gopkg.in/urfave/cli.v1"
)

func init() {
	cmd := gcli.Command{
		Name:      "generateWallet",
		Usage:     "Generate a new wallet from seed.",
		ArgsUsage: "[options]",
		Description: `
		Use caution when using the “-p” command. If you have command history enabled your 
		wallet encryption password can be recovered from the history log. If you do not 
		include the “-p” option you will be prompted to enter your password after you enter 
		your command. 
		
		All results are returned in JSON format. 
                      `,
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "s",
				Usage: "Your seed.",
			},
			gcli.StringFlag{
				Name:  "r",
				Usage: "A random alpha numeric seed will be generated for you.",
			},
			gcli.StringFlag{
				Name:  "rd",
				Usage: "A random seed consisting of 12 dictionary words will be generated for you.",
			},
			gcli.IntFlag{
				Name:  "m",
				Usage: "[numberOfAddresses] Number of addresses to generate. By default 1 address is generated.",
			},
			// gcli.StringFlag{
			// 	Name:  "p",
			// 	Usage: "Password used to encrypt the wallet locally.",
			// },
			gcli.StringFlag{
				Name:  "n",
				Usage: `[walletName] Name of wallet. The final format will be "yourName.wlt". If no wallet name is specified a generic name will be selected.`,
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
			return err
		}
	}

	// get wallet name
	wltName := c.String("n")
	if wltName == "" {
		wltName = defaultWalletName
	} else if wltName == defaultWalletName {
		return fmt.Errorf("wallet of %s name already exist, please choose another one", defaultWalletName)
	}

	// get number of address need to be generated.
	m := c.String("m")
	if m == "" || m == "0" {
		m = "1"
	}

	addrNum, err := strconv.Atoi(m)
	if err != nil {
		return fmt.Errorf("error address number:%v", err)
	}

	// get label
	// label := c.String("l")

	// get password
	// pwd := c.String("p")
	// if pwd == "" {
	// 	// TODO: show message of password request
	// }

	// get seed
	s := c.String("s")
	r := c.Bool("r")
	rd := c.Bool("rd")

	sd, err := makeSeed(s, r, rd)
	if err != nil {
		return err
	}
	wlt := wallet.NewWallet(sd, wltName)
	wlt.GenerateAddresses(addrNum)

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
		return err
	}
	fmt.Println(string(d))

	return nil
}

func makeSeed(s string, r, rd bool) (string, error) {
	if s != "" {
		if r || rd {
			return "", errors.New("seed already specified, must not use -r or -rd again")
		}
		return s, nil
	}

	if r && rd {
		return "", errors.New("for -r and -rd, only one option can be used")
	}

	if r {
		seedRaw := cipher.SumSHA256(secp256k1.RandByte(64))
		return hex.EncodeToString(seedRaw[:]), nil
	}

	if rd {
		return "", errors.New("not support yet")
	}
	return "", errors.New("no seed option found")
}
