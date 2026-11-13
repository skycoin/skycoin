// Package params src/params/distribution.go c1-coin-chain
//
// Package params defines blockchain parameters and distribution configuration.
package params

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
)

var errUnlockedCount = errors.New("unlocked addresses > total distribution addresses")

// Distribution parameters define the initial coin distribution and unlocking schedule
type Distribution struct {
	// MaxCoinSupply is the maximum supply of coins
	MaxCoinSupply uint64
	// InitialUnlockedCount is the initial number of unlocked addresses
	InitialUnlockedCount uint64
	// UnlockAddressRate is the number of addresses to unlock per unlock time interval
	UnlockAddressRate uint64
	// UnlockTimeInterval is the distribution address unlock time interval, measured in seconds
	// Once the InitialUnlockedCount is exhausted,
	// UnlockAddressRate addresses will be unlocked per UnlockTimeInterval
	UnlockTimeInterval uint64

	// Addresses are the distribution addresses that received coins in the
	// first block after the genesis block
	Addresses        []string
	addressesDecoded []cipher.Address
}

// Validate validates Distribution parameters
func (d *Distribution) Validate() error {
	if d.InitialUnlockedCount > uint64(len(d.Addresses)) {
		return errUnlockedCount
	}

	if d.MaxCoinSupply%uint64(len(d.Addresses)) != 0 {
		return errors.New("MaxCoinSupply should be perfectly divisible by len(addresses)")
	}

	if err := d.decodeAddresses(); err != nil {
		return err
	}

	return nil
}

// AddressInitialBalance is the initial balance of each distribution address
func (d *Distribution) AddressInitialBalance() uint64 {
	return d.MaxCoinSupply / uint64(len(d.Addresses))
}

// UnlockedAddresses returns distribution addresses that are unlocked, i.e. they have spendable outputs
func (d *Distribution) UnlockedAddresses() []string {
	// The first InitialUnlockedCount (25) addresses are unlocked by default.
	// Subsequent addresses will be unlocked at a rate of UnlockAddressRate (5) per year,
	// after the InitialUnlockedCount (25) addresses have no remaining balance.
	// The unlock timer will be enabled manually once the
	// InitialUnlockedCount (25) addresses are distributed.

	// NOTE: To have automatic unlocking, transaction verification would have
	// to be handled in visor rather than in coin.Transactions.Visor(), because
	// the coin package is agnostic to the state of the blockchain and cannot reference it.
	// Instead of automatic unlocking, we can hardcode the timestamp at which the first 30%
	// is distributed, then compute the unlocked addresses easily here.

	addrs := make([]string, d.InitialUnlockedCount)
	copy(addrs[:], d.Addresses[:d.InitialUnlockedCount])
	return addrs
}

// LockedAddresses returns distribution addresses that are locked, i.e. they have unspendable outputs
func (d *Distribution) LockedAddresses() ([]string, error) {
	// TODO -- once we reach 30% distribution, we can hardcode the
	// initial timestamp for releasing more coins
	n, err := d.numLocked()
	if err != nil {
		return nil, err
	}
	addrs := make([]string, n)
	copy(addrs, d.Addresses[d.InitialUnlockedCount:])
	return addrs, nil
}

// AddressesDecoded returns a copy of the hardcoded distribution addresses array.
// Each address has 1,000,000 coins. There are 100 addresses.
func (d *Distribution) AddressesDecoded() ([]cipher.Address, error) {
	if err := d.decodeAddresses(); err != nil {
		return nil, err
	}
	addrs := make([]cipher.Address, len(d.addressesDecoded))
	copy(addrs, d.addressesDecoded)
	return addrs, nil
}

// UnlockedAddressesDecoded returns distribution addresses that are unlocked, i.e. they have spendable outputs
func (d *Distribution) UnlockedAddressesDecoded() ([]cipher.Address, error) {
	// The first d.InitialUnlockedCount (25) addresses are unlocked by default.
	// Subsequent addresses will be unlocked at a rate of UnlockAddressRate (5) per year,
	// after the d.InitialUnlockedCount (25) addresses have no remaining balance.
	// The unlock timer will be enabled manually once the
	// d.InitialUnlockedCount (25) addresses are distributed.

	// NOTE: To have automatic unlocking, transaction verification would have
	// to be handled in visor rather than in coin.Transactions.Visor(), because
	// the coin package is agnostic to the state of the blockchain and cannot reference it.
	// Instead of automatic unlocking, we can hardcode the timestamp at which the first 30%
	// is distributed, then compute the unlocked addresses easily here.
	if err := d.decodeAddresses(); err != nil {
		return nil, err
	}
	if d.InitialUnlockedCount > uint64(len(d.addressesDecoded)) {
		return nil, errUnlockedCount
	}
	addrs := make([]cipher.Address, d.InitialUnlockedCount)
	copy(addrs[:], d.addressesDecoded[:d.InitialUnlockedCount])
	return addrs, nil
}

// LockedAddressesDecoded returns distribution addresses that are locked, i.e. they have unspendable outputs
func (d *Distribution) LockedAddressesDecoded() ([]cipher.Address, error) {
	// TODO -- once we reach 30% distribution, we can hardcode the
	// initial timestamp for releasing more coins
	if err := d.decodeAddresses(); err != nil {
		return nil, err
	}
	n, err := d.numLocked()
	if err != nil {
		return nil, err
	}
	addrs := make([]cipher.Address, n)
	copy(addrs, d.addressesDecoded[d.InitialUnlockedCount:])
	return addrs, nil
}

func (d *Distribution) numLocked() (uint64, error) {
	n := uint64(len(d.Addresses))
	if n < d.InitialUnlockedCount {
		return 0, errUnlockedCount
	}
	return n - d.InitialUnlockedCount, nil
}

func (d *Distribution) decodeAddresses() error {
	if len(d.addressesDecoded) == len(d.Addresses) {
		return nil
	}

	decodedAddrs := make([]cipher.Address, len(d.Addresses))
	for i, a := range d.Addresses {
		var err error
		decodedAddrs[i], err = cipher.DecodeBase58Address(a)
		if err != nil {
			return err
		}
	}

	d.addressesDecoded = decodedAddrs
	return nil
}
