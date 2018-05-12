package cli

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/go-bip39"
	secp256k1 "github.com/skycoin/skycoin/src/cipher/secp256k1-go"
	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/urfave/cli"
)

const (
	// AlphaNumericSeedLength is the size of generated alphanumeric seeds, in bytes
	AlphaNumericSeedLength = 64
)

func generateWalletCmd(cfg Config) gcli.Command {
	name := "generateWallet"
	return gcli.Command{
		Name:         "generateWallet",
		Usage:        "Generate a new wallet",
		ArgsUsage:    " ",
		OnUsageError: onCommandUsageError(name),
		Description: fmt.Sprintf(`The default wallet (%s) will
		be created if no wallet and address was specified.

		Use caution when using the "-p" command. If you have command
		history enabled your wallet encryption password can be recovered
		from the history log. If you do not include the "-p" option you will
		be prompted to enter your password after you enter your command.

		All results are returned in JSON format.`, cfg.FullWalletPath()),
		Flags: []gcli.Flag{
			gcli.BoolFlag{
				Name:  "r",
				Usage: "A random alpha numeric seed will be generated for you",
			},
			gcli.BoolFlag{
				Name:  "rd",
				Usage: "A random seed consisting of 12 dictionary words will be generated for you",
			},
			gcli.StringFlag{
				Name:  "s",
				Usage: "Your seed",
			},
			gcli.UintFlag{
				Name:  "n",
				Value: 1,
				Usage: `[numberOfAddresses] Number of addresses to generate
						By default 1 address is generated.`,
			},
			gcli.StringFlag{
				Name:  "f",
				Value: cfg.WalletName,
				Usage: `[walletName] Name of wallet. The final format will be "yourName.wlt".
						 If no wallet name is specified a generic name will be selected.`,
			},
			gcli.StringFlag{
				Name:  "l",
				Usage: "[label] Label used to idetify your wallet.",
			},
			gcli.BoolFlag{
				Name:  "e,encrypt",
				Usage: `Whether creates wallet with encryption `,
			},
			gcli.StringFlag{
				Name:  "x,crypto-type",
				Value: string(wallet.CryptoTypeScryptChacha20poly1305),
				Usage: "[crypto type] The crypto type for wallet encryption, can be scrypt-chacha20poly1305 or sha256-xor",
			},
			gcli.StringFlag{
				Name:  "p",
				Usage: "[password] Wallet password",
			},
		},
		Action: generateWalletHandler,
	}
}

func generateWalletHandler(c *gcli.Context) error {
	cfg := ConfigFromContext(c)

	// create wallet dir if not exist
	if _, err := os.Stat(cfg.WalletDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.WalletDir, 0755); err != nil {
			return errors.New("create dir failed")
		}
	}

	// get wallet name
	wltName := c.String("f")

	// check if the wallet name has wlt extension.
	if !strings.HasSuffix(wltName, ".wlt") {
		return ErrWalletName
	}

	// wallet file should not be a path.
	if filepath.Base(wltName) != wltName {
		return fmt.Errorf("wallet file name must not contain path")
	}

	// check if the wallet file does exist
	if _, err := os.Stat(filepath.Join(cfg.WalletDir, wltName)); err == nil {
		return fmt.Errorf("%v already exist", wltName)
	}

	// check if the wallet dir does exist.
	if _, err := os.Stat(cfg.WalletDir); os.IsNotExist(err) {
		return err
	}

	// get number of address that are need to be generated, if m is 0, set to 1.
	num := c.Uint64("n")
	if num == 0 {
		return errors.New("-n must > 0")
	}

	// get label
	label := c.String("l")

	// get seed
	s := c.String("s")
	r := c.Bool("r")
	rd := c.Bool("rd")

	encrypt := c.Bool("e")

	sd, err := makeSeed(s, r, rd)
	if err != nil {
		return err
	}

	cryptoType, err := wallet.CryptoTypeFromString(c.String("x"))
	if err != nil {
		return err
	}

	pr := NewPasswordReader([]byte(c.String("p")))
	switch pr.(type) {
	case PasswordFromBytes:
		p, err := pr.Password()
		if err != nil {
			return err
		}

		if !encrypt && len(p) != 0 {
			return errors.New("password should not be set as we're not going to create a wallet with encryption")
		}
	}

	var password []byte
	if encrypt {
		var err error
		password, err = pr.Password()
		if err != nil {
			return err
		}
	}

	opts := wallet.Options{
		Label:      label,
		Seed:       sd,
		Encrypt:    encrypt,
		CryptoType: cryptoType,
		Password:   []byte(password),
	}

	wlt, err := GenerateWallet(wltName, opts, num)
	if err != nil {
		return err
	}

	if err := wlt.Save(cfg.WalletDir); err != nil {
		return err
	}

	return printJSON(wallet.NewReadableWallet(wlt))
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
		return MakeAlphanumericSeed(), nil
	}

	// 001, 000
	return bip39.NewDefaultMnemomic()
}

// PUBLIC

// GenerateWallet generates a new wallet with filename walletFile, label, seed and number of addresses.
// Caller should save the wallet file to its chosen directory
func GenerateWallet(walletFile string, opts wallet.Options, numAddrs uint64) (*wallet.Wallet, error) {
	walletFile = filepath.Base(walletFile)

	wlt, err := wallet.NewWallet(walletFile, wallet.Options{
		Seed:  opts.Seed,
		Label: opts.Label,
	})
	if err != nil {
		return nil, err
	}

	if numAddrs > 1 {
		if _, err := wlt.GenerateAddresses(numAddrs - 1); err != nil {
			return nil, err
		}
	}

	if !opts.Encrypt {
		if len(opts.Password) != 0 {
			return nil, wallet.ErrWalletNotEncrypted
		}

		return wlt, nil
	}

	if err := wlt.Lock(opts.Password, opts.CryptoType); err != nil {
		return nil, err
	}

	return wlt, nil
}

// MakeAlphanumericSeed creates a random seed with AlphaNumericSeedLength bytes and hex encodes it
func MakeAlphanumericSeed() string {
	seedRaw := cipher.SumSHA256(secp256k1.RandByte(AlphaNumericSeedLength))
	return hex.EncodeToString(seedRaw[:])
}
