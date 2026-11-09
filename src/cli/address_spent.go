package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/droplet"
)

// AddressSpentResult says whether an address has ever sent coins
type AddressSpentResult struct {
	Address      string `json:"address"`
	Spent        bool   `json:"spent"`
	Received     int    `json:"received_outputs"`
	SpentOutputs int    `json:"spent_outputs"`
	Balance      string `json:"balance"`
}

func addressSpentCmd() *cobra.Command {
	return &cobra.Command{
		Short: "Check whether addresses have ever sent coins",
		Use:   "addressSpent [addresses]",
		Long: `Check whether each address has ever sent coins. An address has sent
coins once any output it received was spent, which also puts its public key
on chain.
    example: addressSpent addr1 addr2`,
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  addrSpent,
	}
}

func addrSpent(_ *cobra.Command, addrs []string) error {
	results := make([]AddressSpentResult, 0, len(addrs))
	for _, a := range addrs {
		if _, err := cipher.DecodeBase58Address(a); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", a, err)
		}
		uxouts, err := apiClient.AddressUxOuts(a)
		if err != nil {
			return err
		}
		r, err := summarizeAddressSpent(a, uxouts)
		if err != nil {
			return err
		}
		results = append(results, r)
	}
	return printJSON(results)
}

func summarizeAddressSpent(addr string, uxouts []readable.SpentOutput) (AddressSpentResult, error) {
	r := AddressSpentResult{Address: addr, Received: len(uxouts)}
	var bal uint64
	for _, ux := range uxouts {
		if ux.SpentBlockSeq > 0 {
			r.Spent = true
			r.SpentOutputs++
		} else {
			bal += ux.Coins
		}
	}
	s, err := droplet.ToString(bal)
	if err != nil {
		return r, err
	}
	r.Balance = s
	return r, nil
}
