package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip32"
)

// NextAddressResult represents the result of finding the next unused address
type NextAddressResult struct {
	Address    string `json:"address"`
	ChildIndex uint32 `json:"child_index"`
	PublicKey  string `json:"public_key"`
}

func nextAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Derive the next unused address from an xpub key",
		Use:   "nextAddress [xpub key]",
		Long: `Derive child addresses from a bip32 xpub key and return the first unused one.

An unused address has zero balance and no transaction history.
Derivation starts at child index 0 (or the specified start index) and scans
forward until an unused address is found.

Requires a running skycoin node to check address activity.

Examples:
  skycoin cli nextAddress xpub6CJWevR9X57j...
  skycoin cli nextAddress xpub6CJWevR9X57j... --start 10
  skycoin cli nextAddress xpub6CJWevR9X57j... --batch-size 50`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE:         findNextAddress,
	}

	cmd.Flags().Uint32("start", 0, "Child index to start scanning from")
	cmd.Flags().Int("batch-size", 20, "Number of addresses to check per batch")

	return cmd
}

func findNextAddress(c *cobra.Command, args []string) error {
	xpubStr := args[0]

	startIdx, err := c.Flags().GetUint32("start")
	if err != nil {
		return err
	}

	batchSize, err := c.Flags().GetInt("batch-size")
	if err != nil {
		return err
	}
	if batchSize < 1 {
		return errors.New("batch-size must be >= 1")
	}

	xpub, err := bip32.DeserializeEncodedPublicKey(xpubStr)
	if err != nil {
		return fmt.Errorf("invalid xpub key: %w", err)
	}

	childIdx := startIdx
	for {
		// Derive a batch of addresses
		var addrs []string
		type childInfo struct {
			addr   string
			index  uint32
			pubKey cipher.PubKey
		}
		var children []childInfo

		for i := 0; i < batchSize; i++ {
			pk, err := xpub.NewPublicChildKey(childIdx)
			if err != nil {
				if bip32.IsImpossibleChildError(err) {
					childIdx++
					i-- // retry with next index
					continue
				}
				return fmt.Errorf("failed to derive child key at index %d: %w", childIdx, err)
			}

			pubKey := cipher.MustNewPubKey(pk.Key)
			addr := cipher.AddressFromPubKey(pubKey)
			addrStr := addr.String()

			addrs = append(addrs, addrStr)
			children = append(children, childInfo{
				addr:   addrStr,
				index:  childIdx,
				pubKey: pubKey,
			})
			childIdx++
		}

		// Check which addresses are unused
		unusedAddrs, err := findUnusedAddrsInList(addrs)
		if err != nil {
			return fmt.Errorf("failed to check address activity: %w", err)
		}

		if len(unusedAddrs) > 0 {
			// Find the first unused address in derivation order
			unusedSet := make(map[string]struct{}, len(unusedAddrs))
			for _, a := range unusedAddrs {
				unusedSet[a] = struct{}{}
			}

			for _, child := range children {
				if _, ok := unusedSet[child.addr]; ok {
					return printJSON(NextAddressResult{
						Address:    child.addr,
						ChildIndex: child.index,
						PublicKey:  child.pubKey.Hex(),
					})
				}
			}
		}

		// Safety: don't scan forever
		if childIdx-startIdx > 10000 {
			return errors.New("scanned 10000 addresses without finding an unused one")
		}
	}
}
