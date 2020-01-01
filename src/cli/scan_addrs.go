package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func walletScanAddressesCmd() *cobra.Command {
	walletScanAddressesCmd := &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "walletScanAddresses [wallet]",
		Short: "Scan addresses ahead for deterministic, bip44 or xpub wallet",
		Long: `Scan addresses ahead for deterministic, bip44 or xpub wallet.

    Warning: if you generate long (over 20) sequences of empty addresses and use
    a later address this can cause the wallet history scanner to miss your addresses,
    if you load the wallet from seed elsewhere. In that case, you'll have to manually
    generate addresses to cover the gap of unused addresses in the sequence.

    BIP44 wallets generate their addresses on the external (0'/0) chain.

    Use caution when using the "-p" command. If you have command
    history enabled your wallet encryption password can be recovered from the
    history log. If you do not include the "-p" option you will be prompted to
    enter your password after you enter your command.`,
		RunE:         runScanAddresses,
		SilenceUsage: true,
	}

	walletScanAddressesCmd.Flags().Uint64P("num", "n", 20, "Number of addresses to scan ahead")
	walletScanAddressesCmd.Flags().StringP("password", "p", "", "wallet password")
	walletScanAddressesCmd.Flags().BoolP("json", "j", false, "Returns the results in json format")

	return walletScanAddressesCmd
}

func runScanAddresses(c *cobra.Command, args []string) error {
	// get the number of addresses to scan ahead
	num, err := c.Flags().GetUint64("num")
	if err != nil {
		return err
	}

	if num == 0 {
		return errors.New("--num or -n must be > 0")
	}

	jsonFmt, err := c.Flags().GetBool("json")
	if err != nil {
		return err
	}

	wltFile := args[0]
	dir, id := filepath.Split(wltFile)
	if dir != "" {
		if _, err := os.Stat(wltFile); os.IsNotExist(err) {
			return fmt.Errorf("wallet file %s does not exist", wltFile)
		}
	}

	rsp, err := apiClient.Wallet(id)
	if err != nil {
		return err
	}

	pr := NewPasswordReader([]byte(c.Flag("password").Value.String()))
	var password []byte
	if rsp.Meta.Encrypted {
		var err error
		password, err = pr.Password()
		if err != nil {
			return err
		}
		defer func() {
			password = []byte("")
		}()
	}

	addrs, err := apiClient.ScanWalletAddresses(id, int(num), string(password))
	if err != nil {
		return err
	}

	if jsonFmt {
		var obj = struct {
			Addresses []string `json:"addresses"`
		}{
			Addresses: addrs,
		}

		return printJSON(obj)
	}

	for _, addr := range addrs {
		fmt.Println(addr)
	}
	return nil
}
