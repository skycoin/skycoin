package params

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
)

type Distribution struct {
	// MaxCoinSupply is the maximum supply of coins
	MaxCoinSupply uint64
	// AddressInitialBalance is the initial balance of each distribution address
	AddressInitialBalance uint64
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

func (d *Distribution) MustValidate() {
	if err := d.Validate(); err != nil {
		panic(err)
	}
}

func (d *Distribution) Validate() error {
	if d.InitialUnlockedCount > uint64(len(d.Addresses)) {
		return errors.New("unlocked addresses > total distribution addresses")
	}

	if d.AddressInitialBalance*uint64(len(d.Addresses)) != d.MaxCoinSupply {
		return errors.New("total balance in distribution addresses != max coin supply")
	}

	if d.MaxCoinSupply%uint64(len(d.Addresses)) != 0 {
		return errors.New("MaxCoinSupply should be perfectly divisible by len(addresses)")
	}

	if err := d.decodeAddresses(); err != nil {
		return err
	}

	return nil
}

// UnlockedAddresses returns distribution addresses that are unlocked, i.e. they have spendable outputs
func (d *Distribution) UnlockedAddresses() []string {
	// The first InitialUnlockedCount (25) addresses are unlocked by default.
	// Subsequent addresses will be unlocked at a rate of UnlockAddressRate (5) per year,
	// after the InitialUnlockedCount (25) addresses have no remaining balance.
	// The unlock timer will be enabled manually once the
	// d.InitialUnlockedCount (25) addresses are distributed.

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
func (d *Distribution) LockedAddresses() []string {
	// TODO -- once we reach 30% distribution, we can hardcode the
	// initial timestamp for releasing more coins
	addrs := make([]string, d.numLocked())
	copy(addrs, d.Addresses[d.InitialUnlockedCount:])
	return addrs
}

// AddressesDecoded returns a copy of the hardcoded distribution addresses array.
// Each address has 1,000,000 coins. There are 100 addresses.
func (d *Distribution) AddressesDecoded() []cipher.Address {
	d.mustDecodeAddresses()
	addrs := make([]cipher.Address, len(d.addressesDecoded))
	copy(addrs, d.addressesDecoded)
	return addrs
}

// UnlockedAddressesDecoded returns distribution addresses that are unlocked, i.e. they have spendable outputs
func (d *Distribution) UnlockedAddressesDecoded() []cipher.Address {
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
	d.mustDecodeAddresses()
	addrs := make([]cipher.Address, d.InitialUnlockedCount)
	copy(addrs[:], d.addressesDecoded[:d.InitialUnlockedCount])
	return addrs
}

// LockedAddressesDecoded returns distribution addresses that are locked, i.e. they have unspendable outputs
func (d *Distribution) LockedAddressesDecoded() []cipher.Address {
	// TODO -- once we reach 30% distribution, we can hardcode the
	// initial timestamp for releasing more coins
	d.mustDecodeAddresses()
	addrs := make([]cipher.Address, d.numLocked())
	copy(addrs, d.addressesDecoded[d.InitialUnlockedCount:])
	return addrs
}

func (d *Distribution) numLocked() uint64 {
	n := uint64(len(d.Addresses))
	if n < d.InitialUnlockedCount {
		panic("number of distribution addresses is less than InitialUnlockedCount")
	}
	return n - d.InitialUnlockedCount
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

func (d *Distribution) mustDecodeAddresses() {
	if err := d.decodeAddresses(); err != nil {
		panic(err)
	}
}

// addressesDecoded is initialized in init.go from params.go's addresses
// var addressesDecoded []cipher.Address

// // GetDistributionAddresses returns a copy of the hardcoded distribution addresses array.
// // Each address has 1,000,000 coins. There are 100 addresses.
// func GetDistributionAddresses(p *Params) []string {
// 	addrs := make([]string, len(p.DistributionAddresses))
// 	for i := range p.DistributionAddresses {
// 		addrs[i] = p.DistributionAddresses[i]
// 	}
// 	return addrs
// }

// // GetUnlockedDistributionAddresses returns distribution addresses that are unlocked, i.e. they have spendable outputs
// func GetUnlockedDistributionAddresses(p *Params) []string {
// 	// The first InitialUnlockedCount (25) addresses are unlocked by default.
// 	// Subsequent addresses will be unlocked at a rate of UnlockAddressRate (5) per year,
// 	// after the InitialUnlockedCount (25) addresses have no remaining balance.
// 	// The unlock timer will be enabled manually once the
// 	// p.InitialUnlockedCount (25) addresses are distributed.

// 	// NOTE: To have automatic unlocking, transaction verification would have
// 	// to be handled in visor rather than in coin.Transactions.Visor(), because
// 	// the coin package is agnostic to the state of the blockchain and cannot reference it.
// 	// Instead of automatic unlocking, we can hardcode the timestamp at which the first 30%
// 	// is distributed, then compute the unlocked addresses easily here.

// 	addrs := make([]string, p.InitialUnlockedCount)
// 	copy(addrs[:], p.DistributionAddresses[:p.InitialUnlockedCount])
// 	return addrs
// }

// // GetLockedDistributionAddresses returns distribution addresses that are locked, i.e. they have unspendable outputs
// func GetLockedDistributionAddresses(p *Params) []string {
// 	// TODO -- once we reach 30% distribution, we can hardcode the
// 	// initial timestamp for releasing more coins
// 	addrs := make([]string, p.DistributionAddressesTotal-p.InitialUnlockedCount)
// 	for i := range addresses[p.InitialUnlockedCount:] {
// 		addrs[i] = addresses[p.InitialUnlockedCount+uint64(i)]
// 	}
// 	return addrs
// }

// // GetDistributionAddressesDecoded returns a copy of the hardcoded distribution addresses array.
// // Each address has 1,000,000 coins. There are 100 addresses.
// func GetDistributionAddressesDecoded(p *Params) []cipher.Address {
// 	addrs := make([]cipher.Address, len(p.addressesDecoded))
// 	copy(addrs, p.addressesDecoded)
// 	return addrs
// }

// // GetUnlockedDistributionAddressesDecoded returns distribution addresses that are unlocked, i.e. they have spendable outputs
// func GetUnlockedDistributionAddressesDecoded(p *Params) []cipher.Address {
// 	// The first p.InitialUnlockedCount (25) addresses are unlocked by default.
// 	// Subsequent addresses will be unlocked at a rate of UnlockAddressRate (5) per year,
// 	// after the p.InitialUnlockedCount (25) addresses have no remaining balance.
// 	// The unlock timer will be enabled manually once the
// 	// p.InitialUnlockedCount (25) addresses are distributed.

// 	// NOTE: To have automatic unlocking, transaction verification would have
// 	// to be handled in visor rather than in coin.Transactions.Visor(), because
// 	// the coin package is agnostic to the state of the blockchain and cannot reference it.
// 	// Instead of automatic unlocking, we can hardcode the timestamp at which the first 30%
// 	// is distributed, then compute the unlocked addresses easily here.

// 	addrs := make([]cipher.Address, p.InitialUnlockedCount)
// 	copy(addrs[:], p.addressesDecoded[:p.InitialUnlockedCount])
// 	return addrs
// }

// // GetLockedDistributionAddressesDecoded returns distribution addresses that are locked, i.e. they have unspendable outputs
// func GetLockedDistributionAddressesDecoded(p *Params) []cipher.Address {
// 	// TODO -- once we reach 30% distribution, we can hardcode the
// 	// initial timestamp for releasing more coins
// 	addrs := make([]cipher.Address, p.DistributionAddressesTotal-p.InitialUnlockedCount)
// 	copy(addrs, p.addressesDecoded[p.InitialUnlockedCount:])
// 	return addrs
// }
