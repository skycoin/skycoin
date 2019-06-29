package cli

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gcli "github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip39"
	secp256k1 "github.com/skycoin/skycoin/src/cipher/secp256k1-go"
	"github.com/skycoin/skycoin/src/wallet"
)

const (
	// AlphaNumericSeedLength is the size of generated alphanumeric seeds, in bytes
	AlphaNumericSeedLength = 64
)

func walletCreateCmd() *gcli.Command {
	walletCreateCmd := &gcli.Command{
		Use:   "walletCreate",
		Short: "Generate a new wallet",
		Long: fmt.Sprintf(`The default wallet (%s) will be created if no wallet and address was specified.

    Use caution when using the "-p" command. If you have command
    history enabled your wallet encryption password can be recovered
    from the history log. If you do not include the "-p" option you will
    be prompted to enter your password after you enter your command.

    All results are returned in JSON format.`, cliConfig.FullWalletPath()),
		SilenceUsage: true,
		RunE:         generateWalletHandler,
	}

	walletCreateCmd.Flags().BoolP("random", "r", false, "A random alpha numeric seed will be generated")
	walletCreateCmd.Flags().BoolP("mnemonic", "m", false, "A mnemonic seed consisting of 12 dictionary words will be generated")
	walletCreateCmd.Flags().StringP("seed", "s", "", "Your seed")
	walletCreateCmd.Flags().Uint64P("num", "n", 1, `Number of addresses to generate. By default 1 address is generated.`)
	walletCreateCmd.Flags().StringP("wallet-file", "f", cliConfig.WalletName, `Name of wallet. The final format will be "yourName.wlt".
If no wallet name is specified a generic name will be selected.`)
	walletCreateCmd.Flags().StringP("label", "l", "", "Label used to idetify your wallet.")
	walletCreateCmd.Flags().BoolP("encrypt", "e", false, "Create encrypted wallet.")
	walletCreateCmd.Flags().StringP("crypto-type", "x", string(wallet.CryptoTypeScryptChacha20poly1305),
		"The crypto type for wallet encryption, can be scrypt-chacha20poly1305 or sha256-xor")
	walletCreateCmd.Flags().StringP("password", "p", "", "Wallet password")

	return walletCreateCmd
}

func generateWalletHandler(c *gcli.Command, _ []string) error {
	// create wallet dir if not exist
	if _, err := os.Stat(cliConfig.WalletDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cliConfig.WalletDir, 0750); err != nil {
			return errors.New("create dir failed")
		}
	}

	// get wallet name
	wltName := c.Flag("wallet-file").Value.String()

	// check if the wallet name has wlt extension.
	if !strings.HasSuffix(wltName, walletExt) {
		return ErrWalletName
	}

	// wallet file should not be a path.
	if filepath.Base(wltName) != wltName {
		return fmt.Errorf("wallet file name must not contain path")
	}

	// check if the wallet file does exist
	if _, err := os.Stat(filepath.Join(cliConfig.WalletDir, wltName)); err == nil {
		return fmt.Errorf("%v already exist", wltName)
	}

	// check if the wallet dir does exist.
	if _, err := os.Stat(cliConfig.WalletDir); os.IsNotExist(err) {
		return err
	}

	// get number of address that are need to be generated, if m is 0, set to 1.
	num, err := c.Flags().GetUint64("num")
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("-n must > 0")
	}

	// get label
	label := c.Flag("label").Value.String()

	// get seed
	s := c.Flag("seed").Value.String()
	random, err := c.Flags().GetBool("random")
	if err != nil {
		return err
	}

	mnemonic, err := c.Flags().GetBool("mnemonic")
	if err != nil {
		return err
	}

	encrypt, err := c.Flags().GetBool("encrypt")
	if err != nil {
		return err
	}

	sd, err := makeSeed(s, random, mnemonic)
	if err != nil {
		return err
	}

	cryptoType, err := wallet.CryptoTypeFromString(c.Flag("crypto-type").Value.String())
	if err != nil {
		return err
	}

	pr := NewPasswordReader([]byte(c.Flag("password").Value.String()))
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
		Password:   password,
	}

	wlt, err := GenerateWallet(wltName, opts, num)
	if err != nil {
		return err
	}

	if err := wallet.Save(wlt, cliConfig.WalletDir); err != nil {
		return err
	}

	return printJSON(wlt.ToReadable())
}

func makeSeed(s string, r, m bool) (string, error) {
	if s != "" {
		// 111, 101, 110
		if r || m {
			return "", errors.New("seed already specified, must not use -r or -m again")
		}
		// 100
		return s, nil
	}

	// 011
	if r && m {
		return "", errors.New("for -r and -m, only one option can be used")
	}

	// 010
	if r {
		return MakeAlphanumericSeed(), nil
	}

	// 001, 000
	return bip39.NewDefaultMnemonic()
}

// PUBLIC

// GenerateWallet generates a new wallet with filename walletFile, label, seed and number of addresses.
// Caller should save the wallet file to its chosen directory
func GenerateWallet(walletFile string, opts wallet.Options, numAddrs uint64) (wallet.Wallet, error) {
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

	if err := wallet.Lock(wlt, opts.Password, opts.CryptoType); err != nil {
		return nil, err
	}

	return wlt, nil
}

// MakeAlphanumericSeed creates a random seed with AlphaNumericSeedLength bytes and hex encodes it
func MakeAlphanumericSeed() string {
	seedRaw := cipher.SumSHA256(secp256k1.RandByte(AlphaNumericSeedLength))
	return hex.EncodeToString(seedRaw[:])
}
