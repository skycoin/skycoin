package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	gcli "github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

func walletAddAddressesCmd() *gcli.Command {
	walletAddAddressesCmd := &gcli.Command{
		Use:   "walletAddAddresses",
		Short: "Generate additional addresses for a wallet",
		Long: fmt.Sprintf(`The default wallet (%s) will be used if no wallet was specified.

    Use caution when using the "-p" command. If you have command
    history enabled your wallet encryption password can be recovered from the
    history log. If you do not include the "-p" option you will be prompted to
    enter your password after you enter your command.`, cliConfig.FullWalletPath()),
		SilenceUsage: true,
		RunE:         generateAddrs,
	}

	walletAddAddressesCmd.Flags().Uint64P("num", "n", 1, "Number of addresses to generate")
	walletAddAddressesCmd.Flags().StringP("wallet-file", "f", cliConfig.FullWalletPath(), "Generate addresses in the wallet")
	walletAddAddressesCmd.Flags().StringP("password", "p", "", "wallet password")
	walletAddAddressesCmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format")

	return walletAddAddressesCmd
}

func generateAddrs(c *gcli.Command, _ []string) error {
	// get number of address that are need to be generated.
	num, err := c.Flags().GetUint64("num")
	if err != nil {
		return err
	}

	if num == 0 {
		return errors.New("-n must > 0")
	}

	jsonFmt, err := c.Flags().GetBool("json")
	if err != nil {
		return err
	}

	w, err := resolveWalletPath(cliConfig, c.Flag("wallet-file").Value.String())
	if err != nil {
		return err
	}

	pr := NewPasswordReader([]byte(c.Flag("password").Value.String()))
	addrs, err := GenerateAddressesInFile(w, num, pr)

	switch err.(type) {
	case nil:
	case WalletLoadError:
		printHelp(c)
		return err
	default:
		return err
	}

	if jsonFmt {
		s, err := FormatAddressesAsJSON(addrs)
		if err != nil {
			return err
		}
		fmt.Println(s)
	} else {
		fmt.Println(FormatAddressesAsJoinedArray(addrs))
	}

	return nil
}

// GenerateAddressesInFile generates addresses in given wallet file
func GenerateAddressesInFile(walletFile string, num uint64, pr PasswordReader) ([]cipher.Addresser, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, WalletLoadError{err}
	}

	switch pr.(type) {
	case nil:
		if wlt.IsEncrypted() {
			return nil, wallet.ErrWalletEncrypted
		}
	case PasswordFromBytes:
		p, err := pr.Password()
		if err != nil {
			return nil, err
		}

		if !wlt.IsEncrypted() && len(p) != 0 {
			return nil, wallet.ErrWalletNotEncrypted
		}
	}

	genAddrsInWallet := func(w *wallet.Wallet, n uint64) ([]cipher.Addresser, error) {
		return w.GenerateAddresses(n)
	}

	if wlt.IsEncrypted() {
		genAddrsInWallet = func(w *wallet.Wallet, n uint64) ([]cipher.Addresser, error) {
			password, err := pr.Password()
			if err != nil {
				return nil, err
			}

			var addrs []cipher.Addresser
			if err := w.GuardUpdate(password, func(wlt *wallet.Wallet) error {
				var err error
				addrs, err = wlt.GenerateAddresses(n)
				return err
			}); err != nil {
				return nil, err
			}

			return addrs, nil
		}
	}

	addrs, err := genAddrsInWallet(wlt, num)
	if err != nil {
		return nil, err
	}

	dir, err := filepath.Abs(filepath.Dir(walletFile))
	if err != nil {
		return nil, err
	}

	if err := wlt.Save(dir); err != nil {
		return nil, WalletSaveError{err}
	}

	return addrs, nil
}

// FormatAddressesAsJSON converts []cipher.Address to strings and formats the array into a standard JSON object wrapper
func FormatAddressesAsJSON(addrs []cipher.Addresser) (string, error) {
	d, err := formatJSON(struct {
		Addresses []string `json:"addresses"`
	}{
		Addresses: AddressesToStrings(addrs),
	})

	if err != nil {
		return "", err
	}

	return string(d), nil
}

// FormatAddressesAsJoinedArray converts []cipher.Address to strings and concatenates them with a comma
func FormatAddressesAsJoinedArray(addrs []cipher.Addresser) string {
	return strings.Join(AddressesToStrings(addrs), ",")
}

// AddressesToStrings converts []cipher.Address to []string
func AddressesToStrings(addrs []cipher.Addresser) []string {
	if addrs == nil {
		return nil
	}

	addrsStr := make([]string, len(addrs))
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}

	return addrsStr
}
