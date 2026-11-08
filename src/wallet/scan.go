package wallet

import "github.com/skycoin/skycoin/src/cipher"

// ScanGap calls generate for the next addresses of a chain until gap of them
// in a row have no activity, the address gap limit of
// https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#address-gap-limit
// It returns the generated addresses up to the last one with activity.
func ScanGap(gap uint64, generate func(n uint64) ([]cipher.Addresser, error), tf TransactionsFinder) ([]cipher.Addresser, error) {
	var addrs []cipher.Addresser
	keep := 0
	for uint64(len(addrs)-keep) < gap {
		next, err := generate(gap - uint64(len(addrs)-keep))
		if err != nil {
			return nil, err
		}
		if len(next) == 0 {
			break
		}

		active, err := tf.AddressesActivity(next)
		if err != nil {
			return nil, err
		}
		for i := len(active) - 1; i >= 0; i-- {
			if active[i] {
				keep = len(addrs) + i + 1
				break
			}
		}
		addrs = append(addrs, next...)
	}

	return addrs[:keep], nil
}
