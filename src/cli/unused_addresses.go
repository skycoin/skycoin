package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/wallet"
)

// UnusedAddressResult represents the result of finding unused addresses
type UnusedAddressResult struct {
	Wallet          string   `json:"wallet"`
	TotalAddresses  int      `json:"total_addresses"`
	UnusedAddresses []string `json:"unused_addresses"`
	UnusedCount     int      `json:"unused_count"`
}

func unusedAddressesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Find unused addresses in a wallet (zero balance, no transaction history)",
		Use:   "unusedAddresses [wallet]",
		Long: `Find addresses in a wallet that have never been used.
An unused address has:
- Zero confirmed balance
- No transaction history (never received or sent coins)

This is useful for generating funding requests where address reuse should be avoided.

Examples:
  skycoin cli unusedAddresses myWallet.wlt
  skycoin cli unusedAddresses myWallet.wlt -n 3
  skycoin cli unusedAddresses myWallet.wlt --generate`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE:         findUnusedAddresses,
	}

	cmd.Flags().IntP("num", "n", 0, "Maximum number of unused addresses to return (0 = all)")
	cmd.Flags().BoolP("generate", "g", false, "Generate new addresses if not enough unused addresses exist")

	return cmd
}

func findUnusedAddresses(c *cobra.Command, args []string) error {
	walletID := args[0]

	num, err := c.Flags().GetInt("num")
	if err != nil {
		return err
	}
	if num < 0 {
		return errors.New("num must be >= 0")
	}

	generate, err := c.Flags().GetBool("generate")
	if err != nil {
		return err
	}

	// Get wallet info
	wlt, err := apiClient.Wallet(walletID)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	// Collect all addresses from wallet
	var allAddrs []string
	for _, e := range wlt.Entries {
		allAddrs = append(allAddrs, e.Address)
	}

	if len(allAddrs) == 0 {
		return errors.New("wallet has no addresses")
	}

	// Find unused addresses
	unusedAddrs, err := findUnusedAddrsInList(allAddrs)
	if err != nil {
		return err
	}

	// If we need more addresses and --generate is set
	if generate && num > 0 && len(unusedAddrs) < num {
		needed := num - len(unusedAddrs)
		_, _ = fmt.Fprintf(c.ErrOrStderr(), "Generating %d new address(es)...\n", needed) //nolint:errcheck

		// Generate new addresses via API (empty password for unencrypted wallets)
		newAddrs, err := apiClient.NewWalletAddress(walletID, "", wallet.OptionGenerateN(uint64(needed))) //nolint:gosec // needed is guaranteed positive by preceding checks
		if err != nil {
			return fmt.Errorf("failed to generate new addresses: %w", err)
		}

		// Add newly generated addresses (they should all be unused)
		unusedAddrs = append(unusedAddrs, newAddrs...)
	}

	// Limit results if num > 0
	if num > 0 && len(unusedAddrs) > num {
		unusedAddrs = unusedAddrs[:num]
	}

	result := UnusedAddressResult{
		Wallet:          walletID,
		TotalAddresses:  len(allAddrs),
		UnusedAddresses: unusedAddrs,
		UnusedCount:     len(unusedAddrs),
	}

	return printJSON(result)
}

// findUnusedAddrsInList checks a list of addresses and returns those that are unused
// (zero balance and no transaction history)
func findUnusedAddrsInList(addrs []string) ([]string, error) {
	if len(addrs) == 0 {
		return nil, nil
	}

	// Get balances for all addresses in one call
	balances, err := GetBalanceOfAddresses(apiClient, addrs)
	if err != nil {
		return nil, fmt.Errorf("failed to get balances: %w", err)
	}

	// Build a map of addresses with zero balance
	zeroBalanceAddrs := make(map[string]bool)
	for _, ab := range balances.Addresses {
		if ab.Confirmed.Coins == "0.000000" {
			zeroBalanceAddrs[ab.Address] = true
		}
	}

	// For addresses with zero balance, check transaction history
	var unusedAddrs []string
	for _, addr := range addrs {
		if !zeroBalanceAddrs[addr] {
			continue // Has balance, not unused
		}

		// Check transaction history
		txns, err := apiClient.TransactionsVerbose([]string{addr})
		if err != nil {
			return nil, fmt.Errorf("failed to get transactions for %s: %w", addr, err)
		}

		// If no transactions, address is unused
		if len(txns) == 0 {
			unusedAddrs = append(unusedAddrs, addr)
		}
	}

	return unusedAddrs, nil
}
