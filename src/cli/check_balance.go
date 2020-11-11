package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/wallet"
)

// Balance represents an coin and hours balance
type Balance struct {
	Coins string `json:"coins"`
	Hours string `json:"hours"`
}

// AddressBalances represents an address's balance
type AddressBalances struct {
	Confirmed Balance `json:"confirmed"`
	Spendable Balance `json:"spendable"`
	Expected  Balance `json:"expected"`
	Address   string  `json:"address"`
}

// BalanceResult represents an set of addresses' balances
type BalanceResult struct {
	Confirmed Balance           `json:"confirmed"`
	Spendable Balance           `json:"spendable"`
	Expected  Balance           `json:"expected"`
	Addresses []AddressBalances `json:"addresses"`
}

func walletBalanceCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Check the balance of a wallet",
		Use:                   "walletBalance [wallet]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE:                  checkWltBalance,
	}
}

func addressBalanceCmd() *cobra.Command {
	return &cobra.Command{
		Short: "Check the balance of specific addresses",
		Use:   "addressBalance [addresses]",
		Long: `Check balance of specific addresses, join multiple addresses with space.
    example: addressBalance "$addr1 $addr2 $addr3"`,
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  addrBalance,
	}
}

func checkWltBalance(c *cobra.Command, args []string) error {
	w := args[0]
	balRlt, err := CheckWalletBalance(apiClient, w)
	switch err.(type) {
	case nil:
	case WalletLoadError:
		printHelp(c)
		return err
	default:
		return err
	}

	return printJSON(balRlt)
}

func addrBalance(_ *cobra.Command, args []string) error {
	numArgs := len(args)

	addrs := make([]string, numArgs)

	var err error
	for i := 0; i < numArgs; i++ {
		addrs[i] = args[i]
		if _, err = cipher.DecodeBase58Address(addrs[i]); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", addrs[i], err)
		}
	}

	balRlt, err := GetBalanceOfAddresses(apiClient, addrs)
	if err != nil {
		return err
	}

	return printJSON(balRlt)
}

// PUBLIC

// CheckWalletBalance returns the total and individual balances of addresses in a wallet file
func CheckWalletBalance(c GetOutputser, id string) (*BalanceResult, error) {
	wlt, err := apiClient.Wallet(id)
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, e := range wlt.Entries {
		addrs = append(addrs, e.Address)
	}

	return GetBalanceOfAddresses(c, addrs)
}

// GetBalanceOfAddresses returns the total and individual balances of a set of addresses
func GetBalanceOfAddresses(c GetOutputser, addrs []string) (*BalanceResult, error) {
	outs, err := c.OutputsForAddresses(addrs)
	if err != nil {
		return nil, err
	}

	return getBalanceOfAddresses(outs, addrs)
}

func getBalanceOfAddresses(outs *readable.UnspentOutputsSummary, addrs []string) (*BalanceResult, error) {
	addrsMap := make(map[string]struct{}, len(addrs))
	for _, a := range addrs {
		addrsMap[a] = struct{}{}
	}

	addrBalances := make(map[string]struct {
		confirmed, spendable, expected wallet.Balance
	}, len(addrs))

	// Count confirmed balances
	for _, o := range outs.HeadOutputs {
		if _, ok := addrsMap[o.Address]; !ok {
			return nil, fmt.Errorf("Found address %s in GetUnspentOutputs result, but this address wasn't requested", o.Address)
		}

		amt, err := droplet.FromString(o.Coins)
		if err != nil {
			return nil, fmt.Errorf("droplet.FromString failed: %v", err)
		}

		b := addrBalances[o.Address]
		b.confirmed.Coins += amt
		b.confirmed.Hours += o.CalculatedHours

		addrBalances[o.Address] = b
	}

	// Count spendable balances
	for _, o := range outs.SpendableOutputs() {
		if _, ok := addrsMap[o.Address]; !ok {
			return nil, fmt.Errorf("Found address %s in GetUnspentOutputs result, but this address wasn't requested", o.Address)
		}

		amt, err := droplet.FromString(o.Coins)
		if err != nil {
			return nil, fmt.Errorf("droplet.FromString failed: %v", err)
		}

		b := addrBalances[o.Address]
		b.spendable.Coins += amt
		b.spendable.Hours += o.CalculatedHours

		addrBalances[o.Address] = b
	}

	// Count predicted balances
	for _, o := range outs.ExpectedOutputs() {
		if _, ok := addrsMap[o.Address]; !ok {
			return nil, fmt.Errorf("Found address %s in GetUnspentOutputs result, but this address wasn't requested", o.Address)
		}

		amt, err := droplet.FromString(o.Coins)
		if err != nil {
			return nil, fmt.Errorf("droplet.FromString failed: %v", err)
		}

		b := addrBalances[o.Address]
		b.expected.Coins += amt
		b.expected.Hours += o.CalculatedHours

		addrBalances[o.Address] = b
	}

	toBalance := func(b wallet.Balance) (Balance, error) {
		coins, err := droplet.ToString(b.Coins)
		if err != nil {
			return Balance{}, err
		}

		return Balance{
			Coins: coins,
			Hours: strconv.FormatUint(b.Hours, 10),
		}, nil
	}

	var totalConfirmed, totalSpendable, totalExpected wallet.Balance
	balRlt := &BalanceResult{
		Addresses: make([]AddressBalances, len(addrs)),
	}

	for i, a := range addrs {
		b := addrBalances[a]
		var err error

		balRlt.Addresses[i].Address = a

		totalConfirmed, err = totalConfirmed.Add(b.confirmed)
		if err != nil {
			return nil, err
		}

		totalSpendable, err = totalSpendable.Add(b.spendable)
		if err != nil {
			return nil, err
		}

		totalExpected, err = totalExpected.Add(b.expected)
		if err != nil {
			return nil, err
		}

		balRlt.Addresses[i].Confirmed, err = toBalance(b.confirmed)
		if err != nil {
			return nil, err
		}

		balRlt.Addresses[i].Spendable, err = toBalance(b.spendable)
		if err != nil {
			return nil, err
		}

		balRlt.Addresses[i].Expected, err = toBalance(b.expected)
		if err != nil {
			return nil, err
		}
	}

	var err error
	balRlt.Confirmed, err = toBalance(totalConfirmed)
	if err != nil {
		return nil, err
	}

	balRlt.Spendable, err = toBalance(totalSpendable)
	if err != nil {
		return nil, err
	}

	balRlt.Expected, err = toBalance(totalExpected)
	if err != nil {
		return nil, err
	}

	return balRlt, nil
}
