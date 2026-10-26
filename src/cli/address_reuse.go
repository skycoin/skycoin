package cli

import (
	"fmt"
	"io"

	"github.com/skycoin/skycoin/src/readable"
)

// addressUxOutser looks up an address's historical outputs.
type addressUxOutser interface {
	AddressUxOuts(addr string) ([]readable.SpentOutput, error)
}

// warnSpentAddresses warns about each destination that has spent before. Its public
// key is on the blockchain, so coins sent to it are exposed to anyone who can derive
// the private key from that key. The warning does not stop the send.
func warnSpentAddresses(c addressUxOutser, w io.Writer, addrs []string) {
	seen := make(map[string]struct{}, len(addrs))
	for _, a := range addrs {
		if _, ok := seen[a]; ok {
			continue
		}
		seen[a] = struct{}{}

		outs, err := c.AddressUxOuts(a)
		if err != nil {
			fmt.Fprintf(w, "warning: could not check whether %s has spent before: %v\n", a, err) //nolint:errcheck
			continue
		}
		for _, o := range outs {
			if o.SpentBlockSeq != 0 {
				fmt.Fprintf(w, "warning: %s has spent before, so its public key is on the blockchain. "+ //nolint:errcheck
					"Coins sent to it are exposed to anyone who can derive the private key from that public key. "+
					"Ask the recipient for an unused address.\n", a)
				break
			}
		}
	}
}
